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
	// ThumbnailContentType 项目缩略图 MIME 类型
	ThumbnailContentType = "image/png"
)

// ProjectStatus 项目状态枚举
type ProjectStatus string

const (
	// ProjectStatusActive 进行中
	ProjectStatusActive ProjectStatus = "active"
	// ProjectStatusCompleted 已完成
	ProjectStatusCompleted ProjectStatus = "completed"
	// ProjectStatusDeleted 已删除
	ProjectStatusDeleted ProjectStatus = "deleted"
)

// String 将项目状态枚举转换为字符串
func (s ProjectStatus) String() string {
	return string(s)
}

// ParseProjectStatus 将外部输入转换为项目状态枚举
func ParseProjectStatus(value string) ProjectStatus {
	switch status := ProjectStatus(strings.TrimSpace(value)); status {
	case ProjectStatusActive, ProjectStatusCompleted, ProjectStatusDeleted:
		return status
	default:
		return ""
	}
}

// ProjectProgress 项目进度枚举
type ProjectProgress uint8

const (
	// ProjectProgressInitializationFinished 项目初始化完成，当前项目基础信息阶段
	ProjectProgressInitializationFinished ProjectProgress = 0
	// ProjectProgressProfileFinished 项目基础信息完成，当前图像资产构建阶段
	ProjectProgressProfileFinished ProjectProgress = 1
	// ProjectProgressAssetsFinished 图像资产构建完成，当前 Agent 智能检测阶段
	ProjectProgressAssetsFinished ProjectProgress = 2
	// ProjectProgressDetectFinished Agent 智能检测完成，当前人工复核确认阶段
	ProjectProgressDetectFinished ProjectProgress = 3
	// ProjectProgressReviewFinished 人工复核确认完成，当前评估报告生成阶段
	ProjectProgressReviewFinished ProjectProgress = 4
	// ProjectProgressReportFinished 评估报告生成完成，当前项目已完成
	ProjectProgressReportFinished ProjectProgress = 5
)

// Uint8 将项目进度枚举转换为 uint8
func (p ProjectProgress) Uint8() uint8 {
	return uint8(p)
}

// ParseProjectProgress 将外部输入转换为项目进度枚举
func ParseProjectProgress(value uint8) ProjectProgress {
	switch progress := ProjectProgress(value); progress {
	case ProjectProgressInitializationFinished,
		ProjectProgressProfileFinished,
		ProjectProgressAssetsFinished,
		ProjectProgressDetectFinished,
		ProjectProgressReviewFinished,
		ProjectProgressReportFinished:
		return progress
	default:
		return ProjectProgressInitializationFinished
	}
}

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
