package rocketmq

import (
	"time"

	"icw_core_biz/consts"
	"icw_core_biz/utils"
)

// startMQLog 开始记录 MQ 调用日志
func startMQLog(method, topic string) func(error) {
	start := time.Now()
	return func(err error) {
		if utils.IsEmptyError(err) {
			MQInfo("[%s] topic=%s cost=%s", method, topic, time.Since(start))
			return
		}
		MQError("[%s] topic=%s cost=%s err=%s", method, topic, time.Since(start), utils.FormatErrorLog(err))
	}
}

// MQInfo 记录 MQ 正常日志
func MQInfo(format string, args ...interface{}) {
	utils.LogInfo(consts.LogScopeMQ, consts.LogColorBoldCyan, format, args...)
}

// MQError 记录 MQ 错误日志
func MQError(format string, args ...interface{}) {
	utils.LogError(consts.LogScopeMQ, format, args...)
}

// MQWarn 记录 MQ 警告日志
func MQWarn(format string, args ...interface{}) {
	utils.LogWarn(consts.LogScopeMQ, format, args...)
}
