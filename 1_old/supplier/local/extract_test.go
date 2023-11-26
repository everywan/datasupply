package local

import (
	"context"
	"encoding/json"
	"testing"

	"git.in.zhihu.com/antispam/datasupply/dtype"
	"github.com/stretchr/testify/assert"
)

func TestExtract(t *testing.T) {
	testCases := []struct {
		name        string
		key         string
		value       interface{}
		expectValue interface{}
		dtype       dtype.DType
	}{
		{"f_string", "f_string", "f_string_example", "f_string_example", dtype.String},
		{"f_int64", "f_int64", 100, int64(100), dtype.Int64},
		{"f_int64_long", "f_int64_long", 1568280116236611582, int64(1568280116236611582), dtype.Int64},
		{"f_float64", "f_float64", "1603767488868.01", float64(1603767488868.01), dtype.Float64},
		{"f_float64_long",
			"f_float64_long", 1603767488868831232.01,
			float64(1603767488868831232.01), dtype.Float64}, // float64 因为精度问题会丢数据
		{"f_bool", "f_bool", true, true, dtype.Bool},
		{"f_array_int", "f_array_int", []int{1, 2, 3}, []int64{1, 2, 3}, dtype.ArrayInt64},
		{"f_array_int_str", "f_array_int_str", "[1,2,3]", []int64{1, 2, 3}, dtype.ArrayInt64}, // 兼容现有系统
		{"f_array_int_long", "f_array_int_long",
			[]int{1234567890123456780, 2234567890123456780, 3234567890123456780},
			[]int64{1234567890123456780, 2234567890123456780, 3234567890123456780}, dtype.ArrayInt64},
		{"f_array_string", "f_array_string", []string{"1", "2", "3"}, []string{"1", "2", "3"}, dtype.ArrayString},
	}

	input := make(map[string]interface{}, len(testCases))
	for _, tcase := range testCases {
		input[tcase.key] = tcase.value
	}
	payload, _ := json.Marshal(input)

	t.Parallel()
	t.Run("multi_extract", func(t *testing.T) {
		params := make([]interface{}, 1, len(input)+1)
		params[0] = payload
		for key := range input {
			params = append(params, key)
		}
		result, err := Extract(context.Background(), params...)
		assert.NoError(t, err)
		for _, tcase := range testCases {
			actualValue, err := dtype.Convert(result[tcase.key], tcase.dtype)
			assert.NoError(t, err)
			assert.Equal(t, tcase.expectValue, actualValue)
		}
	})
	t.Run("single_test", func(t *testing.T) {
		for _, tcase := range testCases {
			params := []interface{}{payload, tcase.key}
			result, err := Extract(context.Background(), params...)
			assert.NoError(t, err)
			actualValue, err := dtype.Convert(result[tcase.key], tcase.dtype)
			assert.NoError(t, err)
			assert.Equal(t, tcase.expectValue, actualValue)
		}
	})
}
