package response

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type OKEnvelope[T any] struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Data    T      `json:"data"`
}

type ErrorEnvelope struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

func OK[T any](c *gin.Context, data T) {
	c.JSON(http.StatusOK, OKEnvelope[T]{
		Code:    0,
		Message: "ok",
		Data:    data,
	})
}

func Error(c *gin.Context, status int, code, message string) {
	c.JSON(status, ErrorEnvelope{
		Code:    code,
		Message: message,
	})
}
