package node

import (
	"context"
	"sync"
	"sync/atomic"
)

type IRuntime interface {
	GetNode() INode
	GetNodeID() string
	AddParam(fieldCode string, value interface{})
	IsReady() bool
	Run(context.Context) Result
	Merge(Result)
	IsPrune() bool
	SetPrune()
}

// runtime node 的运行时
type Runtime struct {
	isPrune   bool
	varParams map[string]interface{} // 存储需要准备的参数, key == paran.id
	node      INode

	nodeid   string
	paramCnt int64
	locker   sync.Locker
}

var _ IRuntime = new(Runtime)

func (runtime *Runtime) GetNode() INode {
	return runtime.node
}

func (runtime *Runtime) GetNodeID() string {
	return runtime.nodeid
}

func (runtime *Runtime) AddParam(fieldCode string, value interface{}) {
	runtime.locker.Lock()
	defer runtime.locker.Unlock()
	runtime.varParams[fieldCode] = value
	// todo 去掉这个吧, 并发写也会 panic...
	atomic.AddInt64(&runtime.paramCnt, -1)
}

// 不管如何操作, 最后一个上游依赖的 IsReady 判断一定是在所有的其他读写之后, 从而保证并发安全.
func (runtime *Runtime) IsReady() bool {
	return atomic.LoadInt64(&runtime.paramCnt) == 0
}

func (runtime *Runtime) Run(ctx context.Context) Result {
	if !runtime.IsReady() {
		return runtime.node.ValueOnError("params not ready")
	}
	return runtime.node.Run(ctx, runtime.varParams)
}

// todo [optimize] 这里写的很差, 需要改. 包括调用这个函数的地方.
func (runtime *Runtime) Merge(results Result) {
	for _, param := range runtime.node.GetParamVariables() {
		if fieldResult, ok := results[param.FieldName]; ok {
			runtime.AddParam(param.FieldName, fieldResult.Value)
		}
	}
}

func (runtime *Runtime) IsPrune() bool {
	return runtime.isPrune
}

func (runtime *Runtime) SetPrune() {
	runtime.isPrune = true
}
