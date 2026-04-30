package response

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"icw_core_biz/pkg/rpc_err"
)

// OKEnvelope 标准成功响应
type OKEnvelope[T any] struct {
	Code    string `json:"code"`
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
		Code:    "OK",
		Message: "success",
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
		Error(c, http.StatusBadRequest, string(rpc_err.DetailBadRequest), errorMessage(rpc_err.DetailBadRequest))
		return false
	}
	return true
}
