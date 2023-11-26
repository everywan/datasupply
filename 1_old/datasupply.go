// Package Datasupply 补数 SDK
package datasupply

import (
	"git.in.zhihu.com/antispam/datasupply/dag"
	"git.in.zhihu.com/antispam/datasupply/log"
	"git.in.zhihu.com/antispam/datasupply/node"
	"git.in.zhihu.com/antispam/datasupply/supplier"
)

type IDatasupply interface {
	BuildRoot(cfg *NodeConfig, options ...node.Option) (node.INode, error)
	BuildNode(cfg *NodeConfig, options ...node.Option) (node.INode, error)
	BuildDAG(cfg *DAGConfig, options ...dag.Option) (dag.IDAG, error)
	GetDAG() dag.IDAG
}

type Datasupply struct {
	root  node.INode
	nodes []node.INode
	dag   dag.IDAG
}

var _ IDatasupply = new(Datasupply)

func New() *Datasupply {
	return &Datasupply{}
}

func (ds *Datasupply) BuildRoot(cfg *NodeConfig, options ...node.Option) (node.INode, error) {
	var err error
	ds.root, err = ds.buildNode(cfg, options...)
	return ds.root, err
}

type NodeConfig struct {
	Supplier supplier.ISupplier
	FuncName string       `json:"func_name"`
	Params   []node.Param `json:"params"`
	Fields   []*node.Field
	Logger   log.ILog
}

func (ds *Datasupply) BuildNode(cfg *NodeConfig, options ...node.Option) (node.INode, error) {
	newNode, err := ds.buildNode(cfg, options...)
	if err != nil {
		return nil, err
	}
	ds.nodes = append(ds.nodes, newNode)
	return newNode, nil
}

func (ds *Datasupply) buildNode(cfg *NodeConfig, options ...node.Option) (node.INode, error) {
	request := &node.CreateNodeRequest{
		FuncName: cfg.FuncName,
		Params:   cfg.Params,
		Fields:   cfg.Fields,
		Supplier: cfg.Supplier,
		Logger:   cfg.Logger,
	}
	newNode, err := node.New(request, options...)
	if err != nil {
		return nil, err
	}

	return newNode, err
}

type DAGConfig struct {
	ID             string
	NodeConcurrent int
	Logger         log.ILog
}

func (ds *Datasupply) BuildDAG(cfg *DAGConfig, options ...dag.Option) (dag.IDAG, error) {
	dag, err := dag.New(&dag.CreateDAGRequest{
		Root:           ds.root,
		Nodes:          ds.nodes,
		ID:             cfg.ID,
		NodeConcurrent: cfg.NodeConcurrent,
		Logger:         cfg.Logger,
	}, options...)
	if err != nil {
		return nil, err
	}
	ds.dag = dag
	return dag, nil
}

func (ds *Datasupply) GetDAG() dag.IDAG {
	return ds.dag
}
