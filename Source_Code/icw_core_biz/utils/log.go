package utils

import (
	"log"

	"icw_core_biz/consts"
)

// LogInfo 输出标准日志
func LogInfo(scope, color, format string, args ...interface{}) {
	logWithScope(logInfoPrefix(scope, color), format, args...)
}

// LogWarn 输出警告日志
func LogWarn(scope, format string, args ...interface{}) {
	logWithScope(logWarnPrefix(scope), format, args...)
}

// LogError 输出错误日志
func LogError(scope, format string, args ...interface{}) {
	logWithScope(logErrorPrefix(scope), format, args...)
}

// LogFault 输出致命错误日志并退出进程
func LogFault(scope, format string, args ...interface{}) {
	log.Printf("%s "+format, append([]interface{}{logFaultPrefix(scope)}, args...)...)
}

// logWithScope 输出日志所属功能域和内容
func logWithScope(prefix, format string, args ...interface{}) {
	log.Printf("%s "+format, append([]interface{}{prefix}, args...)...)
}

// logInfoPrefix 全局标准日志前缀
func logInfoPrefix(scope, color string) string {
	if scope == "" {
		if color == "" {
			return "[INFO]"
		}
		return color + "[INFO]" + consts.LogColorReset
	}
	if color == "" {
		return "[" + scope + " INFO]"
	}
	return color + "[" + scope + " INFO]" + consts.LogColorReset
}

// logErrorPrefix 全局错误日志前缀
func logErrorPrefix(scope string) string {
	if scope == "" {
		return "[" + scope + "ERROR]"
	}
	return consts.LogColorBoldRed + "[" + scope + " ERROR]" + consts.LogColorReset
}

// logWarnPrefix 全局警告日志前缀
func logWarnPrefix(scope string) string {
	if scope == "" {
		return "[" + scope + "WARN]"
	}
	return consts.LogColorBoldYellow + "[" + scope + " WARN]" + consts.LogColorReset
}

// logFaultPrefix 全局致命错误日志前缀
func logFaultPrefix(scope string) string {
	if scope == "" {
		return "[FAULT]"
	}
	return consts.LogColorBoldWhiteOnRed + "[" + scope + " FAULT]" + consts.LogColorBoldWhiteOnRed
}
