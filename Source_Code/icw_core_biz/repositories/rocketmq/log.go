package rocketmq

import (
	"icw_common/consts"
	"icw_common/utils"
)

// MQInfo 输出标准 MQ 日志
func MQInfo(format string, args ...interface{}) {
	utils.LogInfo(consts.LogScopeMQ, consts.LogColorBoldCyan, format, args...)
}

// MQError 输出错误 MQ 日志
func MQError(format string, args ...interface{}) {
	utils.LogError(consts.LogScopeMQ, format, args...)
}
