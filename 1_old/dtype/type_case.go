package dtype

import (
	"encoding/json"
	"errors"
	"fmt"
	"go/types"
	"reflect"
	"strconv"
	"strings"
	"unicode"

	"github.com/tidwall/gjson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func ToInt(value interface{}) (int, error) {
	switch value := value.(type) {
	case int:
		return value, nil
	case int16:
		return int(value), nil
	case int32:
		x := int(value)
		if int32(x) != value {
			return 0, strconv.ErrRange
		}
		return x, nil
	case int64:
		x := int(value)
		if int64(x) != value {
			return 0, strconv.ErrRange
		}
		return x, nil
	case uint16:
		return int(value), nil
	case uint32:
		return int(value), nil
	case uint64:
		return int(value), nil
	case float32:
		return int(value), nil
	case float64:
		return int(value), nil
	case string:
		n, err := parseEToInt64(value)
		return int(n), err
	case []byte:
		n, err := parseEToInt64(ToString(value))
		return int(n), err
	case json.Number:
		n, err := value.Int64()
		return int(n), err
	case nil:
		return 0, nil
	}
	return 0, fmt.Errorf("unexpected type for Int, got type %T", value)
}

func ToIntPointer(value interface{}) (*int, error) {
	if reflect.TypeOf(value).Kind() == reflect.Ptr {
		value = reflect.ValueOf(value).Elem().Interface()
	}

	v, err := ToInt(value)
	if err != nil {
		return nil, err
	}
	return &v, nil
}

func ToInt32(value interface{}) (int32, error) {
	switch value := value.(type) {
	case int:
		return int32(value), nil
	case int16:
		return int32(value), nil
	case int32:
		return value, nil
	case int64:
		x := int32(value)
		if int64(x) != value {
			return 0, strconv.ErrRange
		}
		return x, nil
	case uint16:
		return int32(value), nil
	case uint32:
		return int32(value), nil
	case uint64:
		return int32(value), nil
	case float32:
		return int32(value), nil
	case float64:
		return int32(value), nil
	case string:
		n, err := parseEToInt64(value)
		return int32(n), err
	case []byte:
		n, err := parseEToInt64(ToString(value))
		return int32(n), err
	case json.Number:
		n, err := value.Int64()
		return int32(n), err
	case nil:
		return 0, nil
	}
	return 0, fmt.Errorf("unexpected type for Int, got type %T", value)
}

func ToInt32Pointer(value interface{}) (*int32, error) {
	if reflect.TypeOf(value).Kind() == reflect.Ptr {
		value = reflect.ValueOf(value).Elem().Interface()
	}

	v, err := ToInt32(value)
	if err != nil {
		return nil, err
	}
	return &v, nil
}

func ToInt64(value interface{}) (int64, error) {
	switch value := value.(type) {
	case int:
		return int64(value), nil
	case int16:
		return int64(value), nil
	case int32:
		return int64(value), nil
	case int64:
		return value, nil
	case uint16:
		return int64(value), nil
	case uint32:
		return int64(value), nil
	case uint64:
		return int64(value), nil
	case float32:
		return int64(value), nil
	case float64:
		return int64(value), nil
	case string:
		return parseEToInt64(value)
	case []byte:
		return parseEToInt64(string(value))
	case json.Number:
		return value.Int64()
	case gjson.Result:
		return value.Int(), nil
	case nil:
		return 0, nil
	}
	return 0, fmt.Errorf("unexpected type for ToInt64(, got type %T", value)
}

func ToInt64Pointer(value interface{}) (*int64, error) {
	if reflect.TypeOf(value).Kind() == reflect.Ptr {
		value = reflect.ValueOf(value).Elem().Interface()
	}

	v, err := ToInt64(value)
	if err != nil {
		return nil, err
	}
	return &v, nil
}

// 处理科学计数 1.06059501e+08
func parseEToInt64(value string) (int64, error) {
	if strings.Contains(value, "e+") || strings.Contains(value, "e-") {
		f, err := ToFloat64(value)
		return int64(f), err
	}
	n, err := strconv.ParseInt(value, 10, 64)
	return n, err
}

func ToFloat32(value interface{}) (float32, error) {
	switch value := value.(type) {
	case int:
		return float32(value), nil
	case int16:
		return float32(value), nil
	case int32:
		return float32(value), nil
	case int64:
		return float32(value), nil
	case float32:
		return value, nil
	case float64:
		x := float32(value)
		if float64(x) != value {
			return 0, strconv.ErrRange
		}
		return x, nil
	case string:
		n, err := strconv.ParseFloat(value, 32)
		return float32(n), err
	case []byte:
		n, err := strconv.ParseFloat(string(value), 32)
		return float32(n), err
	case json.Number:
		n, err := value.Float64()
		return float32(n), err
	case nil:
		return 0, nil
	}
	return 0, fmt.Errorf("unexpected type for Float32, got type %T", value)
}

func ToFloat32Pointer(value interface{}) (*float32, error) {
	if reflect.TypeOf(value).Kind() == reflect.Ptr {
		value = reflect.ValueOf(value).Elem().Interface()
	}

	v, err := ToFloat32(value)
	if err != nil {
		return nil, err
	}
	return &v, nil
}

func ToFloat64(value interface{}) (float64, error) {
	switch value := value.(type) {
	case int:
		return float64(value), nil
	case int16:
		return float64(value), nil
	case int32:
		return float64(value), nil
	case int64:
		return float64(value), nil
	case float32:
		return float64(value), nil
	case float64:
		return value, nil
	case string:
		n, err := strconv.ParseFloat(value, 64)
		return n, err
	case []byte:
		n, err := strconv.ParseFloat(string(value), 64)
		return n, err
	case json.Number:
		return value.Float64()
	case gjson.Result:
		return value.Float(), nil
	case nil:
		return 0, nil
	}
	return 0, fmt.Errorf("unexpected type for Float64, got type %T", value)
}

func ToFloat64Pointer(value interface{}) (*float64, error) {
	if reflect.TypeOf(value).Kind() == reflect.Ptr {
		value = reflect.ValueOf(value).Elem().Interface()
	}

	v, err := ToFloat64(value)
	if err != nil {
		return nil, err
	}
	return &v, nil
}

func ToString(value interface{}) string {
	switch value := value.(type) {
	case int:
		return strconv.Itoa(value)
	case int16:
		return strconv.FormatInt(int64(value), 10)
	case int32:
		return strconv.FormatInt(int64(value), 10)
	case int64:
		return strconv.FormatInt(value, 10)
	case float32:
		return strconv.FormatFloat(float64(value), 'f', -1, 32)
	case float64:
		return strconv.FormatFloat(value, 'f', -1, 64)
	case bool:
		return strconv.FormatBool(value)
	case string:
		return value
	case json.Number:
		return value.String()
	case []byte:
		return string(value)
	case []interface{}, []int64, []float64, []string:
		v, _ := json.Marshal(value)
		return string(v)
	case gjson.Result:
		return value.String()
	case nil:
		return ""
	}
	return fmt.Sprintf("%v", value)
}

func ToStringPointer(value interface{}) *string {
	if reflect.TypeOf(value).Kind() == reflect.Ptr {
		value = reflect.ValueOf(value).Elem().Interface()
	}

	v := ToString(value)
	return &v
}

func ToBytes(value interface{}) []byte {
	switch value := value.(type) {
	case int, int16, int32, int64, float32, float64:
		return []byte(ToString(value))
	case string:
		return []byte(value)
	case []byte:
		return value
	case nil:
		return nil
	case json.Number:
		return []byte(value)
	case bool:
		return []byte(strconv.FormatBool(value))
	case types.Nil:
		return nil
	}
	return []byte(ToString(value))
}

func ToBool(value interface{}) (bool, error) {
	switch value := value.(type) {
	case int, int32, int64, float32, float64:
		return value != 0, nil
	case string:
		return strconv.ParseBool(value)
	case []byte:
		return strconv.ParseBool(string(value))
	case bool:
		return value, nil
	case gjson.Result:
		return value.Bool(), nil
	case nil:
		return false, nil
	}
	return false, fmt.Errorf("unexpected type for Bool, got type %T", value)
}

func ToBoolPointer(value interface{}) (*bool, error) {
	if reflect.TypeOf(value).Kind() == reflect.Ptr {
		value = reflect.ValueOf(value).Elem().Interface()
	}

	v, err := ToBool(value)
	if err != nil {
		return nil, err
	}
	return &v, nil
}

func ToStringArray(value interface{}) ([]string, error) {
	switch value := value.(type) {
	case []string:
		return value, nil
	case string:
		var r []string
		err := json.Unmarshal([]byte(value), &r)
		return r, err
	case gjson.Result:
		if value.IsArray() {
			values := value.Array()
			result := make([]string, len(values))
			for i, v := range values {
				result[i] = v.String()
			}
			return result, nil
		}
		if value.Type == gjson.String {
			return ToStringArray(value.String())
		}
	case nil:
		return []string{}, nil
	}

	var result []string
	err := sliceHelper(value, func(n int) { result = make([]string, n) }, func(i int, v interface{}) error {
		result[i] = ToString(v)
		return nil
	})
	return result, err
}

func ToArray(value interface{}) ([]interface{}, error) {
	switch value := value.(type) {
	case []interface{}:
		return value, nil
	case string:
		var r []interface{}
		err := json.Unmarshal([]byte(value), &r)
		return r, err
	case nil:
		return []interface{}{}, nil
	}
	return nil, fmt.Errorf("unexpected type for Arrays, got type %T", value)
}

func ToByteArray(value interface{}) ([][]byte, error) {
	var result [][]byte
	err := sliceHelper(value, func(n int) { result = make([][]byte, n) }, func(i int, v interface{}) error {
		p, ok := v.([]byte)
		if !ok {
			return fmt.Errorf("unexpected element type for Byte Slices, got type %T", v)
		}
		result[i] = p
		return nil
	})
	return result, err
}

func ToIntArray(value interface{}) ([]int, error) {
	switch _value := value.(type) {
	case []int:
		return _value, nil
	case string:
		var r []int
		err := json.Unmarshal([]byte(_value), &r)
		return r, err
	case []interface{}:
		r := make([]int, len(_value))
		for i, v := range _value {
			_v, err := ToInt(v)
			if err != nil {
				return nil, fmt.Errorf("cannot case %T to int array", value)
			}
			r[i] = _v
		}
		return r, nil
	case nil:
		return []int{}, nil
	default:
		i, ok := value.([]int)
		if ok {
			return i, nil
		}
		return nil, fmt.Errorf("case to int array error get type %T", value)
	}
}

func ToInt64Array(value interface{}) ([]int64, error) {
	switch _value := value.(type) {
	case []int64:
		return _value, nil
	case string:
		var i []int64
		err := json.Unmarshal([]byte(_value), &i)
		return i, err
	case nil:
		return []int64{}, nil
	case []interface{}:
		return interfaceArrayToInt64Array(_value)
	case primitive.A:
		return interfaceArrayToInt64Array(_value)
	case gjson.Result:
		if _value.IsArray() {
			values := _value.Array()
			result := make([]int64, len(values))
			for i, v := range values {
				result[i] = v.Int()
			}
			return result, nil
		}
		if _value.Type == gjson.String {
			return ToInt64Array(_value.String())
		}
	}
	i, ok := value.([]int64)
	if ok {
		return i, nil
	}
	return nil, fmt.Errorf("case to int64 array error get %T", value)
}

func interfaceArrayToInt64Array(value []interface{}) ([]int64, error) {
	r := make([]int64, len(value))
	for i, v := range value {
		_v, err := ToInt64(v)
		if err != nil {
			return nil, fmt.Errorf("cannot case %T to int64 array", value)
		}
		r[i] = _v
	}
	return r, nil
}

func ToFloat64Array(value interface{}) ([]float64, error) {
	switch _value := value.(type) {
	case []float64:
		return _value, nil
	case string:
		var f []float64
		err := json.Unmarshal([]byte(_value), &f)
		return f, err
	case nil:
		return []float64{}, nil
	case []interface{}:
		r := make([]float64, len(_value))
		for i, v := range _value {
			_v, err := ToFloat64(v)
			if err != nil {
				return nil, fmt.Errorf("cannot case %T to float64 array", value)
			}
			r[i] = _v
		}
		return r, nil
	default:
		i, ok := value.([]float64)
		if ok {
			return i, nil
		}
		return nil, fmt.Errorf("case to int64 array error get %T", value)
	}
}

func ToMap(value interface{}) (map[string]interface{}, error) {
	if value == nil {
		return nil, nil
	}
	m, ok := value.(map[string]interface{})
	if ok {
		return m, nil
	}

	if reflect.TypeOf(value).Name() == "string" {
		var r map[string]interface{}
		if err := json.Unmarshal([]byte(value.(string)), &r); err == nil {
			return r, err
		}
	}

	values, err := ToArray(value)
	if err != nil {
		return nil, err
	}
	if len(values)%2 != 0 {
		return nil, errors.New("map expects even number of values result")
	}
	m = make(map[string]interface{}, len(values)/2)
	for i := 0; i < len(values); i += 2 {
		key, okKey := values[i].([]byte)
		if !okKey {
			return nil, errors.New("map key not a bulk string value")
		}
		m[string(key)] = values[i+1]
	}
	return m, nil
}

func sliceHelper(value interface{}, makeSlice func(int), assign func(int, interface{}) error) error {
	switch value := value.(type) {
	case []interface{}:
		makeSlice(len(value))
		for i := range value {
			if value[i] == nil {
				continue
			}
			if err := assign(i, value[i]); err != nil {
				return err
			}
		}
		return nil
	case nil:
		return nil
	}
	return fmt.Errorf("unexpected type %T", value)
}

func IsNumber(d interface{}) bool {
	s := ToString(d)
	for _, c := range s {
		if !unicode.IsDigit(c) {
			return false
		}
	}
	return true
}
