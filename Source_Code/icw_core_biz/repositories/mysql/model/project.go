package model

import (
	"context"
	"database/sql"
	"time"

	"icw_common/gen/core/biz"

	"icw_core_biz/repositories/minio"
	"icw_core_biz/repositories/mysql/utils"
	"icw_core_biz/repositories/redis"
)

// ProjectAssetsReadyStats 项目图像状态校验
type ProjectAssetsReadyStats struct {
	PendingImageCount  uint64
	UploadedImageCount uint64
	FailedImageCount   uint64
	EmptyGroupCount    uint64
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
	Progress            bizpb.ProjectProgress_Value
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
		BuiltYear:           utils.NullUint32(record.BuiltYear),
		BuildingDescription: utils.NullString(record.BuildingDescription),
		KnownIssues:         utils.NullString(record.KnownIssues),
		AssessmentGoal:      utils.NullString(record.AssessmentGoal),
		ThumbnailUrl:        "",
		Progress:            record.Progress,
		CreatedAt:           record.CreatedAt.Format(time.DateTime),
		UpdatedAt:           record.UpdatedAt.Format(time.DateTime),
	}
}

// ProjectRecordToDTOWithThumbnail 将 MySQL 数据模型转换为 RPC 数据模型，并添加项目缩略图下载预签名 URL
func ProjectRecordToDTOWithThumbnail(ctx context.Context, minioRepo *minio.Repository, redisRepo *redis.Repository, record *ProjectRecord, ttl time.Duration) (*bizpb.Project, error) {
	item := ProjectRecordToDTO(record)
	if item == nil {
		return nil, nil
	}
	thumbnailURL, err := minio.PresignProjectThumbnailURL(ctx, minioRepo, redisRepo, record.Id, ttl)
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
			Progress:         record.Progress,
			CreatedAt:        record.CreatedAt.Format(time.DateTime),
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
		thumbnailURL, err := minio.PresignProjectThumbnailURL(ctx, minioRepo, redisRepo, record.Id, ttl)
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
		Status:      record.Status,
		UploadedAt:  utils.NullTimeString(record.UploadedAt),
		CreatedAt:   record.CreatedAt.Format(time.DateTime),
	}
	if record.Status != bizpb.ProjectImageStatus_Uploaded {
		return item, nil
	}
	var err error
	item.ThumbnailUrl, err = minio.PresignProjectImageThumbnailURL(ctx, minioRepo, redisRepo, record.ProjectId, record.Uuid, ttl)
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

// ProjectDetectionTaskRecord 项目图像检测主任务记录
type ProjectDetectionTaskRecord struct {
	Id                     uint64
	Uuid                   string
	UserId                 uint64
	ProjectId              uint64
	ImageId                uint64
	ImageUuid              string
	Status                 bizpb.ProjectDetectionTaskStatus_Value
	CorrosionShouldExecute bool
	CorrosionTaskId        sql.NullInt64
	CrackShouldExecute     bool
	CrackTaskId            sql.NullInt64
	StainShouldExecute     bool
	StainTaskId            sql.NullInt64
	FlatnessShouldExecute  bool
	FlatnessTaskId         sql.NullInt64
	SpallingShouldExecute  bool
	SpallingTaskId         sql.NullInt64
	SummaryShouldExecute   bool
	SummaryTaskId          sql.NullInt64
	StartedAt              sql.NullTime
	FinishedAt             sql.NullTime
	CreatedAt              time.Time
	UpdatedAt              time.Time
}

// ProjectDetectionSubTaskRecord 项目图像检测子任务记录
type ProjectDetectionSubTaskRecord struct {
	Id         uint64
	Uuid       string
	MainTaskId uint64
	UserId     uint64
	ProjectId  uint64
	ImageId    uint64
	Status     bizpb.ProjectDetectionSubTaskStatus_Value
	StartedAt  sql.NullTime
	FinishedAt sql.NullTime
	CreatedAt  time.Time
	UpdatedAt  time.Time
}

// ProjectDetectionCorrosionTaskRecord 项目图像金属锈蚀检测子任务记录
type ProjectDetectionCorrosionTaskRecord struct {
	ProjectDetectionSubTaskRecord
	HasCorrosion      sql.NullBool
	CorrosionCount    sql.NullInt64
	MaxConfidence     sql.NullFloat64
	AverageConfidence sql.NullFloat64
	CorrosionPixels   sql.NullInt64
	CorrosionRatio    sql.NullFloat64
	Regions           sql.NullString
	ArtifactSha256Map sql.NullString
	RuntimeSeconds    sql.NullFloat64
}

// ProjectDetectionCrackTaskRecord 项目图像石材裂缝检测子任务记录
type ProjectDetectionCrackTaskRecord struct {
	ProjectDetectionSubTaskRecord
	HasCrack          sql.NullBool
	CrackCount        sql.NullInt64
	CrackPixels       sql.NullInt64
	CrackRatio        sql.NullFloat64
	Regions           sql.NullString
	ArtifactSha256Map sql.NullString
	RuntimeSeconds    sql.NullFloat64
}

// ProjectDetectionStainTaskRecord 项目图像石材污渍检测子任务记录
type ProjectDetectionStainTaskRecord struct {
	ProjectDetectionSubTaskRecord
	HasStain          sql.NullBool
	StainCount        sql.NullInt64
	AverageStainRatio sql.NullFloat64
	MaxStainRatio     sql.NullFloat64
	Regions           sql.NullString
	ArtifactSha256Map sql.NullString
	RuntimeSeconds    sql.NullFloat64
}

// ProjectDetectionFlatnessTaskRecord 项目图像玻璃平整度检测子任务记录
type ProjectDetectionFlatnessTaskRecord struct {
	ProjectDetectionSubTaskRecord
	Result            sql.NullString
	UnevenCount       sql.NullInt64
	Regions           sql.NullString
	ArtifactSha256Map sql.NullString
	RuntimeSeconds    sql.NullFloat64
}

// ProjectDetectionSpallingTaskRecord 项目图像玻璃爆裂检测子任务记录
type ProjectDetectionSpallingTaskRecord struct {
	ProjectDetectionSubTaskRecord
	HasSpalling       sql.NullBool
	Confidence        sql.NullFloat64
	ArtifactSha256Map sql.NullString
	RuntimeSeconds    sql.NullFloat64
}

// ProjectDetectionSummaryTaskRecord 项目图像检测总结任务记录
type ProjectDetectionSummaryTaskRecord struct {
	ProjectDetectionSubTaskRecord
	ResultJson sql.NullString
}
