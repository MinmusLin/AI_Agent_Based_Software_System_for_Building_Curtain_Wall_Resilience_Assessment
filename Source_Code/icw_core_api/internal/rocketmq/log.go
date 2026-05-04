package rocketmq

import (
	"icw_core_biz/consts"
	"icw_core_biz/utils"
)

// MQInfo 输出标准 MQ 日志
func MQInfo(format string, args ...interface{}) {
	utils.LogInfo(consts.LogScopeMQ, consts.LogColorBoldCyan, format, args...)
}

// MQWarn 输出警告 MQ 日志
func MQWarn(format string, args ...interface{}) {
	utils.LogWarn(consts.LogScopeMQ, format, args...)
}

// MQError 输出错误 MQ 日志
func MQError(format string, args ...interface{}) {
	utils.LogError(consts.LogScopeMQ, format, args...)
}

// MQFatal 输出致命 MQ 日志并退出进程
func MQFatal(format string, args ...interface{}) {
	utils.LogFatal(consts.LogScopeMQ, format, args...)
}
