package enum

import (
	"strings"

	"icw_common/gen/core/biz"
)

// EmailCodeSceneString 将邮箱验证码业务场景枚举转换为字符串
func EmailCodeSceneString(scene bizpb.EmailCodeScene_Value) string {
	switch scene {
	case bizpb.EmailCodeScene_Register:
		return "register"
	case bizpb.EmailCodeScene_Login:
		return "login"
	case bizpb.EmailCodeScene_Reset:
		return "reset"
	default:
		return ""
	}
}

// ParseEmailCodeScene 将存储值转换为邮箱验证码业务场景枚举
func ParseEmailCodeScene(value string) bizpb.EmailCodeScene_Value {
	switch strings.TrimSpace(value) {
	case "register":
		return bizpb.EmailCodeScene_Register
	case "login":
		return bizpb.EmailCodeScene_Login
	case "reset":
		return bizpb.EmailCodeScene_Reset
	default:
		return bizpb.EmailCodeScene_Unknown
	}
}

// LoginSceneString 将登录方式枚举转换为字符串
func LoginSceneString(scene bizpb.LoginScene_Value) string {
	switch scene {
	case bizpb.LoginScene_Password:
		return "password"
	case bizpb.LoginScene_Email:
		return "email"
	default:
		return ""
	}
}

// ParseLoginScene 将存储值转换为登录方式枚举
func ParseLoginScene(value string) bizpb.LoginScene_Value {
	switch strings.TrimSpace(value) {
	case "password":
		return bizpb.LoginScene_Password
	case "email":
		return bizpb.LoginScene_Email
	default:
		return bizpb.LoginScene_Unknown
	}
}

// EmailSendStatusString 将邮件发送状态枚举转换为字符串
func EmailSendStatusString(status bizpb.EmailSendStatus_Value) string {
	switch status {
	case bizpb.EmailSendStatus_Success:
		return "success"
	case bizpb.EmailSendStatus_Failed:
		return "failed"
	default:
		return ""
	}
}

// ParseEmailSendStatus 将存储值转换为邮件发送状态枚举
func ParseEmailSendStatus(value string) bizpb.EmailSendStatus_Value {
	switch strings.TrimSpace(value) {
	case "success":
		return bizpb.EmailSendStatus_Success
	case "failed":
		return bizpb.EmailSendStatus_Failed
	default:
		return bizpb.EmailSendStatus_Unknown
	}
}

// ProjectStatusString 将项目状态枚举转换为字符串
func ProjectStatusString(status bizpb.ProjectStatus_Value) string {
	switch status {
	case bizpb.ProjectStatus_Active:
		return "active"
	case bizpb.ProjectStatus_Completed:
		return "completed"
	case bizpb.ProjectStatus_Deleted:
		return "deleted"
	default:
		return ""
	}
}

// ParseProjectStatus 将存储值转换为项目状态枚举
func ParseProjectStatus(value string) bizpb.ProjectStatus_Value {
	switch strings.TrimSpace(value) {
	case "active":
		return bizpb.ProjectStatus_Active
	case "completed":
		return bizpb.ProjectStatus_Completed
	case "deleted":
		return bizpb.ProjectStatus_Deleted
	default:
		return bizpb.ProjectStatus_Unknown
	}
}

// ProjectImageStatusString 将项目图像状态枚举转换为字符串
func ProjectImageStatusString(status bizpb.ProjectImageStatus_Value) string {
	switch status {
	case bizpb.ProjectImageStatus_Pending:
		return "pending"
	case bizpb.ProjectImageStatus_Uploaded:
		return "uploaded"
	case bizpb.ProjectImageStatus_Failed:
		return "failed"
	default:
		return ""
	}
}

// ParseProjectImageStatus 将存储值转换为项目图像状态枚举
func ParseProjectImageStatus(value string) bizpb.ProjectImageStatus_Value {
	switch strings.TrimSpace(value) {
	case "pending":
		return bizpb.ProjectImageStatus_Pending
	case "uploaded":
		return bizpb.ProjectImageStatus_Uploaded
	case "failed":
		return bizpb.ProjectImageStatus_Failed
	default:
		return bizpb.ProjectImageStatus_Unknown
	}
}

// ProjectDetectionTaskStatusString 将项目图像检测主任务状态枚举转换为字符串
func ProjectDetectionTaskStatusString(status bizpb.ProjectDetectionTaskStatus_Value) string {
	switch status {
	case bizpb.ProjectDetectionTaskStatus_Pending:
		return "pending"
	case bizpb.ProjectDetectionTaskStatus_Classifying:
		return "classifying"
	case bizpb.ProjectDetectionTaskStatus_Detecting:
		return "detecting"
	case bizpb.ProjectDetectionTaskStatus_Summarizing:
		return "summarizing"
	case bizpb.ProjectDetectionTaskStatus_Succeeded:
		return "succeeded"
	case bizpb.ProjectDetectionTaskStatus_Failed:
		return "failed"
	default:
		return ""
	}
}

// ParseProjectDetectionTaskStatus 将存储值转换为项目图像检测主任务状态枚举
func ParseProjectDetectionTaskStatus(value string) bizpb.ProjectDetectionTaskStatus_Value {
	switch strings.TrimSpace(value) {
	case "pending":
		return bizpb.ProjectDetectionTaskStatus_Pending
	case "classifying":
		return bizpb.ProjectDetectionTaskStatus_Classifying
	case "detecting":
		return bizpb.ProjectDetectionTaskStatus_Detecting
	case "summarizing":
		return bizpb.ProjectDetectionTaskStatus_Summarizing
	case "succeeded":
		return bizpb.ProjectDetectionTaskStatus_Succeeded
	case "failed":
		return bizpb.ProjectDetectionTaskStatus_Failed
	default:
		return bizpb.ProjectDetectionTaskStatus_Unknown
	}
}

// ProjectDetectionSubTaskStatusString 将项目图像检测子任务状态枚举转换为字符串
func ProjectDetectionSubTaskStatusString(status bizpb.ProjectDetectionSubTaskStatus_Value) string {
	switch status {
	case bizpb.ProjectDetectionSubTaskStatus_Pending:
		return "pending"
	case bizpb.ProjectDetectionSubTaskStatus_Running:
		return "running"
	case bizpb.ProjectDetectionSubTaskStatus_Succeeded:
		return "succeeded"
	case bizpb.ProjectDetectionSubTaskStatus_Failed:
		return "failed"
	default:
		return ""
	}
}

// ParseProjectDetectionSubTaskStatus 将存储值转换为项目图像检测子任务状态枚举
func ParseProjectDetectionSubTaskStatus(value string) bizpb.ProjectDetectionSubTaskStatus_Value {
	switch strings.TrimSpace(value) {
	case "pending":
		return bizpb.ProjectDetectionSubTaskStatus_Pending
	case "running":
		return bizpb.ProjectDetectionSubTaskStatus_Running
	case "succeeded":
		return bizpb.ProjectDetectionSubTaskStatus_Succeeded
	case "failed":
		return bizpb.ProjectDetectionSubTaskStatus_Failed
	default:
		return bizpb.ProjectDetectionSubTaskStatus_Unknown
	}
}
