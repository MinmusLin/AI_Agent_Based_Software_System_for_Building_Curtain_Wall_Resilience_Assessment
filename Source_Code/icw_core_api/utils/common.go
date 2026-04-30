package utils

import (
	"encoding/json"
	"errors"
	"strings"

	"github.com/gin-gonic/gin"

	"icw_core_api/configs"
	bizDto "icw_core_biz/pkg/dto"
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

// GetCurrentUser 从 Gin Context 中获取当前登录用户
func GetCurrentUser(c *gin.Context) (*bizDto.User, error) {
	value, ok := c.Get(configs.ContextUser)
	if !ok || value == nil {
		return nil, errors.New("current user not found in Gin context")
	}
	user, ok := value.(*bizDto.User)
	if !ok || user == nil || user.Id == 0 || user.Email == "" || user.Name == "" {
		return nil, errors.New("invalid current user in Gin context")
	}
	return user, nil
}
