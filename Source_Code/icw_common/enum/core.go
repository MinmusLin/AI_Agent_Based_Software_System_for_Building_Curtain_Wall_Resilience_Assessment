package enum

import (
	"strings"

	"icw_common/gen/core/common"
)

// EmailCodeSceneString 将邮箱验证码业务场景枚举转换为字符串
func EmailCodeSceneString(scene commonpb.EmailCodeScene_Value) string {
	switch scene {
	case commonpb.EmailCodeScene_Register:
		return "register"
	case commonpb.EmailCodeScene_Login:
		return "login"
	case commonpb.EmailCodeScene_Reset:
		return "reset"
	default:
		return ""
	}
}

// ParseEmailCodeScene 将存储值转换为邮箱验证码业务场景枚举
func ParseEmailCodeScene(value string) commonpb.EmailCodeScene_Value {
	switch strings.TrimSpace(value) {
	case "register":
		return commonpb.EmailCodeScene_Register
	case "login":
		return commonpb.EmailCodeScene_Login
	case "reset":
		return commonpb.EmailCodeScene_Reset
	default:
		return commonpb.EmailCodeScene_Unknown
	}
}

// LoginSceneString 将登录方式枚举转换为字符串
func LoginSceneString(scene commonpb.LoginScene_Value) string {
	switch scene {
	case commonpb.LoginScene_Password:
		return "password"
	case commonpb.LoginScene_Email:
		return "email"
	default:
		return ""
	}
}

// ParseLoginScene 将存储值转换为登录方式枚举
func ParseLoginScene(value string) commonpb.LoginScene_Value {
	switch strings.TrimSpace(value) {
	case "password":
		return commonpb.LoginScene_Password
	case "email":
		return commonpb.LoginScene_Email
	default:
		return commonpb.LoginScene_Unknown
	}
}

// EmailSendStatusString 将邮件发送状态枚举转换为字符串
func EmailSendStatusString(status commonpb.EmailSendStatus_Value) string {
	switch status {
	case commonpb.EmailSendStatus_Success:
		return "success"
	case commonpb.EmailSendStatus_Failed:
		return "failed"
	default:
		return ""
	}
}

// ParseEmailSendStatus 将存储值转换为邮件发送状态枚举
func ParseEmailSendStatus(value string) commonpb.EmailSendStatus_Value {
	switch strings.TrimSpace(value) {
	case "success":
		return commonpb.EmailSendStatus_Success
	case "failed":
		return commonpb.EmailSendStatus_Failed
	default:
		return commonpb.EmailSendStatus_Unknown
	}
}

// ProjectStatusString 将项目状态枚举转换为字符串
func ProjectStatusString(status commonpb.ProjectStatus_Value) string {
	switch status {
	case commonpb.ProjectStatus_Active:
		return "active"
	case commonpb.ProjectStatus_Completed:
		return "completed"
	case commonpb.ProjectStatus_Deleted:
		return "deleted"
	default:
		return ""
	}
}

// ParseProjectStatus 将存储值转换为项目状态枚举
func ParseProjectStatus(value string) commonpb.ProjectStatus_Value {
	switch strings.TrimSpace(value) {
	case "active":
		return commonpb.ProjectStatus_Active
	case "completed":
		return commonpb.ProjectStatus_Completed
	case "deleted":
		return commonpb.ProjectStatus_Deleted
	default:
		return commonpb.ProjectStatus_Unknown
	}
}

// ProjectProgressUint8 将项目进度枚举转换为 uint8
func ProjectProgressUint8(progress commonpb.ProjectProgress_Value) uint8 {
	switch progress {
	case commonpb.ProjectProgress_InitializationFinished,
		commonpb.ProjectProgress_ProfileFinished,
		commonpb.ProjectProgress_AssetsFinished,
		commonpb.ProjectProgress_DetectionFinished,
		commonpb.ProjectProgress_ReviewFinished,
		commonpb.ProjectProgress_ReportFinished:
		return uint8(progress)
	default:
		return uint8(commonpb.ProjectProgress_InitializationFinished)
	}
}

// ParseProjectProgress 将存储值转换为项目进度枚举
func ParseProjectProgress(value uint8) commonpb.ProjectProgress_Value {
	switch progress := commonpb.ProjectProgress_Value(value); progress {
	case commonpb.ProjectProgress_InitializationFinished,
		commonpb.ProjectProgress_ProfileFinished,
		commonpb.ProjectProgress_AssetsFinished,
		commonpb.ProjectProgress_DetectionFinished,
		commonpb.ProjectProgress_ReviewFinished,
		commonpb.ProjectProgress_ReportFinished:
		return progress
	default:
		return commonpb.ProjectProgress_InitializationFinished
	}
}

// ProjectImageStatusString 将项目图像状态枚举转换为字符串
func ProjectImageStatusString(status commonpb.ProjectImageStatus_Value) string {
	switch status {
	case commonpb.ProjectImageStatus_Pending:
		return "pending"
	case commonpb.ProjectImageStatus_Uploaded:
		return "uploaded"
	case commonpb.ProjectImageStatus_Failed:
		return "failed"
	default:
		return ""
	}
}

// ParseProjectImageStatus 将存储值转换为项目图像状态枚举
func ParseProjectImageStatus(value string) commonpb.ProjectImageStatus_Value {
	switch strings.TrimSpace(value) {
	case "pending":
		return commonpb.ProjectImageStatus_Pending
	case "uploaded":
		return commonpb.ProjectImageStatus_Uploaded
	case "failed":
		return commonpb.ProjectImageStatus_Failed
	default:
		return commonpb.ProjectImageStatus_Unknown
	}
}

// ProjectDetectionTaskStatusString 将项目图像检测主任务状态枚举转换为字符串
func ProjectDetectionTaskStatusString(status commonpb.ProjectDetectionTaskStatus_Value) string {
	switch status {
	case commonpb.ProjectDetectionTaskStatus_Pending:
		return "pending"
	case commonpb.ProjectDetectionTaskStatus_Classifying:
		return "classifying"
	case commonpb.ProjectDetectionTaskStatus_Detecting:
		return "detecting"
	case commonpb.ProjectDetectionTaskStatus_Summarizing:
		return "summarizing"
	case commonpb.ProjectDetectionTaskStatus_Succeeded:
		return "succeeded"
	case commonpb.ProjectDetectionTaskStatus_Failed:
		return "failed"
	default:
		return ""
	}
}

// ParseProjectDetectionTaskStatus 将存储值转换为项目图像检测主任务状态枚举
func ParseProjectDetectionTaskStatus(value string) commonpb.ProjectDetectionTaskStatus_Value {
	switch strings.TrimSpace(value) {
	case "pending":
		return commonpb.ProjectDetectionTaskStatus_Pending
	case "classifying":
		return commonpb.ProjectDetectionTaskStatus_Classifying
	case "detecting":
		return commonpb.ProjectDetectionTaskStatus_Detecting
	case "summarizing":
		return commonpb.ProjectDetectionTaskStatus_Summarizing
	case "succeeded":
		return commonpb.ProjectDetectionTaskStatus_Succeeded
	case "failed":
		return commonpb.ProjectDetectionTaskStatus_Failed
	default:
		return commonpb.ProjectDetectionTaskStatus_Unknown
	}
}

// ProjectDetectionSubTaskStatusString 将项目图像检测子任务状态枚举转换为字符串
func ProjectDetectionSubTaskStatusString(status commonpb.ProjectDetectionSubTaskStatus_Value) string {
	switch status {
	case commonpb.ProjectDetectionSubTaskStatus_Pending:
		return "pending"
	case commonpb.ProjectDetectionSubTaskStatus_Succeeded:
		return "succeeded"
	case commonpb.ProjectDetectionSubTaskStatus_Failed:
		return "failed"
	default:
		return ""
	}
}

// ParseProjectDetectionSubTaskStatus 将存储值转换为项目图像检测子任务状态枚举
func ParseProjectDetectionSubTaskStatus(value string) commonpb.ProjectDetectionSubTaskStatus_Value {
	switch strings.TrimSpace(value) {
	case "pending":
		return commonpb.ProjectDetectionSubTaskStatus_Pending
	case "succeeded":
		return commonpb.ProjectDetectionSubTaskStatus_Succeeded
	case "failed":
		return commonpb.ProjectDetectionSubTaskStatus_Failed
	default:
		return commonpb.ProjectDetectionSubTaskStatus_Unknown
	}
}

// ProjectDetectionReviewVerdictString 将项目图像检测人工复核结论枚举转换为字符串
func ProjectDetectionReviewVerdictString(verdict commonpb.ProjectDetectionReviewVerdict_Value) string {
	switch verdict {
	case commonpb.ProjectDetectionReviewVerdict_Accurate:
		return "accurate"
	case commonpb.ProjectDetectionReviewVerdict_Inaccurate:
		return "inaccurate"
	default:
		return ""
	}
}

// ParseProjectDetectionReviewVerdict 将存储值转换为项目图像检测人工复核结论枚举
func ParseProjectDetectionReviewVerdict(value string) commonpb.ProjectDetectionReviewVerdict_Value {
	switch strings.TrimSpace(value) {
	case "accurate":
		return commonpb.ProjectDetectionReviewVerdict_Accurate
	case "inaccurate":
		return commonpb.ProjectDetectionReviewVerdict_Inaccurate
	default:
		return commonpb.ProjectDetectionReviewVerdict_Unknown
	}
}

// ProjectReportStatusString 将项目评估报告状态枚举转换为字符串
func ProjectReportStatusString(status commonpb.ProjectReportStatus_Value) string {
	switch status {
	case commonpb.ProjectReportStatus_Pending:
		return "pending"
	case commonpb.ProjectReportStatus_Succeeded:
		return "succeeded"
	case commonpb.ProjectReportStatus_Failed:
		return "failed"
	default:
		return ""
	}
}

// ParseProjectReportStatus 将存储值转换为项目评估报告状态枚举
func ParseProjectReportStatus(value string) commonpb.ProjectReportStatus_Value {
	switch strings.TrimSpace(value) {
	case "pending":
		return commonpb.ProjectReportStatus_Pending
	case "succeeded":
		return commonpb.ProjectReportStatus_Succeeded
	case "failed":
		return commonpb.ProjectReportStatus_Failed
	default:
		return commonpb.ProjectReportStatus_Unknown
	}
}
