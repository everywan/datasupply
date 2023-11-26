package dag

import (
	"context"
	"sync"

	"git.in.zhihu.com/antispam/datasupply/log"
	"git.in.zhihu.com/antispam/datasupply/node"
	"git.in.zhihu.com/antispam/datasupply/utils"
	"golang.org/x/sync/errgroup"
)

//go:generate mockgen -package mock_dag -destination ./mock/runtime.go -source=runtime.go
type IRuntime interface {
	Run(ctx context.Context, paramMap map[string]interface{}) IRuntime
	WaitStage(ctx context.Context, stageName node.SupplyStage) IRuntime
	Wait(ctx context.Context) IRuntime
	GetResultCopy() *Result
	AddNodeResultMonitor(func(node.INode, node.Result))
	Close()
}

type runtime struct {
	// runtime 输入
	id         string
	root       node.INode
	allNodeCnt int32 // 全部待补数字段数量
	concurrent int   // 最多并发执行几个节点, 0 表示同步, 负数表示不限制并发
	logger     log.ILog

	// 运行时数据, 预留接口但先使用默认实现
	allNodeDone  chan struct{}    // dag 终止条件判断
	nsKeeper     *nodeStateKeeper // 节点状态管理
	stageKeeper  IStageKeeper     // 补数阶段管理
	resultKeeper IResultKeeper    // 补数数据管理

	nodeResultMonitors []func(node.INode, node.Result)

	finishLocker   sync.Locker
	supplyFinished bool
}

var _ IRuntime = new(runtime)

// Run 对所有字段进行补数, 直到全部字段补充完毕才会返回. 正常退出指 runtime 把所有字段补完.
func (rt *runtime) Run(ctx context.Context, paramMap map[string]interface{}) IRuntime {
	rootRuntime := rt.root.CreateRuntime()
	for fieldCode, fieldValue := range paramMap {
		rootRuntime.AddParam(fieldCode, fieldValue)
	}
	rt.nsKeeper.Push(rootRuntime, rt.root.GetPriority())
	utils.SafelyGo(
		func() {
			rt.run(ctx)
		}, func(err error) {
			rt.logger.Errorf(ctx, "dag.run [%s] error [%s]", rt.id, err)
		})

	return rt
}

// tips: cnode == currentNode
func (rt *runtime) run(ctx context.Context) {
	defer func() {
		// 必须要执行的两个条件, 否则可能会造成 groutine 泄漏
		rt.stageKeeper.SetAllDone()
		close(rt.allNodeDone)
	}()

	ctx, cancel := context.WithTimeout(ctx, DAGTimtout)
	defer cancel()
	errGroup := errgroup.Group{}
	errGroup.SetLimit(rt.concurrent)
	for ; rt.allNodeCnt > 0; rt.allNodeCnt-- {
		select {
		case <-ctx.Done():
			// 兜底措施, 整体超时则直接返回
			rt.logger.Errorf(ctx, "dag.runtime [%s] timeout, still have %d node not run",
				rt.id, rt.allNodeCnt)
			rt.allNodeCnt = 0
			break
		default:
			cnodeRuntime := rt.nsKeeper.Pop(ctx)
			if cnodeRuntime == nil {
				// 因为并没有运行节点, 所以放回这次次数.
				rt.allNodeCnt++
				rt.logger.Errorf(ctx,
					"dag.runtime [%s] get_node_from_keeper timeout, still have %d node not run",
					rt.id, rt.allNodeCnt)
				break
			}
			errGroup.Go(func() error {
				err := utils.SafelyRun(func() {
					// 节点执行
					var nodeResult node.Result
					if cnodeRuntime.IsPrune() {
						nodeResult = cnodeRuntime.GetNode().ValueOnPrune()
					} else {
						nodeResult = cnodeRuntime.Run(ctx)
						// 检测与当前节点相关的下游节点
						rt.nsKeeper.Detection(cnodeRuntime.GetNode(), nodeResult)
					}

					// 检查是否导出字段
					for _, field := range cnodeRuntime.GetNode().GetFields() {
						if field.NotExport {
							delete(nodeResult, field.Code)
						}
					}
					rt.resultKeeper.Write(ctx, cnodeRuntime.GetNodeID(), nodeResult)
					for _, fn := range rt.nodeResultMonitors {
						fn(cnodeRuntime.GetNode(), nodeResult)
					}

					// dag stage 检测
					rt.stageKeeper.RecordAfterNodeFinish(ctx, cnodeRuntime.GetNode())
				})
				if err != nil {
					rt.logger.Errorf(ctx, "dag.run [%s] error [%s]", rt.id, err)
					return err
				}
				return nil
			})
		}
	}
	// 这个理论上不应该发生, 因为内部所有的 error 都需要捕获
	if err := errGroup.Wait(); err != nil {
		log.Errorf(ctx, "dag.run on errgroup.wait error [%v]", err)
	}
	rt.afterSupplyFinished()
}

func (rt *runtime) Wait(ctx context.Context) IRuntime {
	select {
	case <-rt.allNodeDone:
		return rt
	case <-ctx.Done():
		return rt
	}
}

func (rt *runtime) WaitStage(ctx context.Context, stageName node.SupplyStage) IRuntime {
	// todo [optimize] 这里丢失了信息.
	_ = rt.stageKeeper.WaitFor(ctx, stageName)
	return rt
}

func (rt *runtime) GetResultCopy() *Result {
	return rt.resultKeeper.Read()
}

func (rt *runtime) AddNodeResultMonitor(fn func(node.INode, node.Result)) {
	rt.nodeResultMonitors = append(rt.nodeResultMonitors, fn)
}

func (rt *runtime) afterSupplyFinished() {
	// 如果已经处理过则不处理了
	if rt.supplyFinished {
		return
	}
	// slow
	rt.finishLocker.Lock()
	defer rt.finishLocker.Unlock()
	if rt.supplyFinished {
		return
	}
	rt.supplyFinished = true
	rt.nsKeeper.Close()
}

func (rt *runtime) Close() {
	if !rt.supplyFinished {
		rt.afterSupplyFinished()
	}
}
