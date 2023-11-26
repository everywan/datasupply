package view

import (
	"sort"
	"strconv"
	"strings"

	"git.in.zhihu.com/antispam/datasupply/node"
)

func DisplayDAG(root node.INode) {
	displayDAG([]node.INode{root}, 0)
}

func displayDAG(currentLevels []node.INode, depth int) {
	currentLevels = node.Deduplication(currentLevels)
	sort.Slice(currentLevels, func(i, j int) bool {
		return currentLevels[i].GetID() < currentLevels[j].GetID()
	})

	builder := strings.Builder{}
	builder.WriteString("level ")
	builder.WriteString(strconv.Itoa(depth))
	builder.WriteString(": \n\t")

	nextLevels := []node.INode{}
	for _, cnode := range currentLevels {
		builder.WriteString(cnode.GetID())
		builder.WriteString("\n\t")
		nextLevels = append(nextLevels, cnode.GetNexts()...)
	}
	println(strings.TrimRight(builder.String(), "\n\t"))
	if len(nextLevels) != 0 {
		displayDAG(nextLevels, depth+1)
	}
}
