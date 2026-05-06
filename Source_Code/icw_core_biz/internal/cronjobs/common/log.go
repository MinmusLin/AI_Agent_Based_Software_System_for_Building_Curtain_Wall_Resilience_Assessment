package common

import (
	"time"

	"icw_common/consts"
	"icw_common/utils"
)

// cronLog 输出定时任务执行日志
func cronLog(name string, start time.Time, result interface{}, err error) {
	resultStr := utils.JSONF(result)
	if utils.IsEmptyError(err) {
		if resultStr == "" {
			CronInfo("[%s] %s %13v %s result={}",
				name,
				consts.LogColorBoldBlackOnWhite, time.Since(start), consts.LogColorReset,
			)
			return
		}
		CronInfo("[%s] %s %13v %s result=%s",
			name,
			consts.LogColorBoldBlackOnWhite, time.Since(start), consts.LogColorReset,
			resultStr,
		)
		return
	}
	if resultStr == "" {
		CronError("[%s] %s %13v %s result={} err=%s",
			name,
			consts.LogColorBoldBlackOnWhite, time.Since(start), consts.LogColorReset,
			utils.FormatErrorLog(err),
		)
		return
	}
	CronError("[%s] %s %13v %s result=%s err=%s",
		name,
		consts.LogColorBoldBlackOnWhite, time.Since(start), consts.LogColorReset,
		resultStr,
		utils.FormatErrorLog(err),
	)
}

// CronInfo 输出标准定时任务日志
func CronInfo(format string, args ...interface{}) {
	utils.LogInfo(consts.LogScopeCron, consts.LogColorBoldPurple, format, args...)
}

// CronWarn 输出警告定时任务日志
func CronWarn(format string, args ...interface{}) {
	utils.LogWarn(consts.LogScopeCron, format, args...)
}

// CronError 输出错误定时任务日志
func CronError(format string, args ...interface{}) {
	utils.LogError(consts.LogScopeCron, format, args...)
}

// CronFatal 输出致命错误定时任务日志并退出进程
func CronFatal(format string, args ...interface{}) {
	utils.LogFatal(consts.LogScopeCron, format, args...)
}
