package rocketmq

import (
	"icw_core_biz/consts"
	"icw_core_biz/utils"
)

// MQSuccess 输出成功 MQ 日志
func MQSuccess(format string, args ...interface{}) {
	utils.LogInfo(consts.LogScopeMQ, consts.LogColorBoldCyan, format, args...)
}

// MQError 输出错误 MQ 日志
func MQError(format string, args ...interface{}) {
	utils.LogError(consts.LogScopeMQ, format, args...)
}
