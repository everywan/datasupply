package local

import (
	"context"
	"errors"

	"git.in.zhihu.com/antispam/datasupply/supplier"
)

const ForwardPlugin = "Forward"

func NewForwardPlugin() supplier.IPlugin {
	return supplier.NewDefaultPlugin(ForwardPlugin, Forward)
}

func Forward(_ context.Context, params ...interface{}) (map[string]interface{}, error) {
	if len(params) < 2 {
		return map[string]interface{}{},
			errors.New("forward func at least two param: keys([]string),values(...interface{})")
	}
	keys, ok := params[0].([]string)
	if !ok {
		return map[string]interface{}{},
			errors.New("forward func first param must be keys([]string)")
	}
	if len(keys) > len(params)-1 {
		return map[string]interface{}{},
			errors.New("forward func have not enough params")
	}
	result := make(map[string]interface{}, len(keys))
	for i, key := range keys {
		result[key] = params[i+1]
	}
	return result, nil
}
