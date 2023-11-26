package dag

import (
	"context"
	"sync"
	"sync/atomic"

	"git.in.zhihu.com/antispam/datasupply/constant"
	"git.in.zhihu.com/antispam/datasupply/node"
)

type StageFinishCode int

const (
	StageFinishSuccess    = 200
	StageFinishAllSuccess = 201
	StageFinishTimeout    = 408
)

type IStageKeeper interface {
	RecordAfterNodeFinish(context.Context, node.INode)
	WaitFor(context.Context, node.SupplyStage) StageFinishCode
	SetAllDone()
}

type (
	DefaultStageKeeper struct {
		allNodeDoneCh   chan struct{}
		allNodeDoneOnce sync.Once
		stageMap        map[node.SupplyStage]*stage
	}
	// dag stage 变更
	stage struct {
		cnt  int32
		once sync.Once
		done chan struct{}
	}
)

var _ IStageKeeper = new(DefaultStageKeeper)

func NewDefaultStageKeeper(stageCnts map[node.SupplyStage]int32) *DefaultStageKeeper {
	stageMap := make(map[node.SupplyStage]*stage, len(stageCnts))
	for stageName, cnt := range stageCnts {
		done := make(chan struct{})
		if cnt == 0 {
			done = constant.Closedchan
		}
		stageMap[stageName] = &stage{
			cnt:  cnt,
			done: done,
			once: sync.Once{},
		}
	}
	return &DefaultStageKeeper{
		allNodeDoneCh:   make(chan struct{}),
		allNodeDoneOnce: sync.Once{},
		stageMap:        stageMap,
	}
}

func (keeper *DefaultStageKeeper) RecordAfterNodeFinish(ctx context.Context, cnode node.INode) {
	stage, ok := keeper.stageMap[cnode.GetSupplyStage()]
	if ok {
		if atomic.AddInt32(&stage.cnt, -1) <= 0 {
			stage.once.Do(func() {
				close(stage.done)
			})
		}
	}
}

func (keeper *DefaultStageKeeper) WaitFor(ctx context.Context, stageName node.SupplyStage) StageFinishCode {
	ctx, canel := context.WithTimeout(ctx, DAGTimtout)
	defer canel()
	stage, ok := keeper.stageMap[stageName]
	if !ok {
		select {
		case <-keeper.allNodeDoneCh:
			return StageFinishAllSuccess
		case <-ctx.Done():
			return StageFinishTimeout
		}
	}

	select {
	case <-stage.done:
		return StageFinishSuccess
	case <-keeper.allNodeDoneCh:
		return StageFinishAllSuccess
	case <-ctx.Done():
		return StageFinishTimeout
	}
}

func (keeper *DefaultStageKeeper) SetAllDone() {
	keeper.allNodeDoneOnce.Do(func() {
		close(keeper.allNodeDoneCh)
	})
}
