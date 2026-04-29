package utils

import (
	"encoding/json"
	"strings"

	"github.com/gin-gonic/gin"
)

// JSONF 将任意结构格式化为 JSON 字符串
func JSONF(v interface{}) string {
	bytes, err := json.Marshal(v)
	if err != nil {
		return ""
	}
	return string(bytes)
}

// BearerToken 从 HTTP Header 中解析 Bearer Token
// Header 格式："Authorization: Bearer <token>"
func BearerToken(c *gin.Context) string {
	header := c.GetHeader("Authorization")
	if header == "" {
		return ""
	}
	const prefix = "Bearer "
	if !strings.HasPrefix(header, prefix) {
		return ""
	}
	return strings.TrimSpace(strings.TrimPrefix(header, prefix))
}
