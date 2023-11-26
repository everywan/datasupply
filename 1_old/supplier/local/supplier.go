// package local 是本地供应商, 是一系列不需要外部依赖的本地函数.
package local

import (
	"git.in.zhihu.com/antispam/datasupply/supplier"
)

const SupplierName = "local"

func NewSupplier() *supplier.DefaultSupplier {
	localSupplier := supplier.NewDefaultSupplier(SupplierName,
		[]supplier.IPlugin{
			NewDoSomethingPlugin(),
		})

	return localSupplier
}
