package node

import "time"

// node attr 的 get/set 方法

func (node *Node) GetID() string {
	return node.id
}

func (node *Node) GetPriority() int {
	return node.priority
}

func (node *Node) GetSupplyStage() SupplyStage {
	return node.supplyStage
}

func (node *Node) GetPrevs() []INode {
	return node.prevs
}

func (node *Node) GetNexts() []INode {
	return node.nexts
}

func (node *Node) GetFields() []*Field {
	return node.fields
}

func (node *Node) GetFieldCodes() []string {
	return node.fieldCodes
}

func (node *Node) GetTimeout() time.Duration {
	return node.timeout
}

func (node *Node) GetDelaySupply() time.Duration {
	return node.delaySupply
}

func (node *Node) AddFields(fields ...*Field) {
	fields = append(fields, node.fields...)
	fieldIDSet := make(map[string]struct{}, len(fields))
	newFields := make([]*Field, 0, len(fields))
	newfieldCodes := make([]string, 0, len(fields))
	for _, field := range fields {
		if _, ok := fieldIDSet[field.ID]; ok {
			continue
		}
		newFields = append(newFields, field)
		newfieldCodes = append(newfieldCodes, field.Code)
	}

	node.fields = newFields
	node.fieldCodes = newfieldCodes
}

func (node *Node) RemoveNexts(nodes ...INode) {
	removeNodeMap := make(map[string]struct{}, len(nodes))
	for _, node := range nodes {
		removeNodeMap[node.GetID()] = struct{}{}
	}
	newNexts := make([]INode, 0, len(node.GetNexts())-len(nodes))
	for _, node := range node.GetNexts() {
		if _, ok := removeNodeMap[node.GetID()]; !ok {
			newNexts = append(newNexts, node)
		}
	}
	node.nexts = newNexts
}

func (node *Node) AddPrevs(nodes ...INode) {
	node.prevs = Deduplication(append(node.prevs, nodes...))
}

func (node *Node) AddNexts(nodes ...INode) {
	node.nexts = Deduplication(append(node.nexts, nodes...))
}

func (node *Node) SetPriority(priority int) {
	node.priority = priority
}

func (node *Node) SetSupplyStage(stage SupplyStage) {
	node.supplyStage = stage
}
