package dag

import (
	"context"
	"sync"
	"time"

	"git.in.zhihu.com/antispam/datasupply/log"
	"git.in.zhihu.com/antispam/datasupply/node"
	"git.in.zhihu.com/antispam/datasupply/utils"
)

// 节点状态管理员
type nodeStateKeeper struct {
	// todo [optimize] `优先级+无极调节+无就绪节点时阻塞` 较难实现, 大顶堆 heap 或许可以
	// 目前先简单这么写, 后续再优化.
	ready map[int]chan node.IRuntime // 就绪态

	wait   sync.Map // 等待态, 即依赖部分就绪就绪
	prune  sync.Map // 被剪枝节点. 无需再运行
	closed chan struct{}
}

func newNodeStateKeeper(size int) *nodeStateKeeper {
	if size == 0 {
		size = 500
	}
	return &nodeStateKeeper{
		ready: map[int]chan node.IRuntime{
			node.PriorityHigh: make(chan node.IRuntime, size),
			node.PriorityMid:  make(chan node.IRuntime, size),
			node.PriorityLow:  make(chan node.IRuntime, size),
			node.PriorityMin:  make(chan node.IRuntime, size),
		},
		wait:   sync.Map{},
		prune:  sync.Map{},
		closed: make(chan struct{}),
	}
}

func (nodeStateKeeper *nodeStateKeeper) Push(runtime node.IRuntime, priority int) {
	// todo [optimize] 这个 delay 应该是按照事件接收时间算, 而不是调度触发时间. 需要把 created 传进来.
	delay := runtime.GetNode().GetDelaySupply()
	if delay == 0 {
		nodeStateKeeper.push(runtime, priority)
		return
	}
	utils.SafelyGo(
		func() {
			select {
			case <-time.After(delay):
				nodeStateKeeper.push(runtime, priority)
			case <-nodeStateKeeper.closed:
				log.Warnf(context.Background(),
					"node [%s] delay failed because dag finished", runtime.GetNodeID())
			}
		},
		func(err error) {
			log.Errorf(context.Background(), "node [%s] delay push panic:[%s]", runtime.GetNodeID(), err)
		},
	)
}

func (nodeStateKeeper *nodeStateKeeper) push(runtime node.IRuntime, priority int) {
	select {
	case <-nodeStateKeeper.closed:
		return
	default:
		if priority >= node.PriorityHigh {
			nodeStateKeeper.ready[node.PriorityHigh] <- runtime
		} else if priority >= node.PriorityMid {
			nodeStateKeeper.ready[node.PriorityMid] <- runtime
		} else if priority >= node.PriorityLow {
			nodeStateKeeper.ready[node.PriorityLow] <- runtime
		} else {
			nodeStateKeeper.ready[node.PriorityMin] <- runtime
		}
		return
	}

}

// 当没有已就绪节点时, 阻塞等待.
func (nodeStateKeeper *nodeStateKeeper) Pop(ctx context.Context) node.IRuntime {
	// todo [optimize] 这并不是优先级队列
	select {
	case runtime := <-nodeStateKeeper.ready[node.PriorityHigh]:
		return runtime
	case runtime := <-nodeStateKeeper.ready[node.PriorityMid]:
		return runtime
	case runtime := <-nodeStateKeeper.ready[node.PriorityLow]:
		return runtime
	case runtime := <-nodeStateKeeper.ready[node.PriorityMin]:
		return runtime
	case <-ctx.Done():
		return nil
	case <-nodeStateKeeper.closed:
		return nil
	}
}

// 下游节点就绪态检测
// important 这里要理解清楚, 当一个节点有多个上游节点时的并发问题, 这也是扫描器的关键.
// todo [optimize] 状态转换规则可以尝试用 状态机 等方案优化下.
func (nodeStateKeeper *nodeStateKeeper) Detection(cnode node.INode, nodeResult node.Result) {
	for _, childNode := range cnode.GetNexts() {
		if _, hasPrune := nodeStateKeeper.prune.Load(childNode.GetID()); hasPrune {
			continue
		}

		// 构建 node_runtime
		var childNodeRuntime node.IRuntime
		childNodeRuntime = childNode.CreateRuntime()
		isPrune := false
		for _, param := range childNode.GetParamVariables() {
			fieldResult, ok := nodeResult[param.FieldName]
			if !ok {
				continue
			}
			if !fieldResult.IsSupplySuccess() {
				// TODO [optimize] (同 param.go) 可以看下这个 onerror 有无更好的处理方式.
				var paramValue interface{}
				isPrune, paramValue = param.HandleError()
				if isPrune {
					for _, cnode := range append(childNode.Prune(), childNode) {
						_, loaded := nodeStateKeeper.prune.LoadOrStore(cnode.GetID(), struct{}{})
						if loaded {
							continue
						}
						cruntime := cnode.CreateRuntime()
						cruntime.SetPrune()
						nodeStateKeeper.Push(cruntime, cnode.GetPriority())
						nodeStateKeeper.wait.Delete(cnode.GetID())
					}
					break
				}
				childNodeRuntime.AddParam(param.FieldName, paramValue)
				continue
			}
			childNodeRuntime.AddParam(param.FieldName, fieldResult.Value)
		}
		if isPrune {
			continue
		}

		// -------- 检测 node_runtime.is_ready ------
		if childNodeRuntime.IsReady() {
			nodeStateKeeper.Push(childNodeRuntime, childNode.GetPriority())
			continue
		}
		// 当依赖大于一个时, 需要注意并发问题
		_childNodeRuntime, hasStored := nodeStateKeeper.wait.LoadOrStore(childNode.GetID(),
			childNodeRuntime)
		if !hasStored {
			continue
		}
		childNodeRuntime = _childNodeRuntime.(node.IRuntime)
		// 整合新老 node_runtime
		childNodeRuntime.Merge(nodeResult)
		if !childNodeRuntime.IsReady() {
			continue
		}
		nodeStateKeeper.Push(childNodeRuntime, childNode.GetPriority())
		nodeStateKeeper.wait.Delete(childNode.GetID())
	}
}

func (ns *nodeStateKeeper) Close() {
	close(ns.closed)
	for _, ch := range ns.ready {
		close(ch)
	}
}
