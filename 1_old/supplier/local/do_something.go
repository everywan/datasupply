package local

import (
	"context"

	"git.in.zhihu.com/antispam/datasupply/supplier"
)

func NewDoSomethingPlugin() supplier.IPlugin {
	// 对外展示, 所以 Name 要大写
	return supplier.NewDefaultPlugin("DoSomething", DoSomething)
}

func DoSomething(_ context.Context, _ ...interface{}) (map[string]interface{}, error) {
	println("do something like print this")
	return nil, nil
}
