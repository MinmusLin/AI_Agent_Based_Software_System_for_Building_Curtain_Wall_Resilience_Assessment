package utils

import (
	"encoding/json"
	"fmt"
	"reflect"
	"strings"
)

// If 泛型三元运算符
func If[T any](cond bool, onTrue, onFalse T) T {
	if cond {
		return onTrue
	}
	return onFalse
}

// JSONF 将任意结构格式化为 JSON 字符串
func JSONF(v interface{}) string {
	bytes, err := json.Marshal(v)
	if err != nil {
		return ""
	}
	return string(bytes)
}

// IsNil 判断接口是否为空
func IsNil(value interface{}) bool {
	if value == nil {
		return true
	}
	reflectValue := reflect.ValueOf(value)
	switch reflectValue.Kind() {
	case reflect.Chan, reflect.Func, reflect.Interface, reflect.Map, reflect.Ptr, reflect.Slice:
		return reflectValue.IsNil()
	default:
		return false
	}
}

// IsEmptyError 判断错误是否为空错误
func IsEmptyError(err interface{}) bool {
	msg := strings.TrimSpace(fmt.Sprint(err))
	return msg == "" || msg == "nil" || msg == "<nil>"
}

// FirstNotEmpty 返回第一个非空字符串
func FirstNotEmpty(values ...string) string {
	for _, value := range values {
		value = strings.TrimSpace(value)
		if value != "" {
			return value
		}
	}
	return ""
}
