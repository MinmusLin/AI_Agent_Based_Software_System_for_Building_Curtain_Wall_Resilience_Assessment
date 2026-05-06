package common

import (
	"strings"
	"time"

	reasoningConsts "icw_activity_reasoning/consts"
	"icw_common/consts"
	"icw_common/utils"
)

// ReasoningInfo 输出标准原子检测能力日志
func ReasoningInfo(requestId, taskCode, taskUuid, imageUuid string, artifactCount int, cost time.Duration) {
	utils.LogInfo(consts.LogScopeReasoning, consts.LogColorBoldPurple, "[%s] %s %13v %s [%s] task_uuid=%s image_uuid=%s artifacts=%d",
		requestId,
		consts.LogColorBoldBlackOnWhite, cost, consts.LogColorReset,
		strings.ToUpper(strings.TrimSpace(taskCode)),
		taskUuid,
		imageUuid,
		artifactCount,
	)
}

// ReasoningError 输出失败原子检测能力日志
func ReasoningError(requestId, taskCode, taskUuid, imageUuid string, artifactCount int, cost time.Duration, err error) {
	utils.LogError(consts.LogScopeReasoning, "[%s] %s %13v %s [%s] task_uuid=%s image_uuid=%s artifacts=%d err=%s",
		requestId,
		consts.LogColorBoldBlackOnWhite, cost, consts.LogColorReset,
		strings.ToUpper(strings.TrimSpace(taskCode)),
		taskUuid,
		imageUuid,
		artifactCount,
		utils.FormatErrorLog(err),
	)
}

// CallbackInfo 输出标准回调日志
func CallbackInfo(requestId, taskCode, taskUuid, imageUuid, status string, cost time.Duration) {
	utils.LogInfo(consts.LogScopeCallback, consts.LogColorBoldBlue, "[%s] %s %13v %s [%s] task_uuid=%s image_uuid=%s status=%s%s%s",
		requestId,
		consts.LogColorBoldBlackOnWhite, cost, consts.LogColorReset,
		strings.ToUpper(strings.TrimSpace(taskCode)),
		taskUuid,
		imageUuid,
		utils.If[string](status == reasoningConsts.DetectionStatusSucceeded, consts.LogColorBoldGreen, consts.LogColorBoldRed), status, consts.LogColorReset,
	)
}

// CallbackError 输出失败回调日志
func CallbackError(requestId, taskCode, taskUuid, imageUuid, status string, cost time.Duration, err error) {
	utils.LogError(consts.LogScopeCallback, "[%s] %s %13v %s [%s] task_uuid=%s image_uuid=%s status=%s%s%s err=%s",
		requestId,
		consts.LogColorBoldBlackOnWhite, cost, consts.LogColorReset,
		strings.ToUpper(strings.TrimSpace(taskCode)),
		taskUuid,
		imageUuid,
		utils.If[string](status == reasoningConsts.DetectionStatusSucceeded, consts.LogColorBoldGreen, consts.LogColorBoldRed), status, consts.LogColorReset,
		utils.FormatErrorLog(err),
	)
}
