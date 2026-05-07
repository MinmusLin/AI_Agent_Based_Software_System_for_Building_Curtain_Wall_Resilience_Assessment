package project_detection

import (
	"context"
	"database/sql"

	"github.com/google/uuid"

	"icw_common/enum"
	"icw_common/gen/core/biz"

	"icw_core_biz/repositories/mysql/model"
	"icw_core_biz/repositories/mysql/utils"
)

// FindProjectDetectionSummaryTaskByUuidTx 按总结任务 UUID 查询图像检测总结任务
func FindProjectDetectionSummaryTaskByUuidTx(ctx context.Context, tx *sql.Tx, taskUuid string) (*model.ProjectDetectionSummaryTaskRecord, error) {
	subTask, err := utils.ScanProjectDetectionSubTask(tx.QueryRowContext(ctx, `
		SELECT
			id,
			uuid,
			main_task_id,
			user_id,
			project_id,
			image_id,
			status,
			started_at,
			finished_at,
			created_at,
			updated_at
		FROM project_detection_summary_tasks
		WHERE uuid = ?
		FOR UPDATE
	`, taskUuid))
	if err != nil {
		return nil, err
	}
	if subTask == nil {
		return nil, nil
	}
	return &model.ProjectDetectionSummaryTaskRecord{
		ProjectDetectionSubTaskRecord: *subTask,
	}, nil
}

// CreateProjectDetectionSummaryTaskTx 创建图像检测总结任务
func CreateProjectDetectionSummaryTaskTx(ctx context.Context, tx *sql.Tx, task *model.ProjectDetectionTaskRecord) (*model.ProjectDetectionSummaryTaskRecord, error) {
	taskUuid := uuid.NewString()
	result, err := tx.ExecContext(ctx, `
		INSERT INTO project_detection_summary_tasks(uuid, main_task_id, user_id, project_id, image_id, status)
		VALUES (?, ?, ?, ?, ?, ?)
	`, taskUuid, task.Id, task.UserId, task.ProjectId, task.ImageId, enum.ProjectDetectionSubTaskStatusString(bizpb.ProjectDetectionSubTaskStatus_Pending))
	if err != nil {
		return nil, err
	}

	subTaskId, err := result.LastInsertId()
	if err != nil {
		return nil, err
	}

	result, err = tx.ExecContext(ctx, `
		UPDATE project_detection_tasks
		SET summary_should_execute = 1, summary_task_id = ?, status = ?
		WHERE id = ?
	`, subTaskId, enum.ProjectDetectionTaskStatusString(bizpb.ProjectDetectionTaskStatus_Summarizing), task.Id)
	if err := utils.CheckRowsAffected(result, err); err != nil {
		return nil, err
	}

	return &model.ProjectDetectionSummaryTaskRecord{
		ProjectDetectionSubTaskRecord: model.ProjectDetectionSubTaskRecord{
			Id:         uint64(subTaskId),
			Uuid:       taskUuid,
			MainTaskId: task.Id,
			UserId:     task.UserId,
			ProjectId:  task.ProjectId,
			ImageId:    task.ImageId,
			Status:     bizpb.ProjectDetectionSubTaskStatus_Pending,
		},
	}, nil
}

// UpdateProjectDetectionSummaryTaskStatusTx 更新图像检测总结任务状态
func UpdateProjectDetectionSummaryTaskStatusTx(ctx context.Context, tx *sql.Tx, taskUuid string, status bizpb.ProjectDetectionSubTaskStatus_Value, runtimeSeconds float64) error {
	statusText := enum.ProjectDetectionSubTaskStatusString(status)
	if statusText == "" {
		return model.ErrProjectDetectionSummaryTaskStatusInvalid
	}
	result, err := tx.ExecContext(ctx, `
		UPDATE project_detection_summary_tasks
		SET status = ?,
			finished_at = CASE WHEN ? = ? THEN NOW(3) ELSE finished_at END,
			started_at = CASE WHEN ? = ? AND ? > 0 THEN TIMESTAMPADD(MICROSECOND, -CAST(? * 1000000 AS SIGNED), NOW(3)) ELSE started_at END
		WHERE uuid = ?
	`, statusText,
		statusText,
		enum.ProjectDetectionSubTaskStatusString(bizpb.ProjectDetectionSubTaskStatus_Succeeded),
		statusText,
		enum.ProjectDetectionSubTaskStatusString(bizpb.ProjectDetectionSubTaskStatus_Succeeded),
		runtimeSeconds,
		runtimeSeconds,
		taskUuid)
	return utils.CheckRowsAffected(result, err)
}
