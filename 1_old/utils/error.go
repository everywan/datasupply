package utils

import "encoding/json"

// in 必须是指针
func StructToString(in interface{}) string {
	str, _ := json.Marshal(in)
	return string(str)
}
