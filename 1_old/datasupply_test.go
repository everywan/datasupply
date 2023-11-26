package datasupply

import (
	"context"
	"testing"

	"git.in.zhihu.com/antispam/datasupply/dtype"
	"git.in.zhihu.com/antispam/datasupply/node"
	"git.in.zhihu.com/antispam/datasupply/supplier"
	"git.in.zhihu.com/antispam/datasupply/tests"
	"github.com/stretchr/testify/assert"
	"golang.org/x/sync/errgroup"
)

// 黑盒测试, 不 mock api, 直接进行测试.
/*
	测试的 dag 格式
	    	    root
			   /	\
	    child1_1	child1_2
		/		\	/		\
	child2_1	child2_2	child2_3
			\			/(child1_2)
				child3_1
*/
func TestDAG(t *testing.T) {
	ds := New()

	supplier := tests.DefaultSupplier
	rootParams := []string{"root_in_1"}
	rootFields := []string{"root_out_1"}
	supplier.RegisterPlugin(tests.NewTestPlugin("root_func", rootParams, rootFields))
	_, err := ds.BuildRoot(genNodeCfg("root_func", rootParams, rootFields, supplier)[0])
	if err != nil {
		t.Fatal("create root error, ", err)
	}

	testCases := []struct {
		funcName string
		params   []string
		fields   []string
	}{
		// 第一层节点
		{"child1_1_func", []string{rootFields[0]}, []string{"child1_1_out_1", "child1_1_out_2"}},
		{"child1_2_func", []string{rootFields[0]}, []string{"child1_2_out_1", "child1_2_out_2"}},
		// 第二层节点
		{"child2_1_func", []string{"child1_1_out_1"}, []string{"child2_1_out_1"}},
		{"child2_2_func", []string{"child1_1_out_2", "child1_2_out_1"}, []string{"child2_2_out_1"}},
		{"child2_3_func", []string{"child1_2_out_2"}, []string{"child2_3_out_1"}},
		// 第三层节点
		{"child3_1_func", []string{"child2_1_out_1", "child1_2_out_1"}, []string{"child3_1_out_1"}},
	}
	for _, tcase := range testCases {
		supplier.RegisterPlugin(tests.NewTestPlugin(tcase.funcName, tcase.params, tcase.fields))
		cfgs := genNodeCfg(tcase.funcName, tcase.params, tcase.fields, supplier)
		for _, cfg := range cfgs {
			_, err := ds.BuildNode(cfg)
			if err != nil {
				t.Fatal("create node error, ", err)
			}
		}
	}
	_, err = ds.BuildDAG(&DAGConfig{
		ID:             "tests",
		NodeConcurrent: 50,
	})
	if err != nil {
		t.Fatal("create dag error, ", err)
	}

	// println("datasupply test dag:")
	// view.DisplayDAG(ds.root)

	errGroup := errgroup.Group{}
	n := 2
	errGroup.SetLimit(n)
	for i := 0; i < n; i++ {
		dag := ds.GetDAG()
		result := dag.Supply(context.TODO(), "test", map[string]interface{}{
			rootParams[0]: "x",
		})
		// fmt.Printf("%s\n", utils.StructToString(result))
		for _, tcase := range testCases {
			for _, field := range tcase.fields {
				if _, ok := result.Fields[field]; !ok {
					t.Fatalf("can not find field [%s] in dag result", field)
				}
			}
		}
		for fieldName, fieldData := range result.Fields {
			failReason := fieldData.Meta.GetFailReason()
			if failReason != "" {
				t.Fatalf("field [%s] error: %s", fieldName, failReason)
			}
			if fieldData.Value != "x" {
				t.Fatalf("field [%s] value error: %s", fieldName, fieldData.Value)
			}
		}
	}
	assert.NoError(t, errGroup.Wait())
}

func genNodeCfg(funcName string, _params, _fields []string, supplier supplier.ISupplier) []*NodeConfig {
	params := make([]node.Param, len(_params))
	for i, param := range _params {
		_param, _ := node.NewVariableParam(&node.CreateVarParamRequest{
			ParamName:    param,
			DagFieldName: param,
			ParamType:    dtype.String,
			OnError:      node.ParamOnErrorPrune,
		})
		params[i] = *_param
	}

	cfgs := make([]*NodeConfig, len(_fields))
	for i, field := range _fields {
		cfgs[i] = &NodeConfig{
			Fields: []*node.Field{
				{
					Code:          field,
					FieldOfSupply: field,
					SupplyStage:   node.SupplyStageSync,
					FieldType:     dtype.String,
					OnError:       node.OnErrorDiscard,
				},
			},
			Supplier: supplier,
			FuncName: funcName,
			Params:   params,
		}
	}
	return cfgs
}
