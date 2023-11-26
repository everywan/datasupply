package local

import (
	"context"
	"errors"

	"git.in.zhihu.com/antispam/datasupply/supplier"
	"github.com/tidwall/gjson"
)

func NewExtractPlugin() supplier.IPlugin {
	return supplier.NewDefaultPlugin("Extract", Extract)
}

func Extract(_ context.Context, params ...interface{}) (map[string]interface{}, error) {
	if len(params) < 2 {
		return map[string]interface{}{},
			errors.New("extract func at least two param: payload([]byte),paths(...string)")
	}
	payload, ok := params[0].([]byte)
	if !ok {
		return map[string]interface{}{},
			errors.New("extract func first param must be payload([]byte)")
	}
	paths := make([]string, len(params)-1)
	for i := 1; i < len(params); i++ {
		paths[i-1] = params[i].(string)
	}

	out := make(map[string]interface{}, len(paths))
	for i, field := range gjson.GetManyBytes(payload, paths...) {
		out[paths[i]] = field
	}
	return out, nil
}
