package node

import (
	"errors"
	"reflect"
)

func Deduplication(nodes []INode) []INode {
	m := make(map[string]INode, len(nodes))
	for _, cnode := range nodes {
		m[cnode.GetID()] = cnode
	}
	newNodes := make([]INode, 0, len(m))
	for _, cnode := range m {
		newNodes = append(newNodes, cnode)
	}
	return newNodes
}

// 注意, slice 的空值是 nil, 数组的空值才是 []int{}.
func ParamValueCheckNotZero(value interface{}) error {
	if value == nil {
		return errors.New("is nil")
	}
	if reflect.ValueOf(value).IsZero() {
		return errors.New("is zero")
	}
	return nil
}
