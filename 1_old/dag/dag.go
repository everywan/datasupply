package dag

import (
	"context"
	"errors"
	"fmt"
	"sync"
	"time"

	"git.in.zhihu.com/antispam/datasupply/dtype"
	"git.in.zhihu.com/antispam/datasupply/log"
	"git.in.zhihu.com/antispam/datasupply/node"
)

// 考虑到部分延迟补数为 120s, 设置 150 更合适些.
var DAGTimtout = time.Second * 150

//go:generate mockgen -package mock_dag -destination ./mock/dag.go -source=dag.go
type IDAG interface {
	GetID() string
	// 运行 DAG, 返回 dag.runtime. 通过 runtime 处理和获取运行时数据.
	Run(ctx context.Context, runtimeID string, paramMap map[string]interface{}) IRuntime
	// 补充所有字段
	Supply(ctx context.Context, runtimeID string, paramMap map[string]interface{}) *Result

	// 获取当前字段依赖的字段
	GetFieldRelys(ctx context.Context, field string) ([]string, error)
	// 按字段进行补数
	SupplyField(ctx context.Context, data map[string]interface{}, field string) *node.FieldResult
	// SupplyFields(ctx context.Context, data map[string]interface{}, fields []string) *Result

	GetRoot() node.INode
	Use(...Middleware)
	Update(nodes []node.INode)

	// 需要随 DAG 更新而更新的数据
	// fieldInfo. 依赖 map[field]Node 的函数
	GetField(ctx context.Context, fieldCode string) (*node.Field, error)
	GetNodeByField(ctx context.Context, field string) (node.INode, error)
	// 输入一个字段和该字段的值, 转换 value 为所需的字段类型.
	ConvertField(fieldCode string, value interface{}) (interface{}, error)

	Close()
}

// DAG is 根据节点的依赖关系构建的图.
type DAG struct {
	id             string
	root           node.INode
	nodeConcurrent int // 字段并发执行数
	logger         log.ILog

	mwchain     Handler // middleware chain
	middlewares []Middleware
	mwChainLock sync.Locker

	// 为减少重复计算, 提前计算必要统计数据
	preComputeData
}

var _ IDAG = new(DAG)

func New(request *CreateDAGRequest, _options ...Option) (*DAG, error) {
	request.LoadDefault()
	if err := request.Validate(); err != nil {
		return &DAG{}, err
	}

	// 获取配置
	options := &options{}
	for _, option := range _options {
		option(options)
	}

	root := request.Root
	nodes := request.Nodes
	{
		dagBuilder := newDagBuilder(request.Logger)
		root, _ = dagBuilder.build(root, nodes)
	}

	preComputeData := preCompute(root)

	dag := &DAG{
		id:             request.ID,
		root:           root,
		nodeConcurrent: request.NodeConcurrent,
		logger:         request.Logger,
		preComputeData: preComputeData,
		middlewares:    []Middleware{},
		mwChainLock:    &sync.Mutex{},
	}
	dag.mwchain = dag.handler

	return dag, nil
}

func (dag *DAG) GetID() string {
	return dag.id
}

// 更新节点配置
func (dag *DAG) Update(nodes []node.INode) {
	// todo 将所有受影响的节点 copy 一份新的, 构建新的 root
	// todo 注意顺序的问题. 如果同一时间有多个更改, 则只能是最新的生效.
	// root := node.NewRootNode()
	// dag.root = root
}

func (dag *DAG) GetField(ctx context.Context, fieldCode string) (*node.Field, error) {
	_, field, err := dag.getFieldInfo(ctx, fieldCode)
	return field, err
}

func (dag *DAG) GetNodeByField(ctx context.Context, field string) (node.INode, error) {
	node, _, err := dag.getFieldInfo(ctx, field)
	return node, err
}

func (dag *DAG) ConvertField(fieldCode string, value interface{}) (interface{}, error) {
	field, err := dag.GetField(context.Background(), fieldCode)
	if err != nil {
		return value, err
	}
	return dtype.Convert(value, field.FieldType)
}

func (dag *DAG) getFieldInfo(ctx context.Context, fieldCode string) (node.INode, *node.Field, error) {
	cnode, ok := dag.field2NodeMap[fieldCode]
	field, ok2 := dag.field2FieldMap[fieldCode]
	if ok && ok2 {
		return cnode, field, nil
	}
	for _, nodeField := range dag.GetRoot().GetFields() {
		if fieldCode == nodeField.Code {
			return dag.GetRoot(), nodeField, nil
		}
	}
	for _, cnode := range dag.GetRoot().Prune() {
		for _, nodeField := range cnode.GetFields() {
			if nodeField.Code == fieldCode {
				return cnode, nodeField, nil
			}
		}
	}
	return nil, nil, errors.New("node not found by field " + fieldCode)
}

func (dag *DAG) Supply(ctx context.Context, runtimeID string, paramMap map[string]interface{}) *Result {
	result := dag.Run(ctx, runtimeID, paramMap).Wait(ctx).GetResultCopy()
	return result
}

func (dag *DAG) Run(ctx context.Context, runtimeID string, paramMap map[string]interface{}) IRuntime {
	return dag.mwchain(ctx, runtimeID, paramMap)
}

// 如果需要设置 traceid 等信息, 可以从改造 context 入手.
// todo [optimize] 这里其实有很多的扩展空间, 目前是同步结束就返回结果, 其实 dag 已经支持任意阶段判断
func (dag *DAG) handler(ctx context.Context, runtimeID string, paramMap map[string]interface{}) IRuntime {
	runtime := dag.createRuntime(runtimeID)
	return runtime.Run(ctx, paramMap)
}

func (dag *DAG) createRuntime(runtimeID string) *runtime {
	// 每次运行都需要重新生成
	nsKeeper := newNodeStateKeeper(int(dag.allNodeCnt))
	stageKeeper := NewDefaultStageKeeper(dag.stageNodeCntMap)
	resultKeeper := NewDefaultResultKeeper()
	return &runtime{
		id:           runtimeID,
		root:         dag.root,
		allNodeCnt:   dag.allNodeCnt,
		concurrent:   dag.nodeConcurrent,
		logger:       dag.logger,
		allNodeDone:  make(chan struct{}),
		nsKeeper:     nsKeeper,
		stageKeeper:  stageKeeper,
		resultKeeper: resultKeeper,
		finishLocker: &sync.Mutex{},
	}
}

func (dag *DAG) GetRoot() node.INode {
	return dag.root
}

// 中间件调用链, 按照 Use 的顺序执行. 每次添加 middleware 需要重新构建.
// 可以添加指针指向 chain.last_handler, 这样每次新增 middleware 时, 替换这个指针.
// 考虑到该操作十分低频, 性能提升获取的收益远小于复杂度提升带来的缺点, 故放弃.
func (dag *DAG) Use(middlewares ...Middleware) {
	dag.mwChainLock.Lock()
	defer dag.mwChainLock.Unlock()

	dag.middlewares = append(dag.middlewares, middlewares...)
	mwchain := dag.handler
	for i := len(dag.middlewares) - 1; i >= 0; i-- {
		mwchain = dag.middlewares[i](mwchain)
	}
	dag.mwchain = mwchain
}

func (dag *DAG) GetFieldRelys(ctx context.Context, fieldCode string) ([]string, error) {
	cnode, _, err := dag.getFieldInfo(ctx, fieldCode)
	if err != nil {
		return []string{}, err
	}

	fieldSet := dag.getFieldRelys(ctx, cnode)
	fields := make([]string, 0, len(fieldSet))
	for field := range fieldSet {
		fields = append(fields, field)
	}
	return fields, nil
}

func (dag *DAG) getFieldRelys(ctx context.Context, cnode node.INode) map[string]struct{} {
	// todo [next] 不允许执行 root 节点, 因为 root 的参数是外部输入的, 存在外部输入和内部字段相同的情况, 会导致无限循环.
	// 后续可以考虑兼容这个问题, 要求外部输入与内部字段不同.
	if cnode.GetID() == dag.GetRoot().GetID() {
		return map[string]struct{}{}
	}
	fieldSet := map[string]struct{}{}
	for _, param := range cnode.GetParamVariables() {
		fieldSet[param.FieldName] = struct{}{}
	}
	for _, nodeParent := range cnode.GetPrevs() {
		for field := range dag.getFieldRelys(ctx, nodeParent) {
			fieldSet[field] = struct{}{}
		}
	}
	return fieldSet
}

func (dag *DAG) SupplyField(ctx context.Context, data map[string]interface{}, field string) *node.FieldResult {
	cnode, err := dag.GetNodeByField(ctx, field)
	if err != nil {
		return node.NewFieldResult(&node.NewFieldResultRequest{FailReason: err.Error()})
	}
	nodeResult := dag.runSpecNode(ctx, cnode, data)
	fieldResult := nodeResult[field]
	return fieldResult
}

func (dag *DAG) runSpecNode(ctx context.Context, cnode node.INode, event map[string]interface{}) node.Result {
	rt := cnode.CreateRuntime()
	for _, param := range cnode.GetParamVariables() {
		paramValue, ok := event[param.FieldName]
		if !ok {
			// 找到输出该参数值的上游节点
			cnodeParent, err := dag.GetNodeByField(ctx, param.FieldName)
			if err != nil {
				return cnode.ValueOnError(fmt.Sprintf("get prama [%s] error [%s]", param.FieldName, err.Error()))
			}
			// todo [next] 不允许执行 root 节点, 因为 root 的参数是外部输入的, 存在外部输入和内部字段相同的情况, 会导致无限循环.
			// 后续可以考虑兼容这个问题, 要求外部输入与内部字段不同.
			if cnodeParent.GetID() == dag.GetRoot().GetID() {
				return cnode.ValueOnError("can not run root node")
			}
			// 执行该节点以获取值
			_nodeResult := dag.runSpecNode(ctx, cnodeParent, event)
			_fieldResult := _nodeResult[param.FieldName]
			if !_fieldResult.IsSupplySuccess() {
				return cnode.ValueOnError(fmt.Sprintf("field [%s] supply error [%s]",
					param.FieldName, _fieldResult.Meta.GetFailReason()))
			}
			paramValue = _fieldResult.Value
		}
		rt.AddParam(param.FieldName, paramValue)
	}
	if !rt.IsReady() {
		return cnode.ValueOnError("params lost")
	}
	nodeResult := rt.Run(ctx)
	return nodeResult
}

func (dag *DAG) Close() {}
