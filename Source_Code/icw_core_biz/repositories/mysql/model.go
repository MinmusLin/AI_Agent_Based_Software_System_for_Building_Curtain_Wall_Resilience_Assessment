package mysql

import (
	"database/sql"
	"errors"
	"time"
)

var (
	// ErrRefreshTokenNotReplaceable 旧 Refresh Token 不存在或已吊销，不能继续完成轮换
	ErrRefreshTokenNotReplaceable = errors.New("refresh token not replaceable")
)

// EmailSendStatus 邮件发送状态
type EmailSendStatus string

const (
	// EmailSendStatusSuccess 邮件发送成功
	EmailSendStatusSuccess EmailSendStatus = "success"
	// EmailSendStatusFailed 邮件发送失败
	EmailSendStatusFailed EmailSendStatus = "failed"
)

// ProjectStatus 项目状态
type ProjectStatus string

const (
	// ProjectStatusActive 进行中
	ProjectStatusActive ProjectStatus = "active"
	// ProjectStatusCompleted 已完成
	ProjectStatusCompleted ProjectStatus = "completed"
	// ProjectStatusDeleted 已删除
	ProjectStatusDeleted ProjectStatus = "deleted"
)

// ProjectProgress 项目进度
type ProjectProgress int8

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

// UserRecord 用户记录
type UserRecord struct {
	Id           uint64
	Email        string
	PasswordHash string
	Name         string
}

// RefreshTokenRecord Refresh Token 记录
type RefreshTokenRecord struct {
	TokenId   string
	UserId    uint64
	ExpiresAt time.Time
	RevokedAt sql.NullTime
}

// ProjectRecord 项目记录
type ProjectRecord struct {
	Id                  uint64
	UserId              uint64
	Name                string
	BuildingName        string
	BuildingLocation    string
	BuiltYear           sql.NullInt64
	BuildingDescription string
	KnownIssues         string
	AssessmentGoal      string
	Progress            ProjectProgress
	Status              string
	CreatedAt           time.Time
	UpdatedAt           time.Time
}
