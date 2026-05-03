package middlewares

import (
	"context"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"icw_core_api/consts"
)

// RequestId 为每个 HTTP 请求生成请求 ID，并写入 Header 与请求上下文
func RequestId() gin.HandlerFunc {
	return func(c *gin.Context) {
		requestId := strings.TrimSpace(c.GetHeader(consts.HeaderRequestId))
		if requestId == "" {
			requestId = uuid.NewString()
		}

		c.Header(consts.HeaderRequestId, requestId)
		c.Request = c.Request.WithContext(context.WithValue(c.Request.Context(), consts.ContextRequestId, requestId))
		c.Next()
	}
}
