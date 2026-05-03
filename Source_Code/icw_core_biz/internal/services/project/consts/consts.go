package consts

import (
	"strings"
)

const (
	// DefaultProjectName 默认项目名称
	DefaultProjectName = "新项目"
	// DefaultProjectGroupName 默认图像组名称
	DefaultProjectGroupName = "默认图像组"
	// DefaultNewProjectGroupName 默认新图像组名称
	DefaultNewProjectGroupName = "新图像组"
)

const (
	// ProjectNameMaxLength 项目名称最大字符数
	ProjectNameMaxLength = 32
	// ProjectBuildingNameMaxLength 建筑名称最大字符数
	ProjectBuildingNameMaxLength = 32
	// ProjectBuildingLocationMaxLength 建筑地址最大字符数
	ProjectBuildingLocationMaxLength = 128
	// ProjectBuildingDescriptionMaxLength 建筑描述最大字符数
	ProjectBuildingDescriptionMaxLength = 5000
	// ProjectKnownIssuesMaxLength 已知问题或人工先验描述最大字符数
	ProjectKnownIssuesMaxLength = 5000
	// ProjectAssessmentGoalMaxLength 评估目标或重点关注方向最大字符数
	ProjectAssessmentGoalMaxLength = 5000
	// ProjectGroupNameMaxLength 图像组名称最大字符数
	ProjectGroupNameMaxLength = 32
	// ProjectImageFileNameMaxLength 图像文件名最大字符数
	ProjectImageFileNameMaxLength = 255
)

const (
	// ThumbnailContentType 项目缩略图 MIME 类型
	ThumbnailContentType = "image/png"
	// ProjectImageContentType 项目图像 MIME 类型
	ProjectImageContentType = "image/png"
)

// ProjectImageStatus 项目图像状态枚举
type ProjectImageStatus string

const (
	// ProjectImageStatusPending 图像上传中
	ProjectImageStatusPending ProjectImageStatus = "pending"
	// ProjectImageStatusUploaded 图像上传成功
	ProjectImageStatusUploaded ProjectImageStatus = "uploaded"
	// ProjectImageStatusFailed 图像上传失败
	ProjectImageStatusFailed ProjectImageStatus = "failed"
)

// String 将项目图像状态枚举转换为字符串
func (s ProjectImageStatus) String() string {
	return string(s)
}

// ParseProjectImageStatus 将外部输入转换为项目图像状态枚举
func ParseProjectImageStatus(value string) ProjectImageStatus {
	switch status := ProjectImageStatus(strings.TrimSpace(value)); status {
	case ProjectImageStatusPending, ProjectImageStatusUploaded, ProjectImageStatusFailed:
		return status
	default:
		return ""
	}
}
