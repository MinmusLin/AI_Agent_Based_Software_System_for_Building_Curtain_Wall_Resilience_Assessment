package utils

import (
	"encoding/json"
	"fmt"
	"strings"
)

// JSONF 将任意结构格式化为 JSON 字符串
func JSONF(v interface{}) string {
	bytes, err := json.Marshal(v)
	if err != nil {
		return ""
	}
	return string(bytes)
}

// IsEmptyError 判断错误是否为空错误
func IsEmptyError(err interface{}) bool {
	msg := strings.TrimSpace(fmt.Sprint(err))
	return msg == "" || msg == "nil" || msg == "<nil>"
}
