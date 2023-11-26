package dag

import "git.in.zhihu.com/antispam/datasupply/node"

type preComputeData struct {
	allNodeCnt      int32 // 全部待补数字段数量
	stageNodeCntMap map[node.SupplyStage]int32
	stageNodeIDMap  map[node.SupplyStage][]string
	field2FieldMap  map[string]*node.Field
	field2NodeMap   map[string]node.INode // fieldCode:node
}

func preCompute(root node.INode) preComputeData {

	allNodes := append(root.Prune(), root)
	allNodeCnt := len(allNodes)

	stageNodeIDMap := make(map[node.SupplyStage][]string, len(node.SupplyStageNames))
	{
		for _, cnode := range allNodes {
			nodeids, ok := stageNodeIDMap[cnode.GetSupplyStage()]
			if !ok {
				stageNodeIDMap[cnode.GetSupplyStage()] = []string{cnode.GetID()}
				continue
			}
			stageNodeIDMap[cnode.GetSupplyStage()] = append(nodeids, cnode.GetID())
		}
	}

	stageNodeCntMap := make(map[node.SupplyStage]int32, len(node.SupplyStageNames))
	{
		for stageName, ids := range stageNodeIDMap {
			stageNodeCntMap[stageName] = int32(len(ids))
		}
	}

	field2NodeMap := make(map[string]node.INode, len(allNodes))
	field2FieldMap := make(map[string]*node.Field, len(allNodes))
	{
		for i, cnode := range allNodes {
			for _, field := range cnode.GetFields() {
				field2NodeMap[field.Code] = allNodes[i]
				field2FieldMap[field.Code] = field
			}
		}
	}

	return preComputeData{
		allNodeCnt:      int32(allNodeCnt),
		stageNodeCntMap: stageNodeCntMap,
		stageNodeIDMap:  stageNodeIDMap,
		field2NodeMap:   field2NodeMap,
		field2FieldMap:  field2FieldMap,
	}
}
