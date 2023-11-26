package dag

import (
	"errors"

	"git.in.zhihu.com/antispam/datasupply/log"
	"git.in.zhihu.com/antispam/datasupply/node"
)

type CreateDAGRequest struct {
	// 依赖构建
	Root  node.INode
	Nodes []node.INode

	ID             string
	NodeConcurrent int
	Logger         log.ILog
}

func (request *CreateDAGRequest) LoadDefault() {
	if request.NodeConcurrent < 1 {
		request.NodeConcurrent = 1
	}
	if request.Logger == nil {
		request.Logger = log.NewDefaultLog()
	}
}

func (request *CreateDAGRequest) Validate() error {
	if request.ID == "" {
		return errors.New("create_dag_request.dag must have id")
	}
	if request.Root == nil {
		return errors.New("create_dag_request.root can not be nil")
	}
	if request.Nodes == nil {
		return errors.New("create_dag_request.nodes can not be nil")
	}
	for _, node := range request.Nodes {
		if node == nil {
			return errors.New("create_dag_request.nodes can not continue nil node")
		}
	}

	return nil
}
