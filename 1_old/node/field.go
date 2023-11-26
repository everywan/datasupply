package node

import (
	"errors"
	"fmt"
	"time"

	"git.in.zhihu.com/antispam/datasupply/dtype"
)

const DefaultTimeout = time.Second * 10

// Field 表示需要补数的字段. 每个 node 会产生至少一个 field, 一个 field 对应输入的一行补数配置.
type Field struct {
	// id 和 code 理论上都是唯一的. id 由系统生成, 无意义, 仅区分唯一性,
	// 而 code 具有业务意义和业务稳定性, 可用于系统间的交互.
	ID            string                 `json:"id"`               // 唯一 ID, 系统生成和使用.
	Code          string                 `json:"code"`             // 用户输入和使用. 一般也是唯一的.
	FieldType     dtype.DType            `json:"field_type"`       // 字段类型
	SupplyStage   SupplyStage            `json:"supply_mode"`      // 补数方式
	FieldOfSupply string                 `json:"field_of_supplu"`  // 取补数结果的哪个值作为字段值.
	OnError       OnErrorHandler         `json:"on_error"`         // 错误处理
	DefaultValue  interface{}            `json:"default_value"`    // 默认值
	Timeout       time.Duration          `json:"timeout"`          // 超时时间, 单位毫秒
	NotExport     bool                   `json:"not_export"`       // 该字段是否存到 dag.result 里
	DelaySupply   time.Duration          `json:"delay_supply"`     // 单位毫秒
	AutoNilToZero bool                   `json:"auto_nil_to_zero"` // 当字段时是 nil 是, 是否自动转为类型零值
	Meta          map[string]interface{} `json:"meta"`             // 用户自己定义的元数据
}

func (field *Field) Validate() error {
	if field.ID == "" {
		return errors.New("field must have id")
	}
	if field.Code == "" {
		return errors.New("field must have code")
	}
	if field.FieldOfSupply == "" {
		return errors.New("field must have field_key")
	}
	if int(field.SupplyStage) > len(SupplyStageNames) {
		return fmt.Errorf("supply_mode validate error: %s", field.SupplyStage)
	}
	if int(field.OnError) > len(OnErrorHandlerNames) {
		return fmt.Errorf("on_error validate error: %s", field.OnError)
	}
	return nil
}

func (field *Field) LoadDefault() {
	if field.ID == "" {
		field.ID = field.Code
	}
	if field.Timeout == 0 {
		field.Timeout = DefaultTimeout
	}
}

func (field *Field) ValueOnError(failReason string) *FieldResult {
	return &FieldResult{
		Meta: FieldMeta{
			FailReason: failReason,
		},
		Value: field.valueOnError(),
	}
}

func (field *Field) valueOnError() interface{} {
	switch field.OnError {
	case OnErrorDefault:
		return field.DefaultValue
	}
	return nil
}
