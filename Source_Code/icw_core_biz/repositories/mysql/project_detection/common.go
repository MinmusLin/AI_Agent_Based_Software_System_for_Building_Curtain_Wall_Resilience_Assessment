package project_detection

import (
	"context"
	"database/sql"

	"icw_common/enum"
	"icw_common/gen/activity"
	"icw_common/gen/core/biz"

	"icw_core_biz/repositories/mysql/model"
)

// FindProjectDetectionSubTaskByUuidTx 按子任务代码和子任务 UUID 查询项目图像检测子任务记录
func FindProjectDetectionSubTaskByUuidTx(ctx context.Context, tx *sql.Tx, taskCode, taskUuid string) (*model.ProjectDetectionSubTaskRecord, error) {
	switch enum.ParseDetectionTaskCode(taskCode) {
	case activitypb.DetectionTaskCode_Corrosion:
		return findProjectDetectionCorrosionTaskByUuidTx(ctx, tx, taskUuid)
	case activitypb.DetectionTaskCode_Crack:
		return findProjectDetectionCrackTaskByUuidTx(ctx, tx, taskUuid)
	case activitypb.DetectionTaskCode_Stain:
		return findProjectDetectionStainTaskByUuidTx(ctx, tx, taskUuid)
	case activitypb.DetectionTaskCode_Flatness:
		return findProjectDetectionFlatnessTaskByUuidTx(ctx, tx, taskUuid)
	case activitypb.DetectionTaskCode_Spalling:
		return findProjectDetectionSpallingTaskByUuidTx(ctx, tx, taskUuid)
	default:
		return nil, model.ErrUnsupportedDetectionTaskCode
	}
}

// CreateProjectDetectionSubTaskTx 创建项目图像检测子任务记录
func CreateProjectDetectionSubTaskTx(ctx context.Context, tx *sql.Tx, task *model.ProjectDetectionTaskRecord, taskCode string) (*model.ProjectDetectionSubTaskRecord, error) {
	switch enum.ParseDetectionTaskCode(taskCode) {
	case activitypb.DetectionTaskCode_Corrosion:
		return createProjectDetectionCorrosionTaskTx(ctx, tx, task)
	case activitypb.DetectionTaskCode_Crack:
		return createProjectDetectionCrackTaskTx(ctx, tx, task)
	case activitypb.DetectionTaskCode_Stain:
		return createProjectDetectionStainTaskTx(ctx, tx, task)
	case activitypb.DetectionTaskCode_Flatness:
		return createProjectDetectionFlatnessTaskTx(ctx, tx, task)
	case activitypb.DetectionTaskCode_Spalling:
		return createProjectDetectionSpallingTaskTx(ctx, tx, task)
	default:
		return nil, model.ErrUnsupportedDetectionTaskCode
	}
}

// UpdateProjectDetectionSubTaskResultTx 更新项目图像检测子任务报告与状态
func UpdateProjectDetectionSubTaskResultTx(ctx context.Context, tx *sql.Tx, taskCode, taskUuid string, status bizpb.ProjectDetectionSubTaskStatus_Value, resultJSON string) error {
	switch enum.ParseDetectionTaskCode(taskCode) {
	case activitypb.DetectionTaskCode_Corrosion:
		return updateProjectDetectionCorrosionTaskResultTx(ctx, tx, taskUuid, status, resultJSON)
	case activitypb.DetectionTaskCode_Crack:
		return updateProjectDetectionCrackTaskResultTx(ctx, tx, taskUuid, status, resultJSON)
	case activitypb.DetectionTaskCode_Stain:
		return updateProjectDetectionStainTaskResultTx(ctx, tx, taskUuid, status, resultJSON)
	case activitypb.DetectionTaskCode_Flatness:
		return updateProjectDetectionFlatnessTaskResultTx(ctx, tx, taskUuid, status, resultJSON)
	case activitypb.DetectionTaskCode_Spalling:
		return updateProjectDetectionSpallingTaskResultTx(ctx, tx, taskUuid, status, resultJSON)
	default:
		return model.ErrUnsupportedDetectionTaskCode
	}
}
