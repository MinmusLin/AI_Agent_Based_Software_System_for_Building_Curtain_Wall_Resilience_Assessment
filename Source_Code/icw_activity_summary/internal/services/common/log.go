package common

import (
	"strings"
	"time"

	"icw_common/consts"
	"icw_common/utils"
)

// SummaryInfo 输出标准总结能力日志
func SummaryInfo(requestId, summaryType, targetId string, cost time.Duration, output string) {
	utils.LogInfo(consts.LogScopeSummary, consts.LogColorBoldPurple, "[%s] %s %13v %s summary_type=%s target_id=%s output=%s",
		utils.If[string](requestId == "", "-", requestId),
		consts.LogColorBoldBlackOnWhite, cost, consts.LogColorReset,
		summaryType,
		targetId,
		output,
	)
}

// SummaryError 输出失败总结能力日志
func SummaryError(requestId, summaryType, targetId string, cost time.Duration, output string, err error) {
	utils.LogError(consts.LogScopeSummary, "[%s] %s %13v %s summary_type=%s target_id=%s output=%s err=%s",
		utils.If[string](requestId == "", "-", requestId),
		consts.LogColorBoldBlackOnWhite, cost, consts.LogColorReset,
		summaryType,
		targetId,
		output,
		utils.FormatErrorLog(err),
	)
}

// CallbackInfo 输出标准回调日志
func CallbackInfo(requestId, summaryType, targetId, status string, start time.Time) {
	utils.LogInfo(consts.LogScopeCallback, consts.LogColorBoldBlue, "[%s] %s %13v %s [SUMMARY|%s] target_id=%s status=%s",
		utils.If[string](requestId == "", "-", requestId),
		consts.LogColorBoldBlackOnWhite, time.Since(start), consts.LogColorReset,
		strings.ToUpper(summaryType),
		targetId,
		utils.FormatSuccessLog(status),
	)
}

// CallbackError 输出失败回调日志
func CallbackError(requestId, summaryType, targetId, status string, start time.Time, err error) {
	utils.LogError(consts.LogScopeCallback, "[%s] %s %13v %s [SUMMARY|%s] target_id=%s status=%s err=%s",
		utils.If[string](requestId == "", "-", requestId),
		consts.LogColorBoldBlackOnWhite, time.Since(start), consts.LogColorReset,
		strings.ToUpper(summaryType),
		targetId,
		utils.FormatErrorLog(status),
		utils.FormatErrorLog(err),
	)
}
