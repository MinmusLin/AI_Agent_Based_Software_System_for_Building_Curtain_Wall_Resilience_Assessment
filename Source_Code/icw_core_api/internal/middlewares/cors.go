package middlewares

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// CORS 跨域配置中间件
func CORS() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 回显请求来源 Origin，并允许携带 Cookie 等凭证信息
		origin := c.GetHeader("Origin")
		if origin != "" {
			c.Header("Access-Control-Allow-Origin", origin)
			c.Header("Access-Control-Allow-Credentials", "true")
			// 缓存代理配置：不同 Origin / 预检方法 / 预检头会得到不同响应
			c.Header("Vary", "Origin, Access-Control-Request-Method, Access-Control-Request-Headers")
		}

		// 普通请求没有预检头时，保留项目当前需要的基础请求头
		requestHeaders := c.GetHeader("Access-Control-Request-Headers")
		if requestHeaders == "" {
			requestHeaders = "Content-Type, Authorization"
		}
		c.Header("Access-Control-Allow-Headers", requestHeaders)

		// 没有预检方法时，默认放开所有 HTTP 方法
		requestMethod := c.GetHeader("Access-Control-Request-Method")
		if requestMethod == "" {
			requestMethod = "GET, HEAD, POST, PUT, PATCH, DELETE, CONNECT, OPTIONS, TRACE"
		}
		c.Header("Access-Control-Allow-Methods", requestMethod)

		// OPTIONS 是浏览器预检请求，只需要返回跨域响应头，不进入业务 Handler
		if c.Request.Method == http.MethodOptions {
			c.AbortWithStatus(http.StatusNoContent)
			return
		}

		c.Next()
	}
}
