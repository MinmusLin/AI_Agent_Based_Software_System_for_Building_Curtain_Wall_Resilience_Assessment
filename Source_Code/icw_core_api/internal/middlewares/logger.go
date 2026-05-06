package middlewares

import (
	"fmt"
	"time"

	"github.com/gin-gonic/gin"

	"icw_common/consts"
	"icw_core_api/utils"
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
			if v := utils.GetXRequestId(param.Request.Context()); v != "" {
				requestId = v
			}
		}

		if param.Latency > time.Minute {
			param.Latency = param.Latency.Truncate(time.Second)
		}

		return fmt.Sprintf("%v %s[HTTP INFO]%s [%s] %s %13v %s %s %-7s %s %s %3d %s %s\n%s",
			param.TimeStamp.Format("2006/01/02 15:04:05"),
			consts.LogColorBoldGreen, consts.LogColorReset,
			requestId,
			consts.LogColorBoldBlackOnWhite, param.Latency, consts.LogColorReset,
			methodColor, param.Method, resetColor,
			statusColor, param.StatusCode, resetColor,
			param.Path,
			param.ErrorMessage,
		)
	})
}
