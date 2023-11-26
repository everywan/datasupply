package dag

import (
	"context"
	"fmt"
	"time"

	"git.in.zhihu.com/antispam/datasupply/constant"
	"git.in.zhihu.com/antispam/datasupply/log"
	"git.in.zhihu.com/antispam/datasupply/node"
	"git.in.zhihu.com/antispam/datasupply/statsd"
	"git.in.zhihu.com/antispam/datasupply/utils"
)

type Handler = func(ctx context.Context, runtimeID string, paramMap map[string]interface{}) IRuntime
type Middleware = func(next Handler) Handler

// 打点中间点, 记录次数和耗时
func StatsdMiddleware(statd statsd.IStatsd, dagName string) Middleware {
	dagRunCount := fmt.Sprintf(constant.DAGRunCount, dagName)
	dagRunSpeed := fmt.Sprintf(constant.DAGRunSpeed, dagName)
	stageSpeeds := make(map[string]string, len(node.SupplyStageNames))
	for _, stage := range node.SupplyStageNames {
		stageSpeeds[stage] = fmt.Sprintf(constant.DAGRunStageSpeed, dagName, stage)
	}
	return func(next Handler) Handler {
		return func(ctx context.Context, runtimeID string, paramMap map[string]interface{}) IRuntime {
			statd.Increment(dagRunCount)
			rt := next(ctx, runtimeID, paramMap)
			utils.SafelyGo(
				func() {
					rt.Wait(ctx)
					statd.TimingUtilNow(dagRunSpeed, time.Now())
				},
				func(err error) { log.Error(ctx, "dag.middleware.statsd panic", err) },
			)
			for i, stage := range node.SupplyStageNames {
				i := i
				stage := stage
				utils.SafelyGo(
					func() {
						rt.WaitStage(ctx, node.SupplyStage(i))
						statd.TimingUtilNow(stageSpeeds[stage], time.Now())
					},
					func(err error) { log.Error(ctx, "dag.middleware.statsd panic", err) },
				)
			}
			return rt
		}
	}
}

// 日志中间件
func LogMiddleware(logger log.ILog) Middleware {
	return func(next Handler) Handler {
		return func(ctx context.Context, runtimeID string, paramMap map[string]interface{}) IRuntime {
			logger.Infof(ctx, "dag.runtime [%s] receive event [%s]", runtimeID, utils.StructToString(paramMap))
			rt := next(ctx, runtimeID, paramMap)
			utils.SafelyGo(
				func() {
					result := rt.Wait(ctx).GetResultCopy()
					logger.Infof(ctx, "dag.runtime [%s] supply result: [%s]", runtimeID, utils.StructToString(result))
				},
				func(err error) {
					logger.Errorf(ctx, "dag.middleware.logger panic", err)
				},
			)
			return rt
		}
	}
}
