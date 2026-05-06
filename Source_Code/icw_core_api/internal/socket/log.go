package socket

import (
	"icw_common/consts"
	"icw_common/utils"
)

// WSInfo 输出标准 WebSocket 日志
func WSInfo(format string, args ...interface{}) {
	utils.LogInfo(consts.LogScopeWS, consts.LogColorBoldPink, format, args...)
}

// WSError 输出错误 WebSocket 日志
func WSError(format string, args ...interface{}) {
	utils.LogError(consts.LogScopeWS, format, args...)
}
