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

func BindJSON(c *gin.Context, out interface{}) bool {
	if err := c.ShouldBindJSON(out); err != nil {
		Error(c, http.StatusBadRequest, "BAD_REQUEST", "invalid request body")
		return false
	}
	return true
}

func WriteRPCError(c *gin.Context, err error) {
	msg := "unknown error"
	if err != nil {
		msg = err.Error()
	}
	Error(c, http.StatusInternalServerError, "INTERNAL_ERROR", msg)
}
