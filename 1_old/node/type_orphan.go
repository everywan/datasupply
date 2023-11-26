package node

const (
	OrphanReasonNotEnoughParam = "not_enough_params"
	OrphanReasonAncestorPruned = "ancestor_pruned"
)

type OrphanNode struct {
	Node           INode
	Reason         string
	NotFoundParams []string
}

func NewOrphanNode(node INode, reason string, notFoundParams []string) *OrphanNode {
	return &OrphanNode{
		Node:           node,
		Reason:         reason,
		NotFoundParams: notFoundParams,
	}
}
