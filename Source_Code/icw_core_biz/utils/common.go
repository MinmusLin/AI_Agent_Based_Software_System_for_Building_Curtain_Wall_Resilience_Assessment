package utils

import (
	"encoding/json"
	"fmt"
	"strings"

	"icw_core_biz/consts"
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

// FormatErrorLog 格式化错误日志内容
func FormatErrorLog(err interface{}) string {
	msg := strings.TrimSpace(fmt.Sprint(err))
	if IsEmptyError(err) {
		return msg
	}
	return consts.LogColorBoldPurple + msg + consts.LogColorReset
}
