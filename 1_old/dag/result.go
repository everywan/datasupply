package dag

import (
	"errors"

	"git.in.zhihu.com/antispam/datasupply/constant"
	"git.in.zhihu.com/antispam/datasupply/node"
)

// dag 对外输出的 Result
//go:generate msgp
type Result struct {
	Fields map[string]*node.FieldResult `json:"fields"`
}

func NewResult() *Result {
	return &Result{
		Fields: map[string]*node.FieldResult{},
	}
}

func (result *Result) GetFieldMeta(fieldName string) (node.FieldMeta, error) {
	if result == nil || result.Fields == nil {
		return node.FieldMeta{}, constant.NullPointerError
	}
	field, ok := result.Fields[fieldName]
	if !ok {
		return node.FieldMeta{}, constant.NotFoundError
	}
	return field.Meta, nil
}

func (result *Result) GetFieldValue(fieldName string) (interface{}, error) {
	if result == nil || result.Fields == nil {
		return nil, constant.NullPointerError
	}
	field, ok := result.Fields[fieldName]
	if !ok {
		return nil, constant.NotFoundError
	}
	if !field.IsSupplySuccess() {
		return nil, errors.New("supply error " + field.Meta.GetFailReason())
	}
	return field.Value, nil
}

func (result *Result) GetFieldValues() map[string]interface{} {
	if result == nil || result.Fields == nil {
		return map[string]interface{}{}
	}

	values := make(map[string]interface{}, len(result.Fields))
	for fieldName, fieldData := range result.Fields {
		if fieldData.IsSupplySuccess() {
			values[fieldName] = fieldData.Value
		}
	}
	return values
}

func (result *Result) Clone() *Result {
	if result == nil || result.Fields == nil {
		return &Result{}
	}

	fields := make(map[string]*node.FieldResult, len(result.Fields))
	for field, value := range result.Fields {
		fields[field] = value.Clone()
	}
	return &Result{
		Fields: fields,
	}
}
