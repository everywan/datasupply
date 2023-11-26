package dtype

import (
	"fmt"
)

// todo [optimize] 可以看下怎么结合 reflect.ConvertibleTo/Convert 实现.
// if rtype.ConvertibleTo(dstType) {
// 	rvalue := reflect.ValueOf(value)
// 	return rvalue.Convert(dstType).Interface(), nil
// }
// if rtype.Kind() == reflect.String {
// 	return strconv.ParseInt(value.(string), 10, 64)
// }

// 将当前值转换为 dtype 类型的值. nil 会被转换为相应类型的零值.
func Convert(value interface{}, dtype DType) (interface{}, error) {
	// rtype := reflect.TypeOf(value)
	// if rtype == nil {
	// 	return nil, errors.New("cannot convert input, type nil")
	// }
	switch dtype {
	case String:
		return ToString(value), nil
	case Int64:
		return ToInt64(value)
	case Float64:
		return ToFloat64(value)
	case Bool:
		return ToBool(value)
	case Map:
		return ToMap(value)
	case ArrayInt64:
		return ToInt64Array(value)
	// case ArrayUint64:
	case ArrayString:
		return ToStringArray(value)
	case ArrayByte:
		v, ok := value.([]byte)
		if ok {
			return v, nil
		}
	}
	return value, fmt.Errorf(fmt.Sprintf("convert %v to %s error", value, dtype.String()))
}
