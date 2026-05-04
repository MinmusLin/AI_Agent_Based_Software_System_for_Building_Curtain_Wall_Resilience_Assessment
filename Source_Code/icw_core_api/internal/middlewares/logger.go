package middlewares

import (
	"fmt"
	"time"

	"github.com/gin-gonic/gin"

	"icw_core_api/utils"
	"icw_core_biz/consts"
)

// Logger 保持 Gin 默认日志输出格式，添加请求 ID 输出
func Logger() gin.HandlerFunc {
	return gin.LoggerWithFormatter(func(param gin.LogFormatterParams) string {
		var statusColor, methodColor, resetColor string
		if param.IsOutputColor() {
			statusColor = param.StatusCodeColor()
			methodColor = param.MethodColor()
			resetColor = param.ResetColor()
		}

		requestId := "-"
		if param.Request != nil {
			if v := utils.GetRequestId(param.Request.Context()); v != "" {
				requestId = v
			}
		}

		if param.Latency > time.Minute {
			param.Latency = param.Latency.Truncate(time.Second)
		}

		return fmt.Sprintf("%s[HTTP]%s %v | %s |%s %3d %s| %13v | %15s |%s %-7s %s %#v\n%s",
			consts.LogColorBoldGreen,
			consts.LogColorReset,
			param.TimeStamp.Format("2006/01/02 15:04:05"),
			requestId,
			statusColor, param.StatusCode, resetColor,
			param.Latency,
			param.ClientIP,
			methodColor, param.Method, resetColor,
			param.Path,
			param.ErrorMessage,
		)
	})
}
