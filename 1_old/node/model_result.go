package node

const (
	FieldFailReson_ValueIsNil               = "field_value_is_nil"
	FieldFailReson_NotFoundInSupplyResponse = "field_not_found_in_supply_response"
	FieldFailReson_TypeConvertError         = "type_convert_error"
)

//go:generate msgp
type FieldMeta struct {
	FailReason string `json:"fail_reason"`
}

// func (Meta) Marshal()   {}
// func (Meta) Unmarshal() {}
// func (Meta) String()    {}
func (meta FieldMeta) GetFailReason() string {
	return meta.FailReason
}

// 每一个 Field 对应一个 Result.
//go:generate msgp
type FieldResult struct {
	Meta  FieldMeta
	Value interface{}
}

type NewFieldResultRequest struct {
	Value      interface{}
	FailReason string `json:"fail_reason"`
}

func NewFieldResult(req *NewFieldResultRequest) *FieldResult {
	return &FieldResult{
		Meta: FieldMeta{
			FailReason: req.FailReason,
		},
		Value: req.Value,
	}
}

// 第二种情况是指 接口失败, 但是字段有配置默认值.
func (result *FieldResult) IsSupplySuccess() bool {
	return result.Meta.GetFailReason() == "" || result.Value != nil
}

func (result *FieldResult) Clone() *FieldResult {
	return &FieldResult{
		Meta:  result.Meta,
		Value: result.Value,
	}
}

// Node 的对外输出, 每一个 Node.Result 对应多个 FieldResult
type Result map[string]*FieldResult
