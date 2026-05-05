package utils

import (
	"context"
	"errors"
	"strings"

	"github.com/gin-gonic/gin"

	"icw_common/dto"
	"icw_core_api/consts"
)

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
func GetCurrentUser(c *gin.Context) (*dto.User, error) {
	value, ok := c.Get(consts.ContextUser)
	if !ok || value == nil {
		return nil, errors.New("current user not found in Gin context")
	}
	user, ok := value.(*dto.User)
	if !ok || user == nil || user.Id == 0 || user.Email == "" || user.Name == "" {
		return nil, errors.New("invalid current user in Gin context")
	}
	return user, nil
}

// GetRequestId 从请求上下文中获取请求 ID
func GetRequestId(ctx context.Context) string {
	if ctx == nil {
		return ""
	}
	requestId, ok := ctx.Value(consts.ContextRequestId).(string)
	if !ok {
		return ""
	}
	return strings.TrimSpace(requestId)
}
