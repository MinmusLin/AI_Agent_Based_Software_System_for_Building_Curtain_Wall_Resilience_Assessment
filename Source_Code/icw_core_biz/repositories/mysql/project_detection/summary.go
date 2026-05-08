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
	record := &model.ProjectDetectionSummaryTaskRecord{}
	var status string
	if err := tx.QueryRowContext(ctx, `
		SELECT
			id,
			uuid,
			main_task_id,
			user_id,
			project_id,
			image_id,
			status,
			CAST(result_json AS CHAR),
			started_at,
			finished_at,
			created_at,
			updated_at
		FROM project_detection_summary_tasks
		WHERE uuid = ?
		FOR UPDATE
	`, taskUuid).Scan(
		&record.Id,
		&record.Uuid,
		&record.MainTaskId,
		&record.UserId,
		&record.ProjectId,
		&record.ImageId,
		&status,
		&record.ResultJson,
		&record.StartedAt,
		&record.FinishedAt,
		&record.CreatedAt,
		&record.UpdatedAt,
	); err != nil {
		return nil, err
	}
	record.Status = enum.ParseProjectDetectionSubTaskStatus(status)
	return record, nil
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
func UpdateProjectDetectionSummaryTaskStatusTx(ctx context.Context, tx *sql.Tx, taskUuid string, status bizpb.ProjectDetectionSubTaskStatus_Value, startTime, finishTime bool) error {
	return updateProjectDetectionSubTaskStatusByTableTx(ctx, tx, "project_detection_summary_tasks", taskUuid, status, startTime, finishTime, model.ErrProjectDetectionSummaryTaskStatusInvalid)
}

// UpdateProjectDetectionSummaryTaskResultTx 更新图像检测总结任务报告与状态
func UpdateProjectDetectionSummaryTaskResultTx(ctx context.Context, tx *sql.Tx, taskUuid string, status bizpb.ProjectDetectionSubTaskStatus_Value, resultJSON string) error {
	result, err := tx.ExecContext(ctx, `
		UPDATE project_detection_summary_tasks
		SET result_json = CASE WHEN ? = ? THEN ? ELSE result_json END
		WHERE uuid = ?
	`, enum.ProjectDetectionSubTaskStatusString(status), enum.ProjectDetectionSubTaskStatusString(bizpb.ProjectDetectionSubTaskStatus_Succeeded), utils.JsonStringOrEmptyObject(resultJSON), taskUuid)
	if err := utils.CheckRowsAffected(result, err); err != nil {
		return err
	}

	return UpdateProjectDetectionSummaryTaskStatusTx(ctx, tx, taskUuid, status, false, true)
}
