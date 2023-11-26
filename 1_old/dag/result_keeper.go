package dag

import (
	"context"
	"sync"

	"git.in.zhihu.com/antispam/datasupply/node"
)

type IResultKeeper interface {
	Write(ctx context.Context, nodeID string, result node.Result)
	Read() *Result
}

type DefaultResultKeeper struct {
	data sync.Map
}

var _ IResultKeeper = new(DefaultResultKeeper)

func NewDefaultResultKeeper() *DefaultResultKeeper {
	keeper := &DefaultResultKeeper{
		data: sync.Map{},
	}
	return keeper
}

func (keeper *DefaultResultKeeper) Write(ctx context.Context, nodeID string, result node.Result) {
	keeper.data.Store(nodeID, result)
}

func (keeper *DefaultResultKeeper) Read() *Result {
	result := NewResult()
	keeper.data.Range(func(_, _nodeResult any) bool {
		nodeResult, ok := _nodeResult.(node.Result)
		if !ok {
			return true
		}
		for fieldCode, fieldResult := range nodeResult {
			result.Fields[fieldCode] = fieldResult
		}
		return true
	})
	return result
}
