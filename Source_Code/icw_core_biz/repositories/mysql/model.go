package mysql

import (
	"context"
	"database/sql"
	"errors"
	"time"

	authConsts "icw_core_biz/internal/services/auth/consts"
	projectConsts "icw_core_biz/internal/services/project/consts"
	"icw_core_biz/pkg/dto"
	"icw_core_biz/pkg/dto/project"
	"icw_core_biz/pkg/rpc_err"
	"icw_core_biz/repositories/minio"
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

// UserRecordToDTO 将 MySQL 数据模型转换为 RPC 数据模型
func UserRecordToDTO(user *UserRecord) *dto.User {
	if user == nil {
		return nil
	}
	return &dto.User{
		Id:    user.Id,
		Email: user.Email,
		Name:  user.Name,
	}
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

// ProjectRecordToDTO 将 MySQL 数据模型转换为 RPC 数据模型
func ProjectRecordToDTO(record *ProjectRecord) *project.Project {
	if record == nil {
		return nil
	}
	return &project.Project{
		Id:                  record.Id,
		Name:                record.Name,
		BuildingName:        record.BuildingName,
		BuildingLocation:    record.BuildingLocation,
		BuiltYear:           nullInt64ToUint16(record.BuiltYear),
		BuildingDescription: record.BuildingDescription.String,
		KnownIssues:         record.KnownIssues.String,
		AssessmentGoal:      record.AssessmentGoal.String,
		ThumbnailURL:        "",
		Progress:            record.Progress.Uint8(),
		CreatedAt:           timeToString(record.CreatedAt),
		UpdatedAt:           timeToString(record.UpdatedAt),
	}
}

// ProjectRecordToDTOWithThumbnail 将 MySQL 数据模型转换为 RPC 数据模型，并添加项目缩略图下载预签名 URL
func ProjectRecordToDTOWithThumbnail(ctx context.Context, repo *minio.Repository, record *ProjectRecord, ttl time.Duration) (*project.Project, error) {
	item := ProjectRecordToDTO(record)
	if item == nil {
		return nil, nil
	}
	thumbnailURL, err := minio.PresignProjectThumbnailURL(ctx, repo, record.Id, ttl)
	if err != nil {
		return nil, err
	}
	item.ThumbnailURL = thumbnailURL
	return item, nil
}

// ProjectRecordsToListItemsDTO 将 MySQL 数据模型转换为 RPC 数据模型
func ProjectRecordsToListItemsDTO(records []*ProjectRecord) []*project.ProjectListItem {
	if records == nil {
		return make([]*project.ProjectListItem, 0)
	}
	items := make([]*project.ProjectListItem, 0, len(records))
	for _, record := range records {
		if record == nil {
			continue
		}
		items = append(items, &project.ProjectListItem{
			Id:               record.Id,
			Name:             record.Name,
			BuildingName:     record.BuildingName,
			BuildingLocation: record.BuildingLocation,
			ThumbnailURL:     "",
			Progress:         record.Progress.Uint8(),
			CreatedAt:        timeToString(record.CreatedAt),
		})
	}
	return items
}

// ProjectRecordsToListItemsDTOWithThumbnail 将 MySQL 数据模型转换为 RPC 数据模型，并添加项目缩略图下载预签名 URL
func ProjectRecordsToListItemsDTOWithThumbnail(ctx context.Context, repo *minio.Repository, records []*ProjectRecord, ttl time.Duration) ([]*project.ProjectListItem, error) {
	items := ProjectRecordsToListItemsDTO(records)
	if items == nil {
		return nil, nil
	}
	for _, item := range items {
		if item == nil {
			continue
		}
		thumbnailURL, err := minio.PresignProjectThumbnailURL(ctx, repo, item.Id, ttl)
		if err != nil {
			return nil, err
		}
		item.ThumbnailURL = thumbnailURL
	}
	return items, nil
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

// ProjectGroupRecordToDTO 将 MySQL 数据模型转换为 RPC 数据模型
func ProjectGroupRecordToDTO(record *ProjectGroupRecord, images []*project.ProjectImage) *project.ProjectGroup {
	if record == nil {
		return nil
	}
	if images == nil {
		images = make([]*project.ProjectImage, 0)
	}
	return &project.ProjectGroup{
		Id:        record.Id,
		Name:      record.Name,
		SortOrder: record.SortOrder,
		Images:    images,
	}
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

// ProjectImageRecordToDTO 将 MySQL 数据模型转换为 RPC 数据模型
func ProjectImageRecordToDTO(ctx context.Context, repo *minio.Repository, record *ProjectImageRecord, ttl time.Duration) (*project.ProjectImage, error) {
	if record == nil {
		return nil, nil
	}
	item := &project.ProjectImage{
		Uuid:        record.Uuid,
		FileName:    record.FileName,
		ContentType: record.ContentType,
		SizeBytes:   record.SizeBytes,
		Width:       record.Width,
		Height:      record.Height,
		Metadata:    record.Metadata,
		Status:      record.Status.String(),
		UploadedAt:  timeToString(record.UploadedAt.Time),
		CreatedAt:   timeToString(record.CreatedAt),
	}
	if record.Status != projectConsts.ProjectImageStatusUploaded {
		return item, nil
	}
	thumbnailKey, err := minio.GenProjectImageThumbnailKey(record.ProjectId, record.Uuid)
	if err != nil {
		return nil, rpc_err.BadRequestDefault(err.Error())
	}
	item.ThumbnailURL, err = repo.PresignGetObject(ctx, thumbnailKey, ttl)
	if err != nil {
		return nil, err
	}
	return item, nil
}
