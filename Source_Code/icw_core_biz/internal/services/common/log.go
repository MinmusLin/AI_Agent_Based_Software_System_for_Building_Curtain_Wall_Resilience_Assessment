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
