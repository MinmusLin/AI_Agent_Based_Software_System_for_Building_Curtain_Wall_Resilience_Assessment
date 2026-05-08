package project_detection

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/google/uuid"

	"icw_common/enum"
	"icw_common/gen/activity"
	"icw_common/gen/core/biz"

	"icw_core_biz/repositories/mysql/model"
	"icw_core_biz/repositories/mysql/utils"
)

// FindProjectDetectionSubTaskByIdFromTableTx 按子任务表名和子任务 ID 查询项目图像检测子任务记录
func FindProjectDetectionSubTaskByIdFromTableTx(ctx context.Context, tx *sql.Tx, table string, taskId uint64) (*model.ProjectDetectionSubTaskRecord, error) {
	return model.ScanProjectDetectionSubTask(tx.QueryRowContext(ctx, fmt.Sprintf(`
		SELECT id, uuid, main_task_id, user_id, project_id, image_id, status, started_at, finished_at, created_at, updated_at
		FROM %s
		WHERE id = ?
		LIMIT 1
	`, table), taskId))
}

// findProjectDetectionSubTaskStatusByIdTx 按子任务代码和子任务 ID 查询项目图像检测子任务状态
func findProjectDetectionSubTaskStatusByIdTx(ctx context.Context, tx *sql.Tx, taskCode string, taskId uint64) (bizpb.ProjectDetectionSubTaskStatus_Value, error) {
	schema, err := projectDetectionSubTaskSchemaByCode(taskCode)
	if err != nil {
		return bizpb.ProjectDetectionSubTaskStatus_Unknown, err
	}

	var status string
	if err := tx.QueryRowContext(ctx, fmt.Sprintf(`
		SELECT status
		FROM %s
		WHERE id = ?
		LIMIT 1
	`, schema.table), taskId).Scan(&status); err != nil {
		return bizpb.ProjectDetectionSubTaskStatus_Unknown, err
	}
	return enum.ParseProjectDetectionSubTaskStatus(status), nil
}

// updateProjectDetectionSubTaskStatusByTableTx 按子任务表名和子任务 UUID 更新项目图像检测子任务状态
func updateProjectDetectionSubTaskStatusByTableTx(ctx context.Context, tx *sql.Tx, table, taskUuid string, status bizpb.ProjectDetectionSubTaskStatus_Value, startTime, finishTime bool, invalidStatusErr error) error {
	statusText := enum.ProjectDetectionSubTaskStatusString(status)
	if statusText == "" {
		return invalidStatusErr
	}
	result, err := tx.ExecContext(ctx, fmt.Sprintf(`
		UPDATE %s
		SET status = ?,
			started_at = CASE WHEN ? AND started_at IS NULL THEN NOW(3) ELSE started_at END,
			finished_at = CASE WHEN ? AND finished_at IS NULL THEN NOW(3) ELSE finished_at END
		WHERE uuid = ?
	`, table), statusText, startTime, finishTime, taskUuid)
	return utils.CheckRowsAffected(result, err)
}

// FindProjectDetectionSubTaskByIdTx 按子任务代码和子任务 ID 查询项目图像检测子任务记录
func FindProjectDetectionSubTaskByIdTx(ctx context.Context, tx *sql.Tx, taskCode string, taskId uint64) (*model.ProjectDetectionSubTaskRecord, error) {
	schema, err := projectDetectionSubTaskSchemaByCode(taskCode)
	if err != nil {
		return nil, err
	}
	return FindProjectDetectionSubTaskByIdFromTableTx(ctx, tx, schema.table, taskId)
}

// FindProjectDetectionSubTaskByUuidTx 按子任务代码和子任务 UUID 查询项目图像检测子任务记录
func FindProjectDetectionSubTaskByUuidTx(ctx context.Context, tx *sql.Tx, taskCode, taskUuid string) (*model.ProjectDetectionSubTaskRecord, error) {
	schema, err := projectDetectionSubTaskSchemaByCode(taskCode)
	if err != nil {
		return nil, err
	}
	record, err := model.ScanProjectDetectionSubTask(tx.QueryRowContext(ctx, fmt.Sprintf(`
		SELECT id, uuid, main_task_id, user_id, project_id, image_id, status, started_at, finished_at, created_at, updated_at
		FROM %s
		WHERE uuid = ?
		FOR UPDATE
	`, schema.table), taskUuid))
	if errors.Is(err, sql.ErrNoRows) {
		return nil, nil
	}
	return record, err
}

// CreateProjectDetectionSubTaskTx 创建项目图像检测子任务记录
func CreateProjectDetectionSubTaskTx(ctx context.Context, tx *sql.Tx, task *model.ProjectDetectionTaskRecord, taskCode string) (*model.ProjectDetectionSubTaskRecord, error) {
	schema, err := projectDetectionSubTaskSchemaByCode(taskCode)
	if err != nil {
		return nil, err
	}

	taskUuid := uuid.NewString()
	result, err := tx.ExecContext(ctx, fmt.Sprintf(`
		INSERT INTO %s(uuid, main_task_id, user_id, project_id, image_id, status)
		VALUES (?, ?, ?, ?, ?, ?)
	`, schema.table), taskUuid, task.Id, task.UserId, task.ProjectId, task.ImageId, enum.ProjectDetectionSubTaskStatusString(bizpb.ProjectDetectionSubTaskStatus_Pending))
	if err != nil {
		return nil, err
	}

	subTaskId, err := result.LastInsertId()
	if err != nil {
		return nil, err
	}

	result, err = tx.ExecContext(ctx, fmt.Sprintf(`
		UPDATE project_detection_tasks
		SET %s = 1, %s = ?, status = ?
		WHERE id = ?
	`, schema.shouldExecuteColumn, schema.taskIdColumn), subTaskId, enum.ProjectDetectionTaskStatusString(bizpb.ProjectDetectionTaskStatus_Detecting), task.Id)
	if err := utils.CheckRowsAffected(result, err); err != nil {
		return nil, err
	}

	return &model.ProjectDetectionSubTaskRecord{
		Id:         uint64(subTaskId),
		Uuid:       taskUuid,
		MainTaskId: task.Id,
		UserId:     task.UserId,
		ProjectId:  task.ProjectId,
		ImageId:    task.ImageId,
		Status:     bizpb.ProjectDetectionSubTaskStatus_Pending,
	}, nil
}

// UpdateProjectDetectionSubTaskStatusTx 更新项目图像检测子任务状态
func UpdateProjectDetectionSubTaskStatusTx(ctx context.Context, tx *sql.Tx, taskCode, taskUuid string, status bizpb.ProjectDetectionSubTaskStatus_Value, startTime, finishTime bool) error {
	schema, err := projectDetectionSubTaskSchemaByCode(taskCode)
	if err != nil {
		return err
	}
	return updateProjectDetectionSubTaskStatusByTableTx(ctx, tx, schema.table, taskUuid, status, startTime, finishTime, model.ErrProjectDetectionSubTaskStatusInvalid)
}

// UpdateProjectDetectionSubTaskResultTx 更新项目图像检测子任务报告与状态
func UpdateProjectDetectionSubTaskResultTx(ctx context.Context, tx *sql.Tx, taskCode, taskUuid string, status bizpb.ProjectDetectionSubTaskStatus_Value, resultJSON, artifactSha256Map string) error {
	switch enum.ParseDetectionTaskCode(taskCode) {
	case activitypb.DetectionTaskCode_Corrosion:
		return updateProjectDetectionCorrosionTaskResultTx(ctx, tx, taskUuid, status, resultJSON, artifactSha256Map)
	case activitypb.DetectionTaskCode_Crack:
		return updateProjectDetectionCrackTaskResultTx(ctx, tx, taskUuid, status, resultJSON, artifactSha256Map)
	case activitypb.DetectionTaskCode_Stain:
		return updateProjectDetectionStainTaskResultTx(ctx, tx, taskUuid, status, resultJSON, artifactSha256Map)
	case activitypb.DetectionTaskCode_Flatness:
		return updateProjectDetectionFlatnessTaskResultTx(ctx, tx, taskUuid, status, resultJSON, artifactSha256Map)
	case activitypb.DetectionTaskCode_Spalling:
		return updateProjectDetectionSpallingTaskResultTx(ctx, tx, taskUuid, status, resultJSON, artifactSha256Map)
	default:
		return model.ErrUnsupportedDetectionTaskCode
	}
}
