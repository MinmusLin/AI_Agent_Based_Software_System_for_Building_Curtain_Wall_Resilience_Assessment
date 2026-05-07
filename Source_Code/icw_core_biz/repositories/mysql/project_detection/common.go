package project_detection

import (
	"context"
	"database/sql"
	"errors"
	"icw_common/enum"
	"icw_common/gen/activity"
	"icw_common/gen/core/biz"
	"icw_core_biz/repositories/mysql"
)

var (
	// ErrUnsupportedDetectionTaskCode 不支持的原子检测能力代码
	ErrUnsupportedDetectionTaskCode = errors.New("unsupported detection task code")
)

func createProjectDetectionSubTaskTx(ctx context.Context, tx *sql.Tx, task *mysql.ProjectDetectionTaskRecord, taskCode string) (*mysql.ProjectDetectionSubTaskRecord, error) {
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
		return nil, ErrUnsupportedDetectionTaskCode
	}
}

func updateProjectDetectionSubTaskResultTx(ctx context.Context, tx *sql.Tx, taskCode, taskUuid string, status bizpb.ProjectDetectionSubTaskStatus_Value, resultJSON string) error {
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
		return ErrUnsupportedDetectionTaskCode
	}
}

func findProjectDetectionSubTaskByUuidTx(ctx context.Context, tx *sql.Tx, taskCode, taskUuid string) (*mysql.ProjectDetectionSubTaskRecord, error) {
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
		return nil, ErrUnsupportedDetectionTaskCode
	}
}
