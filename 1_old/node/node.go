/* Package node.
datasupply 有两种节点: root 和 普通节点.

用户输入
字段a1, 从函数 f1 取, 参数 p1p2p3
字段a2, 从函数 f1 取, 参数 p1p2p3
字段a3, 从函数 f2 取, 参数 p1p2

datasupply 会整合输入, 合并函数+参数相同的节点 a1/a2.

输入一系列的补数配置, 每个补数配置可表示为 (字段,函数,参数), 其中 函数+参数 相同的节点会被合并为一个节点,
每个节点有一个或多个字段, 每个字段对应用户输入的初始字段.

node 并发执行是安全的, 并发读改是不安全的.
*/
package node

import (
	"context"
	"fmt"
	"sync"
	"time"

	"git.in.zhihu.com/antispam/datasupply/dtype"
	"git.in.zhihu.com/antispam/datasupply/log"
	"git.in.zhihu.com/antispam/datasupply/supplier"
	"git.in.zhihu.com/antispam/datasupply/utils"
)

// 通用节点
type INode interface {
	INodeParams

	// 属性操作
	GetID() string
	GetPriority() int
	GetSupplyStage() SupplyStage
	GetPrevs() []INode
	GetNexts() []INode
	GetFields() []*Field
	GetFieldCodes() []string
	GetTimeout() time.Duration
	GetDelaySupply() time.Duration

	// 动作
	CreateRuntime() IRuntime
	Use(...IMiddleware)
	Run(ctx context.Context, paramMap map[string]interface{}) Result
	Prune() []INode
	ValueOnPrune() Result
	ValueOnError(failReason string) Result

	// 修改 node 配置. 注意, node 不一定是并发安全的.
	AddFields(...*Field)
	RemoveNexts(...INode)
	AddPrevs(...INode)
	AddNexts(...INode)
	SetPriority(int)
	SetSupplyStage(SupplyStage)
}

// important Node 不是并发安全的.
type Node struct {
	// 自动生成字段
	id string // 节点的唯一标识, 一般是可执行函数的唯一值: `supply_source+func_name+params`

	// 输入字段
	supplier supplier.ISupplier // 定义补数方式如 抽取, 本地函数, 远程函数等.
	funcName string             // 函数名称
	*NodeParams

	// todo [optimize] concurrent 没有用起来.
	// define from option
	concurrent int // 节点并发运行限制
	// 1-100, 当用户未指定时, 根据权重字段自动计算.
	// 权重字段如 field.is_sync,len(fields),len(nexts) 等影响整体吞吐的指标.
	priority int

	// 从 fields 中提取/计算得出的字段
	fields      []*Field      // 节点对外输出的字段, code 重复会覆盖
	fieldCodes  []string      // field.code 集合
	supplyStage SupplyStage   // min(fields.supplyStage)
	timeout     time.Duration // max(fields.timeout)
	delaySupply time.Duration // max(delaySupply)

	// 中间件
	mwchain     Handler
	middlewares []IMiddleware
	mwChainLock sync.Locker
	logger      log.ILog

	prevs []INode
	nexts []INode
}

var _ INode = new(Node)

func New(request *CreateNodeRequest, options ...Option) (*Node, error) {
	request.LoadDefault()
	if err := request.Validate(); err != nil {
		return &Node{}, err
	}

	params := NewNodeParams()
	params.AddFuncParams(SupplierFunc, request.Params)

	var timeout, delaySupply time.Duration
	supplyStage := SupplyStageLazy
	fieldCodes := make([]string, len(request.Fields))
	for i, field := range request.Fields {
		fieldCodes[i] = field.Code
		if field.Timeout > timeout {
			timeout = field.Timeout
		}
		if supplyStage > field.SupplyStage {
			supplyStage = field.SupplyStage
		}
		if time.Duration(field.DelaySupply) > delaySupply {
			delaySupply = field.DelaySupply
		}
	}

	logger := request.Logger
	if logger == nil {
		logger = log.NewDefaultLog()
	}

	node := &Node{
		id:          request.GenNodeID(),
		supplier:    request.Supplier,
		funcName:    request.FuncName,
		NodeParams:  params,
		fields:      request.Fields,
		fieldCodes:  fieldCodes,
		supplyStage: supplyStage,
		timeout:     timeout,
		delaySupply: delaySupply,
		middlewares: []IMiddleware{},
		mwChainLock: &sync.Mutex{},
		logger:      logger,
	}
	node.mwchain = node.handler
	for _, option := range options {
		option(node)
	}
	return node, nil
}

func (node *Node) CreateRuntime() IRuntime {
	paramsLen := len(node.GetParamVariables())
	return &Runtime{
		varParams: map[string]interface{}{},
		node:      node,
		nodeid:    node.GetID(),
		paramCnt:  int64(paramsLen),
		locker:    &sync.Mutex{},
	}
}

func (node *Node) Use(middlewares ...IMiddleware) {
	node.mwChainLock.Lock()
	defer node.mwChainLock.Unlock()

	node.middlewares = append(node.middlewares, middlewares...)
	mwchain := node.handler
	for i := len(node.middlewares) - 1; i >= 0; i-- {
		mw := node.middlewares[i]
		node.AddFuncParams(mw.Name(), mw.Params())
		mwchain = mw.MiddlewareFunc()(mwchain)
	}
	node.mwchain = mwchain
}

func (node *Node) Run(ctx context.Context, paramMap map[string]interface{}) Result {
	ctx, canel := context.WithTimeout(ctx, node.timeout)

	valueCh := make(chan Result, 1)
	utils.SafelyGo(
		func() {
			valueCh <- node.mwchain(ctx, paramMap)
			close(valueCh)
		},
		func(err error) {
			node.logger.Errorf(ctx, "node [%s] panic, error [%v]",
				node.id, err)
			valueCh <- node.ValueOnError("node_run_err:" + fmt.Sprint(err))
			close(valueCh)
		})

	var r Result
	select {
	case <-ctx.Done():
		r = node.ValueOnError("timeout")
	case r = <-valueCh:
	}

	canel()
	return r
}

func (node *Node) handler(ctx context.Context, paramMap map[string]interface{}) Result {
	paramsCfg := node.GetParamsByFunc(SupplierFunc)
	params := make([]interface{}, len(paramsCfg))
	for i, param := range paramsCfg {
		switch param.Kind {
		case ParamConstant:
			params[i] = param.Value
		case ParamVariable:
			params[i] = paramMap[param.FieldName]
			if err := param.ValueCheck(params[i]); err != nil {
				return node.ValueOnError("param_value_check_error: " + err.Error())
			}
		}
	}

	supplyFields, err := node.supplier.Supply(ctx, node.funcName, params)
	if err != nil {
		node.logger.Errorf(ctx, "node [%s] supplier error, func [%s], params [%s], error [%v]",
			node.id, node.funcName, utils.StructToString(params), err)
		return node.ValueOnError("supplier_error: " + err.Error())
	}

	result := make(Result, len(node.fields))
	for _, field := range node.fields {
		fieldCode := field.Code
		fieldValue, ok := supplyFields[field.FieldOfSupply]
		if !ok {
			node.logger.Warnf(ctx, "field [%s] not found in func [%s] response", fieldCode, node.funcName)
			result[fieldCode] = field.ValueOnError(FieldFailReson_NotFoundInSupplyResponse)
			continue
		}
		if fieldValue == nil {
			node.logger.Warnf(ctx, "field [%s] value is nil", fieldCode)
			if !field.AutoNilToZero {
				result[fieldCode] = field.ValueOnError(FieldFailReson_ValueIsNil)
				continue
			}
		}
		if fieldValue, err = dtype.Convert(fieldValue, field.FieldType); err != nil {
			node.logger.Warnf(ctx, "match field [%s][%v] type [%s] error", fieldCode, fieldValue, field.FieldType)
			result[fieldCode] = field.ValueOnError(FieldFailReson_TypeConvertError)
			continue
		}
		result[fieldCode] = &FieldResult{
			Value: fieldValue,
		}
	}

	return result
}

// todo [optimize] 可以优化下性能
func (node *Node) Prune() []INode {
	nextNodes := node.GetNexts()
	nodeMap := make(map[string]INode, len(nextNodes))
	for _, cnode := range nextNodes {
		nodeMap[cnode.GetID()] = cnode
		for _, _cnode := range cnode.Prune() {
			nodeMap[_cnode.GetID()] = _cnode
		}
	}
	pruneNodes := make([]INode, 0, len(nodeMap))
	for _, cnode := range nodeMap {
		pruneNodes = append(pruneNodes, cnode)
	}
	return pruneNodes
}

func (node *Node) ValueOnPrune() Result {
	return node.ValueOnError("prune")
}

func (node *Node) ValueOnError(failReason string) Result {
	result := make(Result, len(node.fields))
	for _, field := range node.fields {
		result[field.Code] = field.ValueOnError(failReason)
	}
	return result
}
