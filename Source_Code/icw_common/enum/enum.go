package enum

import (
	"strconv"
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

// AvatarTypeString 将用户头像类型枚举转换为字符串
func AvatarTypeString(avatarType commonpb.AvatarType_Value) string {
	switch avatarType {
	case commonpb.AvatarType_Custom:
		return "custom"
	case commonpb.AvatarType_Default:
		return "default"
	case commonpb.AvatarType_None:
		return "none"
	default:
		return ""
	}
}

// ParseAvatarType 将存储值转换为用户头像类型枚举
func ParseAvatarType(value string) commonpb.AvatarType_Value {
	switch strings.TrimSpace(value) {
	case "custom":
		return commonpb.AvatarType_Custom
	case "default":
		return commonpb.AvatarType_Default
	case "none":
		return commonpb.AvatarType_None
	default:
		return commonpb.AvatarType_Unknown
	}
}

// ProjectEventTypeString 将项目事件类型枚举转换为字符串
func ProjectEventTypeString(eventType commonpb.ProjectEventType_Value) string {
	switch eventType {
	case commonpb.ProjectEventType_ImageStatusChanged:
		return "project_image_status_changed"
	case commonpb.ProjectEventType_DetectionTaskStatusChanged:
		return "project_detection_task_status_changed"
	case commonpb.ProjectEventType_ReportStatusChanged:
		return "project_report_status_changed"
	default:
		return ""
	}
}

// ParseProjectEventType 将存储值转换为项目事件类型枚举
func ParseProjectEventType(value string) commonpb.ProjectEventType_Value {
	switch strings.TrimSpace(value) {
	case "project_image_status_changed":
		return commonpb.ProjectEventType_ImageStatusChanged
	case "project_detection_task_status_changed":
		return commonpb.ProjectEventType_DetectionTaskStatusChanged
	case "project_report_status_changed":
		return commonpb.ProjectEventType_ReportStatusChanged
	default:
		return commonpb.ProjectEventType_Unknown
	}
}

// SocketScopeString 将 WebSocket 连接范围枚举转换为字符串
func SocketScopeString(scope commonpb.SocketScope_Value) string {
	switch scope {
	case commonpb.SocketScope_ProjectAssets:
		return "ws_project_assets"
	case commonpb.SocketScope_ProjectDetection:
		return "ws_project_detection"
	case commonpb.SocketScope_ProjectReport:
		return "ws_project_report"
	default:
		return ""
	}
}

// ParseSocketScope 将存储值转换为 WebSocket 连接范围枚举
func ParseSocketScope(value string) commonpb.SocketScope_Value {
	value = strings.TrimSpace(value)
	if numericValue, err := strconv.Atoi(value); err == nil {
		switch scope := commonpb.SocketScope_Value(numericValue); scope {
		case commonpb.SocketScope_ProjectAssets,
			commonpb.SocketScope_ProjectDetection,
			commonpb.SocketScope_ProjectReport:
			return scope
		default:
			return commonpb.SocketScope_Unknown
		}
	}
	switch value {
	case "ws_project_assets":
		return commonpb.SocketScope_ProjectAssets
	case "ws_project_detection":
		return commonpb.SocketScope_ProjectDetection
	case "ws_project_report":
		return commonpb.SocketScope_ProjectReport
	default:
		return commonpb.SocketScope_Unknown
	}
}

// TaskStatusString 将任务状态枚举转换为字符串
func TaskStatusString(status commonpb.TaskStatus_Value) string {
	switch status {
	case commonpb.TaskStatus_Succeeded:
		return "succeeded"
	case commonpb.TaskStatus_Failed:
		return "failed"
	default:
		return ""
	}
}

// ParseTaskStatus 将存储值转换为任务状态枚举
func ParseTaskStatus(value string) commonpb.TaskStatus_Value {
	switch strings.TrimSpace(value) {
	case "succeeded":
		return commonpb.TaskStatus_Succeeded
	case "failed":
		return commonpb.TaskStatus_Failed
	default:
		return commonpb.TaskStatus_Unknown
	}
}

// DetectionNodeCodeString 将项目图像检测节点编码枚举转换为字符串
func DetectionNodeCodeString(code commonpb.DetectionNodeCode_Value) string {
	switch code {
	case commonpb.DetectionNodeCode_Classification:
		return "classification"
	case commonpb.DetectionNodeCode_Reasoning:
		return "reasoning"
	case commonpb.DetectionNodeCode_Summary:
		return "summary"
	default:
		return ""
	}
}

// ParseDetectionNodeCode 将存储值转换为项目图像检测节点编码枚举
func ParseDetectionNodeCode(value string) commonpb.DetectionNodeCode_Value {
	switch strings.TrimSpace(value) {
	case "classification":
		return commonpb.DetectionNodeCode_Classification
	case "reasoning":
		return commonpb.DetectionNodeCode_Reasoning
	case "summary":
		return commonpb.DetectionNodeCode_Summary
	default:
		return commonpb.DetectionNodeCode_Unknown
	}
}

// DetectionTaskCodeString 将原子检测能力代码枚举转换为字符串
func DetectionTaskCodeString(code commonpb.DetectionTaskCode_Value) string {
	switch code {
	case commonpb.DetectionTaskCode_Corrosion:
		return "corrosion"
	case commonpb.DetectionTaskCode_Crack:
		return "crack"
	case commonpb.DetectionTaskCode_Stain:
		return "stain"
	case commonpb.DetectionTaskCode_Flatness:
		return "flatness"
	case commonpb.DetectionTaskCode_Spalling:
		return "spalling"
	default:
		return ""
	}
}

// ParseDetectionTaskCode 将存储值转换为原子检测能力代码枚举
func ParseDetectionTaskCode(value string) commonpb.DetectionTaskCode_Value {
	switch strings.TrimSpace(value) {
	case "corrosion":
		return commonpb.DetectionTaskCode_Corrosion
	case "crack":
		return commonpb.DetectionTaskCode_Crack
	case "stain":
		return commonpb.DetectionTaskCode_Stain
	case "flatness":
		return commonpb.DetectionTaskCode_Flatness
	case "spalling":
		return commonpb.DetectionTaskCode_Spalling
	default:
		return commonpb.DetectionTaskCode_Unknown
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
