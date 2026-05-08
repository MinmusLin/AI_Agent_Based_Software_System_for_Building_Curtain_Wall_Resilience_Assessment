package project_detection

import (
	"context"
	"database/sql"

	"icw_common/enum"
	"icw_common/gen/core/biz"

	"icw_core_biz/repositories/mysql/model"
	"icw_core_biz/repositories/mysql/utils"
)

// FindProjectDetectionTaskByIdForUpdateTx 按主任务 ID 查询并锁定项目图像检测主任务记录
func FindProjectDetectionTaskByIdForUpdateTx(ctx context.Context, tx *sql.Tx, taskId uint64) (*model.ProjectDetectionTaskRecord, error) {
	return model.ScanProjectDetectionTask(tx.QueryRowContext(ctx, `
		SELECT
			id,
			uuid,
			user_id,
			project_id,
			image_id,
			image_uuid,
			status,
			corrosion_should_execute,
			corrosion_task_id,
			crack_should_execute,
			crack_task_id,
			stain_should_execute,
			stain_task_id,
			flatness_should_execute,
			flatness_task_id,
			spalling_should_execute,
			spalling_task_id,
			summary_should_execute,
			summary_task_id,
			started_at,
			finished_at,
			created_at,
			updated_at
		FROM project_detection_tasks
		WHERE id = ?
		FOR UPDATE
	`, taskId))
}

// FindProjectDetectionTaskByUuidTx 按主任务 UUID 查询并锁定项目图像检测主任务记录
func FindProjectDetectionTaskByUuidTx(ctx context.Context, tx *sql.Tx, taskUuid string) (*model.ProjectDetectionTaskRecord, error) {
	return model.ScanProjectDetectionTask(tx.QueryRowContext(ctx, `
		SELECT
			id,
			uuid,
			user_id,
			project_id,
			image_id,
			image_uuid,
			status,
			corrosion_should_execute,
			corrosion_task_id,
			crack_should_execute,
			crack_task_id,
			stain_should_execute,
			stain_task_id,
			flatness_should_execute,
			flatness_task_id,
			spalling_should_execute,
			spalling_task_id,
			summary_should_execute,
			summary_task_id,
			started_at,
			finished_at,
			created_at,
			updated_at
		FROM project_detection_tasks
		WHERE uuid = ?
		FOR UPDATE
	`, taskUuid))
}

// UpdateProjectDetectionTaskStatusTx 更新项目图像检测主任务状态
func UpdateProjectDetectionTaskStatusTx(ctx context.Context, tx *sql.Tx, taskId uint64, status bizpb.ProjectDetectionTaskStatus_Value) error {
	statusText := enum.ProjectDetectionTaskStatusString(status)
	if statusText == "" {
		return model.ErrProjectDetectionTaskStatusInvalid
	}

	result, err := tx.ExecContext(ctx, `
		UPDATE project_detection_tasks
		SET status = ?,
			started_at = CASE WHEN started_at IS NULL AND ? IN (?, ?, ?) THEN NOW(3) ELSE started_at END,
			finished_at = CASE WHEN ? IN (?, ?) THEN NOW(3) ELSE finished_at END
		WHERE id = ?
	`, statusText,
		statusText,
		enum.ProjectDetectionTaskStatusString(bizpb.ProjectDetectionTaskStatus_Classifying),
		enum.ProjectDetectionTaskStatusString(bizpb.ProjectDetectionTaskStatus_Detecting),
		enum.ProjectDetectionTaskStatusString(bizpb.ProjectDetectionTaskStatus_Summarizing),
		statusText,
		enum.ProjectDetectionTaskStatusString(bizpb.ProjectDetectionTaskStatus_Succeeded),
		enum.ProjectDetectionTaskStatusString(bizpb.ProjectDetectionTaskStatus_Failed),
		taskId)

	return utils.CheckRowsAffected(result, err)
}
