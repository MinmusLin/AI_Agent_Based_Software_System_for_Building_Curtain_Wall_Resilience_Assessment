package common

import (
	"icw_common/consts"
	"icw_common/utils"
)

// RpcInfo 输出标准 RPC 日志
func RpcInfo(format string, args ...interface{}) {
	utils.LogInfo(consts.LogScopeRPC, consts.LogColorBoldGreen, format, args...)
}

// RpcWarn 输出警告 RPC 日志
func RpcWarn(format string, args ...interface{}) {
	utils.LogWarn(consts.LogScopeRPC, format, args...)
}

// RpcError 输出错误 RPC 日志
func RpcError(format string, args ...interface{}) {
	utils.LogError(consts.LogScopeRPC, format, args...)
}

// RpcFatal 输出致命错误 RPC 日志并退出进程
func RpcFatal(format string, args ...interface{}) {
	utils.LogFatal(consts.LogScopeRPC, format, args...)
}
