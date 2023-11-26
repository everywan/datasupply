package view

import (
	"sort"

	"git.in.zhihu.com/antispam/datasupply/node"
	"github.com/go-echarts/go-echarts/v2/charts"
	"github.com/go-echarts/go-echarts/v2/opts"
)

func DispalyDAGEchart(root node.INode) *charts.Tree {
	tree := []opts.TreeData{
		{
			Name:     "Root",
			Children: echartTree(root.GetNexts()),
		},
	}
	return GetDAGEchartTree(tree)
}

func echartTree(currentLevels []node.INode) []*opts.TreeData {
	currentLevels = node.Deduplication(currentLevels)
	sort.Slice(currentLevels, func(i, j int) bool {
		return currentLevels[i].GetID() < currentLevels[j].GetID()
	})

	treeCurrent := make([]*opts.TreeData, 0, len(currentLevels))
	for _, cnode := range currentLevels {
		treeCurrent = append(treeCurrent, &opts.TreeData{
			Name:     cnode.GetID(),
			Children: echartTree(cnode.GetNexts()),
		})
	}

	return treeCurrent
}

func GetDAGEchartTree(tree []opts.TreeData) *charts.Tree {
	graph := charts.NewTree()
	graph.SetGlobalOptions(
		charts.WithInitializationOpts(opts.Initialization{Width: "100%", Height: "95vh"}),
		charts.WithTitleOpts(opts.Title{Title: "basic tree example"}),
		charts.WithTooltipOpts(opts.Tooltip{Show: false}),
	)
	graph.AddSeries("tree", tree).
		SetSeriesOptions(
			charts.WithTreeOpts(
				opts.TreeChart{
					Layout:           "orthogonal",
					Orient:           "LR",
					InitialTreeDepth: -1,
					Leaves: &opts.TreeLeaves{
						Label: &opts.Label{Show: true, Position: "right", Color: "Black"},
					},
				},
			),
			charts.WithLabelOpts(opts.Label{Show: true, Position: "top", Color: "Black"}),
		)
	return graph
}
