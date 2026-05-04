package socket

import (
	"icw_core_biz/consts"
	"icw_core_biz/utils"
)

// WSInfo 输出标准 WebSocket 日志
func WSInfo(format string, args ...interface{}) {
	utils.LogInfo(consts.LogScopeWS, consts.LogColorBoldPink, format, args...)
}

// WSWarn 输出警告 WebSocket 日志
func WSWarn(format string, args ...interface{}) {
	utils.LogWarn(consts.LogScopeWS, format, args...)
}

// WSError 输出错误 WebSocket 日志
func WSError(format string, args ...interface{}) {
	utils.LogError(consts.LogScopeWS, format, args...)
}

// WSFatal 输出致命 WebSocket 日志并退出进程
func WSFatal(format string, args ...interface{}) {
	utils.LogFatal(consts.LogScopeWS, format, args...)
}
