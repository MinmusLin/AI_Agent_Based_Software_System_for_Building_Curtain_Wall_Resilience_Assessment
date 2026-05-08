package utils

import (
	"strings"
)

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
