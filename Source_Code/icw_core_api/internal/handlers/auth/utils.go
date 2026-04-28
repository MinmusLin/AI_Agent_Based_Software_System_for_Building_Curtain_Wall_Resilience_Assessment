package auth

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"

	"icw_core_api/internal/response"
)

func bindJSON(c *gin.Context, out interface{}) bool {
	if err := c.ShouldBindJSON(out); err != nil {
		response.Error(c, http.StatusBadRequest, "BAD_REQUEST", "invalid request body")
		return false
	}
	return true
}

func writeRPCError(c *gin.Context, err error) {
	msg := "unknown error"
	if err != nil {
		msg = err.Error()
	}
	response.Error(c, http.StatusInternalServerError, "INTERNAL_ERROR", msg)
}

func bearerToken(c *gin.Context) string {
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
