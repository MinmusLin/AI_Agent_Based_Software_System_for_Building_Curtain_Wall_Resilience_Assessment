package enum

import (
	"strings"

	"icw_common/gen/activity"
)

// DetectionStatusString 将任务状态枚举转换为字符串
func DetectionStatusString(status activitypb.DetectionStatus_Value) string {
	switch status {
	case activitypb.DetectionStatus_Succeeded:
		return "succeeded"
	case activitypb.DetectionStatus_Failed:
		return "failed"
	default:
		return ""
	}
}

// ParseDetectionStatus 将存储值转换为任务状态枚举
func ParseDetectionStatus(value string) activitypb.DetectionStatus_Value {
	switch strings.TrimSpace(value) {
	case "succeeded":
		return activitypb.DetectionStatus_Succeeded
	case "failed":
		return activitypb.DetectionStatus_Failed
	default:
		return activitypb.DetectionStatus_Unknown
	}
}

// DetectionTaskCodeString 将原子检测能力代码枚举转换为字符串
func DetectionTaskCodeString(code activitypb.DetectionTaskCode_Value) string {
	switch code {
	case activitypb.DetectionTaskCode_Corrosion:
		return "corrosion"
	case activitypb.DetectionTaskCode_Crack:
		return "crack"
	case activitypb.DetectionTaskCode_Stain:
		return "stain"
	case activitypb.DetectionTaskCode_Flatness:
		return "flatness"
	case activitypb.DetectionTaskCode_Spalling:
		return "spalling"
	default:
		return ""
	}
}

// ParseDetectionTaskCode 将存储值转换为原子检测能力代码枚举
func ParseDetectionTaskCode(value string) activitypb.DetectionTaskCode_Value {
	switch strings.TrimSpace(value) {
	case "corrosion":
		return activitypb.DetectionTaskCode_Corrosion
	case "crack":
		return activitypb.DetectionTaskCode_Crack
	case "stain":
		return activitypb.DetectionTaskCode_Stain
	case "flatness":
		return activitypb.DetectionTaskCode_Flatness
	case "spalling":
		return activitypb.DetectionTaskCode_Spalling
	default:
		return activitypb.DetectionTaskCode_Unknown
	}
}
