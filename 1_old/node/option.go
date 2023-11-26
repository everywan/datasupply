package node

type Option func(*Node)

func SetConcurrent(concurrent int) Option {
	return func(node *Node) {
		node.concurrent = concurrent
	}
}

// priority 范围 [1-100]
func SetPriority(priority int) Option {
	return func(node *Node) {
		// 最大值为 100
		if priority > PriorityMax {
			priority = PriorityMax
		} else if priority < PriorityMin {
			priority = PriorityMin
		}
		node.priority = priority
	}
}
