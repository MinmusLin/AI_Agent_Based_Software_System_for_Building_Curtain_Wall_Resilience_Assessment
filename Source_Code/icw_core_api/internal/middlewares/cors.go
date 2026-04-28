package middlewares

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func CORS() gin.HandlerFunc {
	return func(c *gin.Context) {
		origin := c.GetHeader("Origin")
		if origin != "" {
			c.Header("Access-Control-Allow-Origin", origin)
			c.Header("Access-Control-Allow-Credentials", "true")
			c.Header("Vary", "Origin, Access-Control-Request-Method, Access-Control-Request-Headers")
		}

		requestHeaders := c.GetHeader("Access-Control-Request-Headers")
		if requestHeaders == "" {
			requestHeaders = "Content-Type, Authorization"
		}
		c.Header("Access-Control-Allow-Headers", requestHeaders)

		requestMethod := c.GetHeader("Access-Control-Request-Method")
		if requestMethod == "" {
			requestMethod = "GET, HEAD, POST, PUT, PATCH, DELETE, CONNECT, OPTIONS, TRACE"
		}
		c.Header("Access-Control-Allow-Methods", requestMethod)

		if c.Request.Method == http.MethodOptions {
			c.AbortWithStatus(http.StatusNoContent)
			return
		}

		c.Next()
	}
}
