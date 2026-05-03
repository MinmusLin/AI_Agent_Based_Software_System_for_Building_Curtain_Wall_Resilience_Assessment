package mysql

import (
	"database/sql"
	"errors"
	"time"

	authConsts "icw_core_biz/internal/services/auth/consts"
	projectConsts "icw_core_biz/internal/services/project/consts"
	"icw_core_biz/pkg/dto"
)

var (
	// ErrRefreshTokenNotReplaceable 旧 Refresh Token 不存在或已吊销，不能继续完成轮换
	ErrRefreshTokenNotReplaceable = errors.New("refresh token not replaceable")
)

// UserRecord 用户记录
type UserRecord struct {
	Id           uint64
	Email        string
	PasswordHash string
	Name         string
	LastLoginAt  sql.NullTime
	CreatedAt    time.Time
	UpdatedAt    time.Time
}

// RefreshTokenRecord Refresh Token 记录
type RefreshTokenRecord struct {
	Id                uint64
	TokenId           string
	UserId            uint64
	TokenHash         string
	ExpiresAt         time.Time
	RevokedAt         sql.NullTime
	CreatedAt         time.Time
	UpdatedAt         time.Time
	ReplacedByTokenId sql.NullString
}

// EmailSendLogRecord 邮件发送记录
type EmailSendLogRecord struct {
	Id            uint64
	ReceiverEmail string
	SenderEmail   string
	Scene         string
	EmailCode     string
	Status        authConsts.EmailSendStatus
	ErrorMessage  sql.NullString
	CreatedAt     time.Time
	UpdatedAt     time.Time
}

// ProjectRecord 项目记录
type ProjectRecord struct {
	Id                  uint64
	UserId              uint64
	Name                string
	BuildingName        string
	BuildingLocation    string
	BuiltYear           sql.NullInt64
	BuildingDescription sql.NullString
	KnownIssues         sql.NullString
	AssessmentGoal      sql.NullString
	Progress            dto.ProjectProgress
	Status              dto.ProjectStatus
	CreatedAt           time.Time
	UpdatedAt           time.Time
}

// ProjectGroupRecord 项目图像组记录
type ProjectGroupRecord struct {
	Id        uint64
	ProjectId uint64
	UserId    uint64
	Name      string
	SortOrder string
	CreatedAt time.Time
	UpdatedAt time.Time
}

// ProjectImageRecord 项目图像记录
type ProjectImageRecord struct {
	Id          uint64
	GroupId     uint64
	ProjectId   uint64
	UserId      uint64
	Uuid        string
	FileName    string
	ContentType string
	SizeBytes   uint64
	Width       uint32
	Height      uint32
	Metadata    string
	Status      projectConsts.ProjectImageStatus
	UploadedAt  sql.NullTime
	CreatedAt   time.Time
	UpdatedAt   time.Time
}
