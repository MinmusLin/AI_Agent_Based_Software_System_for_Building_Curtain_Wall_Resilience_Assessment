package mysql

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"

	"github.com/google/uuid"

	"icw_common/enum"
	"icw_common/gen/core/biz"
)

// createProjectDetectionStainTaskTx 创建项目图像石材污渍检测子任务记录
func createProjectDetectionStainTaskTx(ctx context.Context, tx *sql.Tx, task *ProjectDetectionTaskRecord) (*ProjectDetectionSubTaskRecord, error) {
	taskUuid := uuid.NewString()

	result, err := tx.ExecContext(ctx, `
		INSERT INTO project_detection_stain_tasks(uuid, main_task_id, user_id, project_id, image_id, status)
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
		SET stain_should_execute = 1, stain_task_id = ?, status = ?
		WHERE id = ?
	`, subTaskId, enum.ProjectDetectionTaskStatusString(bizpb.ProjectDetectionTaskStatus_Detecting), task.Id)
	if err := checkRowsAffected(result, err); err != nil {
		return nil, err
	}

	return &ProjectDetectionSubTaskRecord{
		Id:         uint64(subTaskId),
		Uuid:       taskUuid,
		MainTaskId: task.Id,
		UserId:     task.UserId,
		ProjectId:  task.ProjectId,
		ImageId:    task.ImageId,
		Status:     bizpb.ProjectDetectionSubTaskStatus_Pending,
	}, nil
}

// findProjectDetectionStainTaskByUuidTx 按子任务 UUID 查询项目图像石材污渍检测子任务记录
func findProjectDetectionStainTaskByUuidTx(ctx context.Context, tx *sql.Tx, taskUuid string) (*ProjectDetectionSubTaskRecord, error) {
	record, err := scanProjectDetectionSubTask(tx.QueryRowContext(ctx, `
		SELECT id, uuid, main_task_id, user_id, project_id, image_id, status, started_at, finished_at, created_at, updated_at
		FROM project_detection_stain_tasks
		WHERE uuid = ?
		FOR UPDATE
	`, taskUuid))
	if errors.Is(err, sql.ErrNoRows) {
		return nil, nil
	}
	return record, err
}

// findProjectDetectionStainTaskStatusByIdTx 按子任务 ID 查询项目图像石材污渍检测子任务状态
func findProjectDetectionStainTaskStatusByIdTx(ctx context.Context, tx *sql.Tx, taskId uint64) (bizpb.ProjectDetectionSubTaskStatus_Value, error) {
	var status string
	if err := tx.QueryRowContext(ctx, `
		SELECT status
		FROM project_detection_stain_tasks
		WHERE id = ?
		LIMIT 1
	`, taskId).Scan(&status); err != nil {
		return bizpb.ProjectDetectionSubTaskStatus_Unknown, err
	}
	return enum.ParseProjectDetectionSubTaskStatus(status), nil
}

// updateProjectDetectionStainTaskResultTx 更新项目图像石材污渍检测子任务报告与状态
func updateProjectDetectionStainTaskResultTx(ctx context.Context, tx *sql.Tx, taskUuid string, status bizpb.ProjectDetectionSubTaskStatus_Value, resultJSON string) error {
	statusText := enum.ProjectDetectionSubTaskStatusString(status)

	// 检测失败时，只更新任务状态
	if status != bizpb.ProjectDetectionSubTaskStatus_Succeeded {
		result, err := tx.ExecContext(ctx, `
			UPDATE project_detection_stain_tasks
			SET status = ?
			WHERE uuid = ?
		`, statusText, taskUuid)

		return checkRowsAffected(result, err)
	}

	// 检测成功时，更新任务状态、开始时间、完成时间和检测报告
	report := struct {
		HasStain          bool            `json:"has_stain"`
		StainCount        uint64          `json:"stain_count"`
		AverageStainRatio float64         `json:"average_stain_ratio"`
		MaxStainRatio     float64         `json:"max_stain_ratio"`
		Regions           json.RawMessage `json:"regions"`
		RuntimeSeconds    float64         `json:"runtime_seconds"`
	}{}
	if err := json.Unmarshal([]byte(resultJSON), &report); err != nil {
		return err
	}

	result, err := tx.ExecContext(ctx, `
		UPDATE project_detection_stain_tasks
		SET status = ?,
			has_stain = ?,
			stain_count = ?,
			average_stain_ratio = ?,
			max_stain_ratio = ?,
			regions = ?,
			runtime_seconds = ?,
			started_at = TIMESTAMPADD(MICROSECOND, -CAST(? * 1000000 AS SIGNED), NOW(3)),
			finished_at = NOW(3)
		WHERE uuid = ?
	`, statusText, report.HasStain, report.StainCount, report.AverageStainRatio, report.MaxStainRatio, jsonOrEmptyArray(report.Regions), report.RuntimeSeconds, report.RuntimeSeconds, taskUuid)

	return checkRowsAffected(result, err)
}
