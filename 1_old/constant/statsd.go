package constant

// 指标

// dag 指标
const (
	DAGRunCount      = "dag.%s.run.count"
	DAGRunSpeed      = "dag.%s.run.speed"
	DAGRunStageSpeed = "dag.%s.run.%s.speed"
)

// node 指标
const (
	NodeRunCount         = "node.%s.run.count"     // node_name
	NodeRunSpeed         = "node.%s.run.speed"     // node_name
	NodeRunError         = "node.%s.run.error"     // node_name
	NodeFieldSupplyError = "field.%s.supply.error" // field_name
)
