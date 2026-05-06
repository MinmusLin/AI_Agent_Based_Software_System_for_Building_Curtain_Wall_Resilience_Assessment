package mysql

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"icw_common/consts"
	"icw_common/enum"
	"icw_common/gen/core/biz"
	"icw_common/utils"
	"icw_core_biz/repositories/minio"
	"icw_core_biz/repositories/redis"
)

var (
	// ErrRefreshTokenNotReplaceable 旧 Refresh Token 不存在或已吊销，不能继续完成轮换
	ErrRefreshTokenNotReplaceable = errors.New("refresh token not replaceable")
	// ErrProjectGroupCannotDeleteLast 项目应至少存在一个图像组
	ErrProjectGroupCannotDeleteLast = errors.New("project must keep at least one group")
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
func UserRecordToDTO(user *UserRecord) *bizpb.User {
	if user == nil {
		return nil
	}
	return &bizpb.User{
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
	Status        bizpb.EmailSendStatus_Value
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
	Progress            consts.ProjectProgress
	Status              bizpb.ProjectStatus_Value
	CreatedAt           time.Time
	UpdatedAt           time.Time
}

// ProjectRecordToDTO 将 MySQL 数据模型转换为 RPC 数据模型
func ProjectRecordToDTO(record *ProjectRecord) *bizpb.Project {
	if record == nil {
		return nil
	}
	return &bizpb.Project{
		Id:                  record.Id,
		Name:                record.Name,
		BuildingName:        record.BuildingName,
		BuildingLocation:    record.BuildingLocation,
		BuiltYear:           utils.If[uint32](record.BuiltYear.Valid && record.BuiltYear.Int64 >= 0, uint32(record.BuiltYear.Int64), 0),
		BuildingDescription: record.BuildingDescription.String,
		KnownIssues:         record.KnownIssues.String,
		AssessmentGoal:      record.AssessmentGoal.String,
		ThumbnailUrl:        "",
		Progress:            uint32(record.Progress.Uint8()),
		CreatedAt:           timeToString(record.CreatedAt),
		UpdatedAt:           timeToString(record.UpdatedAt),
	}
}

// ProjectRecordToDTOWithThumbnail 将 MySQL 数据模型转换为 RPC 数据模型，并添加项目缩略图下载预签名 URL
func ProjectRecordToDTOWithThumbnail(ctx context.Context, minioRepo *minio.Repository, redisRepo *redis.Repository, record *ProjectRecord, ttl time.Duration) (*bizpb.Project, error) {
	item := ProjectRecordToDTO(record)
	if item == nil {
		return nil, nil
	}
	thumbnailURL, err := minio.PresignProjectThumbnailURL(ctx, minioRepo, redisRepo, record.UserId, record.Id, ttl)
	if err != nil {
		return nil, err
	}
	item.ThumbnailUrl = thumbnailURL
	return item, nil
}

// ProjectRecordsToListItemsDTO 将 MySQL 数据模型转换为 RPC 数据模型
func ProjectRecordsToListItemsDTO(records []*ProjectRecord) []*bizpb.ProjectListItem {
	if records == nil {
		return make([]*bizpb.ProjectListItem, 0)
	}
	items := make([]*bizpb.ProjectListItem, 0, len(records))
	for _, record := range records {
		if record == nil {
			continue
		}
		items = append(items, &bizpb.ProjectListItem{
			Id:               record.Id,
			Name:             record.Name,
			BuildingName:     record.BuildingName,
			BuildingLocation: record.BuildingLocation,
			ThumbnailUrl:     "",
			Progress:         uint32(record.Progress.Uint8()),
			CreatedAt:        timeToString(record.CreatedAt),
		})
	}
	return items
}

// ProjectRecordsToListItemsDTOWithThumbnail 将 MySQL 数据模型转换为 RPC 数据模型，并添加项目缩略图下载预签名 URL
func ProjectRecordsToListItemsDTOWithThumbnail(ctx context.Context, minioRepo *minio.Repository, redisRepo *redis.Repository, records []*ProjectRecord, ttl time.Duration) ([]*bizpb.ProjectListItem, error) {
	items := ProjectRecordsToListItemsDTO(records)
	if items == nil {
		return nil, nil
	}
	itemIndex := 0
	for _, record := range records {
		if record == nil {
			continue
		}
		thumbnailURL, err := minio.PresignProjectThumbnailURL(ctx, minioRepo, redisRepo, record.UserId, record.Id, ttl)
		if err != nil {
			return nil, err
		}
		items[itemIndex].ThumbnailUrl = thumbnailURL
		itemIndex++
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
func ProjectGroupRecordToDTO(ctx context.Context, minioRepo *minio.Repository, redisRepo *redis.Repository, record *ProjectGroupRecord, images []*ProjectImageRecord, ttl time.Duration) (*bizpb.ProjectGroup, error) {
	if record == nil {
		return nil, nil
	}
	imageItems, err := ProjectImageRecordsToDTO(ctx, minioRepo, redisRepo, images, ttl)
	if err != nil {
		return nil, err
	}
	return &bizpb.ProjectGroup{
		Id:        record.Id,
		Name:      record.Name,
		SortOrder: record.SortOrder,
		Images:    imageItems,
	}, nil
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
	Status      bizpb.ProjectImageStatus_Value
	UploadedAt  sql.NullTime
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

// ProjectImageCreateRecord 项目图像创建记录
type ProjectImageCreateRecord struct {
	ImageUuid   string
	FileName    string
	ContentType string
	SizeBytes   uint64
	Width       uint32
	Height      uint32
	Metadata    string
}

// ProjectImageRecordToDTO 将 MySQL 数据模型转换为 RPC 数据模型
func ProjectImageRecordToDTO(ctx context.Context, minioRepo *minio.Repository, redisRepo *redis.Repository, record *ProjectImageRecord, ttl time.Duration) (*bizpb.ProjectImage, error) {
	if record == nil {
		return nil, nil
	}
	item := &bizpb.ProjectImage{
		Uuid:        record.Uuid,
		FileName:    record.FileName,
		ContentType: record.ContentType,
		SizeBytes:   record.SizeBytes,
		Width:       record.Width,
		Height:      record.Height,
		Metadata:    record.Metadata,
		Status:      enum.ProjectImageStatusString(record.Status),
		UploadedAt:  timeToString(record.UploadedAt.Time),
		CreatedAt:   timeToString(record.CreatedAt),
	}
	if record.Status != bizpb.ProjectImageStatus_Uploaded {
		return item, nil
	}
	var err error
	item.ThumbnailUrl, err = minio.PresignProjectImageThumbnailURL(ctx, minioRepo, redisRepo, record.UserId, record.ProjectId, record.Uuid, ttl)
	if err != nil {
		return nil, err
	}
	return item, nil
}

// ProjectImageRecordsToDTO 将 MySQL 数据模型转换为 RPC 数据模型
func ProjectImageRecordsToDTO(ctx context.Context, minioRepo *minio.Repository, redisRepo *redis.Repository, records []*ProjectImageRecord, ttl time.Duration) ([]*bizpb.ProjectImage, error) {
	if records == nil {
		return make([]*bizpb.ProjectImage, 0), nil
	}
	items := make([]*bizpb.ProjectImage, 0, len(records))
	for _, record := range records {
		if record == nil {
			continue
		}
		item, err := ProjectImageRecordToDTO(ctx, minioRepo, redisRepo, record, ttl)
		if err != nil {
			return nil, err
		}
		items = append(items, item)
	}
	return items, nil
}

// ProjectAssetsReadyStats 项目图像状态校验
type ProjectAssetsReadyStats struct {
	PendingImageCount  uint64
	UploadedImageCount uint64
	FailedImageCount   uint64
	EmptyGroupCount    uint64
}
