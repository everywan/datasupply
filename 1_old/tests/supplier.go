package tests

import (
	"context"
	"errors"

	"git.in.zhihu.com/antispam/datasupply/supplier"
	"git.in.zhihu.com/antispam/datasupply/supplier/local"
)

var DefaultSupplier = NewTestSupplier()

func NewTestSupplier() supplier.ISupplier {
	return supplier.NewDefaultSupplier(
		"supplier_tests",
		[]supplier.IPlugin{
			local.NewDoSomethingPlugin(),
		},
	)
}

func NewTestPlugin(name string, params, fields []string) supplier.IPlugin {
	return supplier.NewDefaultPlugin(
		name,
		func(ctx context.Context, args ...interface{}) (map[string]interface{}, error) {
			if len(args) != len(params) {
				return map[string]interface{}{}, errors.New("param len error")
			}
			result := make(map[string]interface{}, len(args))
			for _, field := range fields {
				result[field] = "x"
			}
			return result, nil
		})
}
