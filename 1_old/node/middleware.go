package node

import (
	"context"
	"fmt"
	"math/rand"
	"time"

	"git.in.zhihu.com/antispam/datasupply/constant"
	"git.in.zhihu.com/antispam/datasupply/log"
	"git.in.zhihu.com/antispam/datasupply/statsd"
	"git.in.zhihu.com/antispam/datasupply/utils"
)

type Handler = func(ctx context.Context, paramMap map[string]interface{}) Result
type MiddlewareFunc = func(next Handler) Handler

type IMiddleware interface {
	// 中间件需要外界输入的参数
	Name() string
	Params() []Param
	MiddlewareFunc() MiddlewareFunc
}

type Middleware struct {
	name   string
	params []Param
	mf     MiddlewareFunc
}

var _ IMiddleware = new(Middleware)

func NewMiddleware(name string, params []Param, mf MiddlewareFunc) *Middleware {
	if name == "" {
		name = fmt.Sprint(rand.Intn(1000))
	}
	if params == nil {
		params = []Param{}
	}
	return &Middleware{
		name:   name,
		params: params,
		mf:     mf,
	}
}

func (mw *Middleware) Name() string {
	return mw.name
}

func (mw *Middleware) Params() []Param {
	return mw.params
}

func (mw *Middleware) MiddlewareFunc() MiddlewareFunc {
	return mw.mf
}

// 打点中间点. 记录 node 执行次数, 失败次数, 执行时间(包含下游中间件), 字段执行失败统计.
func StatsdMiddleware(statd statsd.IStatsd, nodeName string) IMiddleware {
	nodeRunCount := fmt.Sprintf(constant.NodeRunCount, nodeName)
	nodeRunSpeed := fmt.Sprintf(constant.NodeRunSpeed, nodeName)
	nodeRunError := fmt.Sprintf(constant.NodeRunError, nodeName)
	mf := func(next Handler) Handler {
		return func(ctx context.Context, paramMap map[string]interface{}) Result {
			statd.Increment(nodeRunCount)
			defer statd.TimingUtilNow(nodeRunSpeed, time.Now())
			result := next(ctx, paramMap)
			// 记录字段补数失败次数
			hasFieldError := false
			for field, fieldResult := range result {
				if !fieldResult.IsSupplySuccess() {
					hasFieldError = true
					statd.Increment(fmt.Sprintf(constant.NodeFieldSupplyError, nodeName+"."+field))
				}
			}
			if hasFieldError {
				statd.Increment(nodeRunError)
			}
			return result
		}
	}
	return NewMiddleware("statsd", []Param{}, mf)
}

// 日志中间件, 记录 params 和 result
func LogMiddleware(logger log.ILog, nodeName string) IMiddleware {
	mf := func(next Handler) Handler {
		return func(ctx context.Context, paramMap map[string]interface{}) Result {
			logger.Infof(ctx, "node [%s] start, params: [%s]",
				nodeName, utils.StructToString(paramMap))
			result := next(ctx, paramMap)
			logger.Infof(ctx, "node [%s] end, result: [%s]",
				nodeName, utils.StructToString(result))
			return result
		}
	}
	return NewMiddleware("logger", []Param{}, mf)
}
