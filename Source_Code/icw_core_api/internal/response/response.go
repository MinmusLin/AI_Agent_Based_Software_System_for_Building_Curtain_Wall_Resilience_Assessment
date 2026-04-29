package response

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

// OKEnvelope 标准成功响应
type OKEnvelope[T any] struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Data    T      `json:"data"`
}

// ErrorEnvelope 标准失败响应
type ErrorEnvelope struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

// OK 写入标准成功响应
func OK[T any](c *gin.Context, data T) {
	c.JSON(http.StatusOK, OKEnvelope[T]{
		Code:    0,
		Message: "ok",
		Data:    data,
	})
}

// Error 写入标准失败响应
func Error(c *gin.Context, status int, code, message string) {
	c.JSON(status, ErrorEnvelope{
		Code:    code,
		Message: message,
	})
}

// BindJSON 绑定 JSON 请求体
func BindJSON(c *gin.Context, out interface{}) bool {
	if err := c.ShouldBindJSON(out); err != nil {
		Error(c, http.StatusBadRequest, "BAD_REQUEST", "bad request")
		return false
	}
	return true
}

// WriteRPCError 将 RPC 层返回的基础错误转换为 API 层的 HTTP 响应
func WriteRPCError(c *gin.Context, err error) {
	msg := "server internal error"
	if err != nil {
		msg = err.Error()
	}

	switch {
	case strings.Contains(msg, "invalid request"):
		// 无效请求
		Error(c, http.StatusBadRequest, "BAD_REQUEST", msg)
	case strings.Contains(msg, "unauthorized"):
		// 身份认证未通过
		Error(c, http.StatusUnauthorized, "UNAUTHORIZED", msg)
	case strings.Contains(msg, "account locked"):
		// 账号锁定（密码登录失败次数达上限 / 邮箱验证码发送次数达上限）
		Error(c, http.StatusLocked, "ACCOUNT_LOCKED", msg)
	default:
		// 服务器内部错误
		Error(c, http.StatusInternalServerError, "INTERNAL_ERROR", msg)
	}
}
