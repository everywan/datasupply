package node

import (
	"context"
	"errors"
	"testing"

	"git.in.zhihu.com/antispam/datasupply/dtype"
	"git.in.zhihu.com/antispam/datasupply/tests"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type NodeTestSuite struct {
	suite.Suite
}

func (suite *NodeTestSuite) SetupTest() {}

func (suite *NodeTestSuite) TestRun() {
	param2, _ := NewVariableParam(&CreateVarParamRequest{
		ParamName:    "node_param2_name",
		DagFieldName: "parent_out_field1",
		ParamType:    dtype.String,
		OnError:      ParamOnErrorPrune,
	})
	params := []Param{*NewConstantParam(1, dtype.Int64), *param2}
	params[1].AddValueCheckFns(func(i interface{}) error {
		if i == "" {
			return errors.New("can not be empty")
		}
		return nil
	})
	fields := []*Field{
		{
			Code:          "out_field1",
			FieldOfSupply: "field_of_supply1",
			SupplyStage:   SupplyStageSync,
			FieldType:     dtype.String,
			OnError:       OnErrorDiscard,
		},
	}
	supplier := tests.NewTestSupplier()
	supplier.RegisterPlugin(
		tests.NewTestPlugin("DoSomething",
			[]string{params[0].FieldName, params[1].FieldName},
			[]string{fields[0].FieldOfSupply},
		),
	)
	node, err := New(&CreateNodeRequest{
		FuncName: "DoSomething",
		Params:   params,
		Supplier: supplier,
		Fields:   fields,
		Logger:   tests.DefaultLogger,
	})
	assert.NoError(suite.T(), err, "create node error")

	suite.Run("normal", func() {
		result := node.Run(context.Background(), map[string]interface{}{
			params[1].FieldName: "case1",
		})
		assert.Equal(suite.T(), len(fields), len(result))
		for _, field := range fields {
			fieldResult, ok := result[field.Code]
			assert.True(suite.T(), ok)
			assert.Equal(suite.T(), "x", fieldResult.Value)
		}
	})

	// suite.Run("node_id_check", func() {
	// 	assert.Equal(suite.T(), "supplier_tests_DoSomething_", node.GetID())
	// })

	suite.Run("value_check_error", func() {
		result := node.Run(context.Background(), map[string]interface{}{
			params[1].FieldName: "",
		})
		assert.Equal(suite.T(), len(fields), len(result))
		for _, field := range fields {
			fieldResult, ok := result[field.Code]
			assert.True(suite.T(), ok)
			assert.Equal(suite.T(), nil, fieldResult.Value)
			assert.NotEmpty(suite.T(), fieldResult.Meta.FailReason)
		}
	})
}

func TestNode(t *testing.T) {
	suite.Run(t, new(NodeTestSuite))
}
