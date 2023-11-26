package node

import (
	"errors"
	"fmt"
	"strings"

	"git.in.zhihu.com/antispam/datasupply/log"
	"git.in.zhihu.com/antispam/datasupply/supplier"
)

const SupplierFunc = "supplier"

type CreateNodeRequest struct {
	FuncName    string             `json:"func_name"`
	Params      []Param            `json:"params"`
	Supplier    supplier.ISupplier `json:"supplier"`
	Fields      []*Field           `json:"fields"`
	Middlewares []interface{}
	Logger      log.ILog
}

func (request *CreateNodeRequest) LoadDefault() {
	for i := range request.Fields {
		request.Fields[i].LoadDefault()
	}
	for i := range request.Params {
		request.Params[i].LoadDefault()
	}
	if request.Logger == nil {
		request.Logger = log.NewDefaultLog()
	}
}

// nodeid == supply_name + func_name + supplier_params
func (req *CreateNodeRequest) GenNodeID() string {
	builder := strings.Builder{}
	builder.WriteString(req.Supplier.GetName())
	builder.WriteString("_")
	builder.WriteString(req.FuncName)
	builder.WriteString("_")
	for _, param := range req.Params {
		builder.WriteString(param.ID)
		builder.WriteString("_")
	}
	id := strings.TrimRight(builder.String(), "_")
	return id
}

func (req *CreateNodeRequest) Validate() error {
	if req.Supplier == nil {
		return errors.New("must have supplier")
	}

	// check fields
	fieldIDSet := make(map[string]struct{}, len(req.Fields))
	for _, field := range req.Fields {
		if _, ok := fieldIDSet[field.Code]; ok {
			return fmt.Errorf("field repeat. field: %+v", field)
		}
		fieldIDSet[field.Code] = struct{}{}

		if err := field.Validate(); err != nil {
			return err
		}
	}

	paramIDSet := make(map[string]struct{}, len(req.Params))
	for _, param := range req.Params {
		// 当为变量值时, 不允许重复
		if param.Kind == ParamVariable {
			if _, ok := paramIDSet[param.ID]; ok {
				return fmt.Errorf("param repeat. field: %+v", param)
			}
			paramIDSet[param.ID] = struct{}{}
		}

		if err := param.Validate(); err != nil {
			return err
		}
	}

	return nil
}
