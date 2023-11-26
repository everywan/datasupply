package dag

import (
	"context"

	"git.in.zhihu.com/antispam/datasupply/log"
	"git.in.zhihu.com/antispam/datasupply/node"
)

type dagBuilder struct {
	logger log.ILog
}

func newDagBuilder(logger log.ILog) *dagBuilder {
	return &dagBuilder{
		logger: logger,
	}
}

func (builder dagBuilder) build(root node.INode, nodes []node.INode) (node.INode, []node.INode) {
	// 节点集优化
	nodes = builder.mergeNodes(nodes)

	// 构建节点依赖关系
	builder.analyseNodeDep(root, nodes)

	// 孤儿节点检查和移除
	orphanNodeMap := builder.analyseOrphanNode(nodes)
	builder.removeOrphanNode(orphanNodeMap)

	// 重新设置节点补数阶段
	builder.resetSupplyMode(root)

	// 计算节点权重
	builder.calculateNodePriority(nodes)
	return root, nodes
}

// 合并有相同函数调用的节点
func (builder dagBuilder) mergeNodes(nodes []node.INode) []node.INode {
	nodeMap := make(map[string]node.INode, len(nodes))
	for _, _cnode := range nodes {
		nodeid := _cnode.GetID()
		cnode, ok := nodeMap[nodeid]
		if !ok {
			nodeMap[nodeid] = _cnode
			continue
		}
		cnode.AddFields(_cnode.GetFields()...)
	}
	newNodes := make([]node.INode, 0, len(nodeMap))
	for _, cnode := range nodeMap {
		newNodes = append(newNodes, cnode)
	}
	return newNodes
}

func (builder dagBuilder) analyseNodeDep(root node.INode, nodes []node.INode) {
	// dag 内节点可以获取的参数集合, 包括外界输入的 inputs, 节点产生的 fields
	fieldMap := make(map[string]node.INode)
	for _, cnode := range append(nodes, root) {
		for _, fieldCode := range cnode.GetFieldCodes() {
			fieldMap[fieldCode] = cnode
		}
	}
	for _, cnode := range nodes {
		paramVars := cnode.GetParamVariables()
		if len(paramVars) == 0 {
			builder.logger.Infof(context.Background(),
				"node [%s] have zero var_params, add to root's child", cnode.GetID())
			root.AddNexts(cnode)
			cnode.AddPrevs(root)
			continue
		}
		for _, param := range paramVars {
			parentNode, ok := fieldMap[param.FieldName]
			if !ok {
				// param_var 不是其他节点产生的 && 不是外界输入的, 跳过执行
				builder.logger.Warnf(context.Background(),
					"node [%s] param [%s] not found in dag", cnode.GetID(), param.FieldName)
				continue
			}
			parentNode.AddNexts(cnode)
			cnode.AddPrevs(parentNode)
		}
	}
}

/*
	孤儿节点检查. 满足如下任意条件即为孤儿节点
	1. 至少有一个变量参数即不是由其他节点产生, 也不是外界输入的节点叫做孤立节点.
	2. 祖先节点中存在孤儿节点.
*/
func (builder dagBuilder) analyseOrphanNode(nodes []node.INode) map[string]*node.OrphanNode {
	orphanNodeMap := map[string]*node.OrphanNode{}
	for _, cnode := range nodes {
		// 当 node.params 与 node.prevs.fields 不匹配时, 剪枝并记录.
		paramVars := cnode.GetParamVariables()
		paramFoundMap := make(map[string]bool, len(paramVars))
		for _, param := range paramVars {
			paramFoundMap[param.FieldName] = false
		}
		// 从所有父节点中寻找当前参数, 找到则将该参数标记为 true.
		for _, parentNode := range cnode.GetPrevs() {
			for _, fieldCode := range parentNode.GetFieldCodes() {
				if _, ok := paramFoundMap[fieldCode]; ok {
					paramFoundMap[fieldCode] = true
				}
			}
		}

		// 确认是否存在找不到的参数.
		notFoundParams := []string{}
		for param, isFound := range paramFoundMap {
			if !isFound {
				builder.logger.Warnf(context.Background(),
					"found orphan node: [%s], reason: param [%s] not found in dag", cnode.GetID(), param)
				notFoundParams = append(notFoundParams, param)
			}
		}
		if len(notFoundParams) == 0 {
			continue
		}

		// 生成孤儿节点
		orphanNodeMap[cnode.GetID()] = node.NewOrphanNode(cnode,
			node.OrphanReasonNotEnoughParam, notFoundParams)
		for _, pruneNode := range cnode.Prune() {
			orphanNodeMap[pruneNode.GetID()] = node.NewOrphanNode(pruneNode,
				node.OrphanReasonAncestorPruned, []string{})
		}
	}

	return orphanNodeMap
}

// 移除孤儿节点. 如果一个孤儿节点有父节点, 则从其父节点中将该节点剥离
func (builder dagBuilder) removeOrphanNode(orphanNodes map[string]*node.OrphanNode) {
	for _, orphanNode := range orphanNodes {
		for _, parentNode := range orphanNode.Node.GetPrevs() {
			if _, ok := orphanNodes[parentNode.GetID()]; ok {
				// 如果父节点也是孤儿节点, 则无需处理.
				continue
			}
			parentNode.RemoveNexts(orphanNode.Node)
		}
	}
}

func (builder dagBuilder) resetSupplyMode(root node.INode) {
	bottomNodes := []node.INode{}
	for _, cnode := range root.Prune() {
		if len(cnode.GetNexts()) == 0 {
			bottomNodes = append(bottomNodes, cnode)
		}
	}

	i := 0
	for {
		if i >= len(bottomNodes) {
			break
		}
		cnode := bottomNodes[i]
		for _, cPrevNode := range cnode.GetPrevs() {
			// 当父节点补数阶段在子节点阶段之后, 调整父节点补数阶段与子节点相同
			// parent.stage 在 child.stage 之后时, 将父节点 stage 调前.
			if cPrevNode.GetSupplyStage() > cnode.GetSupplyStage() {
				cPrevNode.SetSupplyStage(cnode.GetSupplyStage())
				log.Warnf(context.Background(), "cnode [%s] supply_stage reset to [%s.%s]",
					cPrevNode.GetID(), cnode.GetID(), cnode.GetSupplyStage())
			}
			bottomNodes = append(bottomNodes, cPrevNode)
		}
		i++
	}
}

// 计算节点权重, 根据不同需求进行分级, 每个分级内根据影响吞吐量的权重字段进行详细划分.
// 目前先简单设为两级. 后续有需求再优化.
func (builder dagBuilder) calculateNodePriority(nodes []node.INode) {
	for _, cnode := range nodes {
		// 手动指定权重, 则不重新生成
		if cnode.GetPriority() != 0 {
			continue
		}
		switch cnode.GetSupplyStage() {
		case node.SupplyStageSync:
			cnode.SetPriority(node.PriorityHigh)
		case node.SupplyStageAsync:
			cnode.SetPriority(node.PriorityMid)
		case node.SupplyStageStore:
			cnode.SetPriority(node.PriorityLow)
		default:
			cnode.SetPriority(node.PriorityMin)
		}
	}
}
