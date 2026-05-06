package common

import (
	"time"

	"icw_common/consts"
	"icw_common/enum"
	"icw_common/gen/activity"
	"icw_common/utils"
)

// ClassificationInfo 输出标准分类能力日志
func ClassificationInfo(requestId, taskUuid, imageUuid string, taskCodes []string, cost time.Duration) {
	utils.LogInfo(consts.LogScopeClassification, consts.LogColorBoldPurple, "[%s] %s %13v %s task_uuid=%s image_uuid=%s task_codes=%s",
		requestId,
		consts.LogColorBoldBlackOnWhite, cost, consts.LogColorReset,
		taskUuid,
		imageUuid,
		utils.JSONF(taskCodes),
	)
}

// ClassificationError 输出失败分类能力日志
func ClassificationError(requestId, taskUuid, imageUuid string, taskCodes []string, cost time.Duration, err error) {
	utils.LogError(consts.LogScopeClassification, "[%s] %s %13v %s task_uuid=%s image_uuid=%s task_codes=%d err=%s",
		requestId,
		consts.LogColorBoldBlackOnWhite, cost, consts.LogColorReset,
		taskUuid,
		imageUuid,
		utils.JSONF(taskCodes),
		utils.FormatErrorLog(err),
	)
}

// CallbackInfo 输出标准回调日志
func CallbackInfo(requestId, taskUuid, imageUuid, status string, start time.Time) {
	utils.LogInfo(consts.LogScopeCallback, consts.LogColorBoldBlue, "[%s] %s %13v %s [CLASSIFICATION] task_uuid=%s image_uuid=%s status=%s%s%s",
		requestId,
		consts.LogColorBoldBlackOnWhite, time.Since(start), consts.LogColorReset,
		taskUuid,
		imageUuid,
		utils.If[string](status == enum.DetectionStatusString(activitypb.DetectionStatus_Succeeded), consts.LogColorBoldGreen, consts.LogColorBoldRed), status, consts.LogColorReset,
	)
}

// CallbackError 输出失败回调日志
func CallbackError(requestId, taskUuid, imageUuid, status string, start time.Time, err error) {
	utils.LogError(consts.LogScopeCallback, "[%s] %s %13v %s [CLASSIFICATION] task_uuid=%s image_uuid=%s status=%s%s%s err=%s",
		requestId,
		consts.LogColorBoldBlackOnWhite, time.Since(start), consts.LogColorReset,
		taskUuid,
		imageUuid,
		utils.If[string](status == enum.DetectionStatusString(activitypb.DetectionStatus_Succeeded), consts.LogColorBoldGreen, consts.LogColorBoldRed), status, consts.LogColorReset,
		utils.FormatErrorLog(err),
	)
}
