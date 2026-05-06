package enum

import (
	"strings"

	"icw_common/gen/core/biz"
)

// EmailCodeSceneString 将邮箱验证码业务场景枚举转换为字符串
func EmailCodeSceneString(scene bizpb.EmailCodeScene) string {
	switch scene {
	case bizpb.EmailCodeScene_EMAIL_CODE_SCENE_REGISTER:
		return "register"
	case bizpb.EmailCodeScene_EMAIL_CODE_SCENE_LOGIN:
		return "login"
	case bizpb.EmailCodeScene_EMAIL_CODE_SCENE_RESET:
		return "reset"
	default:
		return ""
	}
}

// ParseEmailCodeScene 将存储值转换为邮箱验证码业务场景枚举
func ParseEmailCodeScene(value string) bizpb.EmailCodeScene {
	switch strings.TrimSpace(value) {
	case "register":
		return bizpb.EmailCodeScene_EMAIL_CODE_SCENE_REGISTER
	case "login":
		return bizpb.EmailCodeScene_EMAIL_CODE_SCENE_LOGIN
	case "reset":
		return bizpb.EmailCodeScene_EMAIL_CODE_SCENE_RESET
	default:
		return bizpb.EmailCodeScene_EMAIL_CODE_SCENE_UNKNOWN
	}
}

// LoginSceneString 将登录方式枚举转换为字符串
func LoginSceneString(scene bizpb.LoginScene) string {
	switch scene {
	case bizpb.LoginScene_LOGIN_SCENE_PASSWORD:
		return "password"
	case bizpb.LoginScene_LOGIN_SCENE_EMAIL:
		return "email"
	default:
		return ""
	}
}

// ParseLoginScene 将存储值转换为登录方式枚举
func ParseLoginScene(value string) bizpb.LoginScene {
	switch strings.TrimSpace(value) {
	case "password":
		return bizpb.LoginScene_LOGIN_SCENE_PASSWORD
	case "email":
		return bizpb.LoginScene_LOGIN_SCENE_EMAIL
	default:
		return bizpb.LoginScene_LOGIN_SCENE_UNKNOWN
	}
}

// EmailSendStatusString 将邮件发送状态枚举转换为字符串
func EmailSendStatusString(status bizpb.EmailSendStatus) string {
	switch status {
	case bizpb.EmailSendStatus_EMAIL_SEND_STATUS_SUCCESS:
		return "success"
	case bizpb.EmailSendStatus_EMAIL_SEND_STATUS_FAILED:
		return "failed"
	default:
		return ""
	}
}

// ParseEmailSendStatus 将存储值转换为邮件发送状态枚举
func ParseEmailSendStatus(value string) bizpb.EmailSendStatus {
	switch strings.TrimSpace(value) {
	case "success":
		return bizpb.EmailSendStatus_EMAIL_SEND_STATUS_SUCCESS
	case "failed":
		return bizpb.EmailSendStatus_EMAIL_SEND_STATUS_FAILED
	default:
		return bizpb.EmailSendStatus_EMAIL_SEND_STATUS_UNKNOWN
	}
}

// ProjectStatusString 将项目状态枚举转换为字符串
func ProjectStatusString(status bizpb.ProjectStatus) string {
	switch status {
	case bizpb.ProjectStatus_PROJECT_STATUS_ACTIVE:
		return "active"
	case bizpb.ProjectStatus_PROJECT_STATUS_COMPLETED:
		return "completed"
	case bizpb.ProjectStatus_PROJECT_STATUS_DELETED:
		return "deleted"
	default:
		return ""
	}
}

// ParseProjectStatus 将存储值转换为项目状态枚举
func ParseProjectStatus(value string) bizpb.ProjectStatus {
	switch strings.TrimSpace(value) {
	case "active":
		return bizpb.ProjectStatus_PROJECT_STATUS_ACTIVE
	case "completed":
		return bizpb.ProjectStatus_PROJECT_STATUS_COMPLETED
	case "deleted":
		return bizpb.ProjectStatus_PROJECT_STATUS_DELETED
	default:
		return bizpb.ProjectStatus_PROJECT_STATUS_UNKNOWN
	}
}

// ProjectImageStatusString 将项目图像状态枚举转换为字符串
func ProjectImageStatusString(status bizpb.ProjectImageStatus) string {
	switch status {
	case bizpb.ProjectImageStatus_PROJECT_IMAGE_STATUS_PENDING:
		return "pending"
	case bizpb.ProjectImageStatus_PROJECT_IMAGE_STATUS_UPLOADED:
		return "uploaded"
	case bizpb.ProjectImageStatus_PROJECT_IMAGE_STATUS_FAILED:
		return "failed"
	default:
		return ""
	}
}

// ParseProjectImageStatus 将存储值转换为项目图像状态枚举
func ParseProjectImageStatus(value string) bizpb.ProjectImageStatus {
	switch strings.TrimSpace(value) {
	case "pending":
		return bizpb.ProjectImageStatus_PROJECT_IMAGE_STATUS_PENDING
	case "uploaded":
		return bizpb.ProjectImageStatus_PROJECT_IMAGE_STATUS_UPLOADED
	case "failed":
		return bizpb.ProjectImageStatus_PROJECT_IMAGE_STATUS_FAILED
	default:
		return bizpb.ProjectImageStatus_PROJECT_IMAGE_STATUS_UNKNOWN
	}
}

// ProjectDetectionTaskStatusString 将项目图像检测主任务状态枚举转换为字符串
func ProjectDetectionTaskStatusString(status bizpb.ProjectDetectionTaskStatus) string {
	switch status {
	case bizpb.ProjectDetectionTaskStatus_PROJECT_DETECTION_TASK_STATUS_PENDING:
		return "pending"
	case bizpb.ProjectDetectionTaskStatus_PROJECT_DETECTION_TASK_STATUS_CLASSIFYING:
		return "classifying"
	case bizpb.ProjectDetectionTaskStatus_PROJECT_DETECTION_TASK_STATUS_DETECTING:
		return "detecting"
	case bizpb.ProjectDetectionTaskStatus_PROJECT_DETECTION_TASK_STATUS_SUMMARIZING:
		return "summarizing"
	case bizpb.ProjectDetectionTaskStatus_PROJECT_DETECTION_TASK_STATUS_SUCCEEDED:
		return "succeeded"
	case bizpb.ProjectDetectionTaskStatus_PROJECT_DETECTION_TASK_STATUS_FAILED:
		return "failed"
	default:
		return ""
	}
}

// ParseProjectDetectionTaskStatus 将存储值转换为项目图像检测主任务状态枚举
func ParseProjectDetectionTaskStatus(value string) bizpb.ProjectDetectionTaskStatus {
	switch strings.TrimSpace(value) {
	case "pending":
		return bizpb.ProjectDetectionTaskStatus_PROJECT_DETECTION_TASK_STATUS_PENDING
	case "classifying":
		return bizpb.ProjectDetectionTaskStatus_PROJECT_DETECTION_TASK_STATUS_CLASSIFYING
	case "detecting":
		return bizpb.ProjectDetectionTaskStatus_PROJECT_DETECTION_TASK_STATUS_DETECTING
	case "summarizing":
		return bizpb.ProjectDetectionTaskStatus_PROJECT_DETECTION_TASK_STATUS_SUMMARIZING
	case "succeeded":
		return bizpb.ProjectDetectionTaskStatus_PROJECT_DETECTION_TASK_STATUS_SUCCEEDED
	case "failed":
		return bizpb.ProjectDetectionTaskStatus_PROJECT_DETECTION_TASK_STATUS_FAILED
	default:
		return bizpb.ProjectDetectionTaskStatus_PROJECT_DETECTION_TASK_STATUS_UNKNOWN
	}
}

// ProjectDetectionSubTaskStatusString 将项目图像检测子任务状态枚举转换为字符串
func ProjectDetectionSubTaskStatusString(status bizpb.ProjectDetectionSubTaskStatus) string {
	switch status {
	case bizpb.ProjectDetectionSubTaskStatus_PROJECT_DETECTION_SUB_TASK_STATUS_PENDING:
		return "pending"
	case bizpb.ProjectDetectionSubTaskStatus_PROJECT_DETECTION_SUB_TASK_STATUS_RUNNING:
		return "running"
	case bizpb.ProjectDetectionSubTaskStatus_PROJECT_DETECTION_SUB_TASK_STATUS_SUCCEEDED:
		return "succeeded"
	case bizpb.ProjectDetectionSubTaskStatus_PROJECT_DETECTION_SUB_TASK_STATUS_FAILED:
		return "failed"
	default:
		return ""
	}
}

// ParseProjectDetectionSubTaskStatus 将存储值转换为项目图像检测子任务状态枚举
func ParseProjectDetectionSubTaskStatus(value string) bizpb.ProjectDetectionSubTaskStatus {
	switch strings.TrimSpace(value) {
	case "pending":
		return bizpb.ProjectDetectionSubTaskStatus_PROJECT_DETECTION_SUB_TASK_STATUS_PENDING
	case "running":
		return bizpb.ProjectDetectionSubTaskStatus_PROJECT_DETECTION_SUB_TASK_STATUS_RUNNING
	case "succeeded":
		return bizpb.ProjectDetectionSubTaskStatus_PROJECT_DETECTION_SUB_TASK_STATUS_SUCCEEDED
	case "failed":
		return bizpb.ProjectDetectionSubTaskStatus_PROJECT_DETECTION_SUB_TASK_STATUS_FAILED
	default:
		return bizpb.ProjectDetectionSubTaskStatus_PROJECT_DETECTION_SUB_TASK_STATUS_UNKNOWN
	}
}
