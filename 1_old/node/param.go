package node

import (
	"errors"
	"fmt"
	"strconv"

	"git.in.zhihu.com/antispam/datasupply/dtype"
	"git.in.zhihu.com/antispam/datasupply/internal/deepcopy"
)

// todo [optimize] supplier/middleware 只需要关注 funcParam 就可以了,
// 不需要关心 caller, caller 是 node 需要关心的.
// 这里的可优化之处在于, 1. 拆分 funcParam 和 nodeparam 的概念. 2. 更好的 param_id 生成方案.

// 节点运行需要的参数列表. 一个节点内有一个或多个函数, 多个函数可能有同名参数值, 参数值相同但处理方式不同.
//go:generate mockgen -package mock_node -destination ./mock/param.go -source=param.go
type INodeParams interface {
	AddFuncParams(funcName string, funcParams []Param) INodeParams
	// 存储 node 用到的所有函数的参数列表. key: func, value: func.params
	GetParamMap() map[string][]Param
	// 拉平后的参数列表
	GetParams() []Param
	// 拉平后的变量参数列表
	GetParamVariables() []Param
	GetParamsByFunc(funcName string) []Param
}

type NodeParams struct {
	paramsMap map[string][]Param
	params    []Param
	paramsVar []Param
}

var _ INodeParams = new(NodeParams)

func NewNodeParams() *NodeParams {
	return &NodeParams{
		paramsMap: map[string][]Param{},
		params:    []Param{},
		paramsVar: []Param{},
	}
}

func (params *NodeParams) AddFuncParams(funcName string, _funcParams []Param) INodeParams {
	funcParams := make([]Param, len(_funcParams))
	for i := range _funcParams {
		funcParams[i] = deepcopy.Copy(_funcParams[i]).(Param)
		funcParams[i].ID = funcName + "_" + funcParams[i].ID
	}
	params.paramsMap[funcName] = funcParams

	paramsMap := make(map[string]struct{}, len(params.params))
	for _, param := range params.params {
		paramsMap[param.ID] = struct{}{}
	}

	// 修改关联数据
	for _, param := range funcParams {
		if _, ok := paramsMap[param.ID]; !ok {
			params.params = append(params.params, param)
			if param.Kind == ParamVariable {
				params.paramsVar = append(params.paramsVar, param)
			}
		}
	}

	return params
}

func (params *NodeParams) GetParams() []Param {
	return params.params
}

func (params *NodeParams) GetParamVariables() []Param {
	return params.paramsVar
}

func (params *NodeParams) GetParamMap() map[string][]Param {
	return params.paramsMap
}

func (params *NodeParams) GetParamsByFunc(funcName string) []Param {
	return params.paramsMap[funcName]
}

// --------------------- 函数参数 ------------------------

// 函数运行需要的参数, 约定函数签名为 fn(ctx, params...), 不存在同名入参.
// 函数参数以数组的形式存储, 所以不需要存储参数的名称, 只关心值.
type Param struct {
	// 作为 nodeParam 和 funcParam 时, id 会发生变化.
	ID        string      `json:"id"`
	Kind      ParamKind   `json:"kind"`  // constant, var
	Value     interface{} `json:"value"` // 常量时有值, 变量时为 nil
	ValueType dtype.DType `json:"value_type"`

	// 变量时才有的分支
	FieldName string `json:"field_name"` // 变量时, 取 dag.result.field 作为参数值
	// TODO [optimize] 可以看下这个 onerror 有无更好的处理方式. 不应该放在 node 上, 因为不同的下游可能有不同的处理.
	OnError ParamOnErrorHandler `json:"on_error"`

	// 验证函数
	ValueCheckFns []func(value interface{}) error `json:"-"`
}

func NewConstantParam(value interface{}, valueType dtype.DType) *Param {
	return &Param{
		ID:        fmt.Sprintf("const_%s_%v", valueType, value),
		Kind:      ParamConstant,
		Value:     value,
		ValueType: valueType,
	}
}

// paramName 在函数中参数的名字. fieldName 在 dag 中寻找那个字段作为参数名.
type CreateVarParamRequest struct {
	ParamName    string              // 参数的名称, 用于确定参数在当前函数的唯一id.
	DagFieldName string              // 参数的值, 将 dag.result.field 作为参数值
	ParamType    dtype.DType         // 参数类型
	OnError      ParamOnErrorHandler // 参数值错误时的处理方式
}

func (req *CreateVarParamRequest) Validate() error {
	if req.ParamName == "" {
		return errors.New("param_name can not be empty")
	}
	if req.DagFieldName == "" {
		return errors.New("dag_field_name can not be empty")
	}
	if _, exist := dtype.GetDtype(req.ParamType.String()); !exist {
		return errors.New("param_type not defined")
	}
	if int(req.OnError) >= len(ParamOnErrorHandlerNames) {
		return errors.New("on_error not found")
	}
	return nil
}

func NewVariableParam(request *CreateVarParamRequest) (*Param, error) {
	if err := request.Validate(); err != nil {
		return &Param{}, err
	}
	return &Param{
		ID:        fmt.Sprintf("var_%s_%s", request.ParamName, request.DagFieldName),
		Kind:      ParamVariable,
		ValueType: request.ParamType,
		FieldName: request.DagFieldName,
		OnError:   request.OnError,
	}, nil
}

func (param *Param) LoadDefault() {}

func (param *Param) Validate() error {
	return nil
}

func (param *Param) AddValueCheckFns(fn func(interface{}) error) {
	// 常量不需要验证值
	if param.Kind == ParamConstant {
		return
	}
	if param.ValueCheckFns == nil {
		param.ValueCheckFns = make([]func(interface{}) error, 0, 1)
	}
	param.ValueCheckFns = append(param.ValueCheckFns, fn)
}

func (param *Param) ValueCheck(v interface{}) error {
	for _, fn := range param.ValueCheckFns {
		if err := fn(v); err != nil {
			return err
		}
	}
	return nil
}

func (param *Param) HandleError() (isPrune bool, paramValue interface{}) {
	switch param.OnError {
	case ParamOnErrorPrune:
		return true, nil
	// case ParamOnErrorDefault:
	default:
		return false, nil
	}
}

type ParamKind int

const (
	ParamConstant ParamKind = iota
	ParamVariable
)

var ParamKindNames = []string{
	ParamConstant: "constant",
	ParamVariable: "variable",
}

func (s ParamKind) String() string {
	if int(s) < len(ParamKindNames) {
		return ParamKindNames[s]
	}
	return "param_kind_" + strconv.Itoa(int(s))
}

type ParamOnErrorHandler int

const (
	// 参数错误时, 进行剪枝, 在上层处理.
	ParamOnErrorPrune ParamOnErrorHandler = iota
	// 参数错误时, 跳过该 caller 执行, 返回 caller 默认值. 默认值逻辑如下.
	// if(caller==midware){return true}, if(caller==supply){return field.default()}
	ParamOnErrorSkip
)

var ParamOnErrorHandlerNames = []string{
	ParamOnErrorPrune: "prune",
	ParamOnErrorSkip:  "skip",
}

func (s ParamOnErrorHandler) String() string {
	if int(s) < len(ParamOnErrorHandlerNames) {
		return ParamOnErrorHandlerNames[s]
	}
	return "param_on_error_handler_" + strconv.Itoa(int(s))
}
