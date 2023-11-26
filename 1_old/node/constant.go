package node

import (
	"encoding/json"
	"errors"
	"strconv"
)

// todo [optimize] 这个是不是应该放在 dag 或单独出来?
// 补数阶段共分为 同步->异步->落库->延迟, 数值从小到大
type SupplyStage int

const (
	SupplyStageSync SupplyStage = iota
	SupplyStageAsync
	SupplyStageStore
	SupplyStageLazy
)

var SupplyStageNames = []string{
	SupplyStageSync:  "sync",
	SupplyStageAsync: "async",
	SupplyStageStore: "store",
	SupplyStageLazy:  "lazy",
}

func (s SupplyStage) String() string {
	if int(s) < len(SupplyStageNames) {
		return SupplyStageNames[s]
	}
	return "supply_mode_" + strconv.Itoa(int(s))
}

func (s SupplyStage) MarshalJSON() ([]byte, error) {
	return json.Marshal(s.String())
}

func (s *SupplyStage) UnmarshalJSON(b []byte) error {
	str := ""
	if err := json.Unmarshal(b, &str); err != nil {
		return err
	}
	for mode, name := range SupplyStageNames {
		if name == str {
			*s = SupplyStage(mode)
			return nil
		}
	}
	return errors.New("unknown supply mode " + str)
}

type OnErrorHandler int

const (
	OnErrorDiscard OnErrorHandler = iota
	OnErrorDefault
)

var OnErrorHandlerNames = []string{
	OnErrorDiscard: "discard",
	OnErrorDefault: "default",
}

func (s OnErrorHandler) String() string {
	if int(s) < len(OnErrorHandlerNames) {
		return OnErrorHandlerNames[s]
	}
	return "on_error_handler_" + strconv.Itoa(int(s))
}

func (s OnErrorHandler) MarshalJSON() ([]byte, error) {
	return json.Marshal(s.String())
}

func (s *OnErrorHandler) UnmarshalJSON(b []byte) error {
	str := ""
	if err := json.Unmarshal(b, &str); err != nil {
		return err
	}
	for mode, name := range OnErrorHandlerNames {
		if name == str {
			*s = OnErrorHandler(mode)
			return nil
		}
	}
	return errors.New("unknown on error handler " + str)
}

const (
	PriorityMax  = 100
	PriorityHigh = 75
	PriorityMid  = 50
	PriorityLow  = 25
	PriorityMin  = 1
)
