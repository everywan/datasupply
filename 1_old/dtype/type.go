package dtype

import (
	"encoding/json"
	"errors"
	"fmt"
	"reflect"
	"strconv"
)

// 与外界交互使用易读性更好的 string, 内部计算使用更高效的 uint.
// 通过实现 Marshaler/Stringer 等接口实现内外交换.
type DType reflect.Kind

var _ json.Marshaler = new(DType)
var _ json.Unmarshaler = new(DType)
var _ fmt.Stringer = new(DType)

// var _ fmt.Formatter = new(DType)

const (
	// system default, [0-100)
	Bool    DType = DType(reflect.Bool)
	Int64   DType = DType(reflect.Int64)
	Uint64  DType = DType(reflect.Uint64)
	Float64 DType = DType(reflect.Float64)
	Map     DType = DType(reflect.Map)
	String  DType = DType(reflect.String)

	// datasupply define, [100-200)
	ArrayInt64  DType = iota + 100 // []int64
	ArrayUint64                    // []uint64
	ArrayString                    // []string
	ArrayByte                      // []byte

	// user define, [200,x)
)

var toNames = map[DType]string{
	Bool:        "bool",
	Int64:       "int64",
	Uint64:      "uint64",
	Float64:     "float64",
	Map:         "map",
	String:      "string",
	ArrayInt64:  "[]int64",
	ArrayUint64: "[]uint64",
	ArrayString: "[]string",
	ArrayByte:   "[]byte",
}
var toTypes = map[string]DType{
	"bool":     Bool,
	"int64":    Int64,
	"uint64":   Uint64,
	"float64":  Float64,
	"map":      Map,
	"string":   String,
	"[]int64":  ArrayInt64,
	"[]uint64": ArrayUint64,
	"[]string": ArrayString,
	"[]byte":   ArrayByte,
}

func NewDType(dtype uint, dtypeStr string) (DType, error) {
	if _, ok := toNames[DType(dtype)]; ok {
		return 0, errors.New("dtype has been defined")
	}
	if _, ok := toTypes[dtypeStr]; ok {
		return 0, errors.New("dtype has been defined")
	}
	toNames[DType(dtype)] = dtypeStr
	toTypes[dtypeStr] = DType(dtype)
	return DType(dtype), nil
}

func GetDtype(dtypeStr string) (DType, bool) {
	dtype, isExist := toTypes[dtypeStr]
	return dtype, isExist
}

func (dtype DType) String() string {
	str, ok := toNames[dtype]
	if ok {
		return str
	}
	return "dtype_" + strconv.Itoa(int(dtype))
}

func (dtype DType) MarshalJSON() ([]byte, error) {
	return json.Marshal(dtype.String())
}

func (dtype *DType) UnmarshalJSON(b []byte) error {
	str := ""
	if err := json.Unmarshal(b, &str); err != nil {
		return err
	}
	var ok bool
	*dtype, ok = toTypes[str]
	if !ok {
		return errors.New("unknown dtype name " + str)
	}
	return nil
}
