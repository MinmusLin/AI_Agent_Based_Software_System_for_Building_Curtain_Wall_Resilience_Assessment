package model

import (
	"database/sql"
	"time"

	"icw_common/gen/core/biz"
)

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
	RuntimeSeconds    sql.NullFloat64
}

// ProjectDetectionCrackTaskRecord 项目图像石材裂缝检测子任务记录
type ProjectDetectionCrackTaskRecord struct {
	ProjectDetectionSubTaskRecord
	HasCrack       sql.NullBool
	CrackCount     sql.NullInt64
	CrackPixels    sql.NullInt64
	CrackRatio     sql.NullFloat64
	Regions        sql.NullString
	RuntimeSeconds sql.NullFloat64
}

// ProjectDetectionStainTaskRecord 项目图像石材污渍检测子任务记录
type ProjectDetectionStainTaskRecord struct {
	ProjectDetectionSubTaskRecord
	HasStain          sql.NullBool
	StainCount        sql.NullInt64
	AverageStainRatio sql.NullFloat64
	MaxStainRatio     sql.NullFloat64
	Regions           sql.NullString
	RuntimeSeconds    sql.NullFloat64
}

// ProjectDetectionFlatnessTaskRecord 项目图像玻璃平整度检测子任务记录
type ProjectDetectionFlatnessTaskRecord struct {
	ProjectDetectionSubTaskRecord
	Result         sql.NullString
	UnevenCount    sql.NullInt64
	Regions        sql.NullString
	RuntimeSeconds sql.NullFloat64
}

// ProjectDetectionSpallingTaskRecord 项目图像玻璃爆裂检测子任务记录
type ProjectDetectionSpallingTaskRecord struct {
	ProjectDetectionSubTaskRecord
	HasSpalling    sql.NullBool
	Confidence     sql.NullFloat64
	RuntimeSeconds sql.NullFloat64
}

// ProjectDetectionSummaryTaskRecord 项目图像检测总结任务记录
type ProjectDetectionSummaryTaskRecord struct {
	ProjectDetectionSubTaskRecord
}
