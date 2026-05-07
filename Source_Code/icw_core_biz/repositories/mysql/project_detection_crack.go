package mysql

import (
	"context"
	"database/sql"
	"encoding/json"

	"github.com/google/uuid"

	"icw_common/enum"
	"icw_common/gen/core/biz"
)

// createProjectDetectionCrackTaskTx 创建石材裂缝检测子任务
func createProjectDetectionCrackTaskTx(ctx context.Context, tx *sql.Tx, task *ProjectDetectionTaskRecord) (*ProjectDetectionSubTaskRecord, error) {
	taskUuid := uuid.NewString()
	result, err := tx.ExecContext(ctx, `
		INSERT INTO project_detection_crack_tasks(uuid, main_task_id, user_id, project_id, image_id, status)
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
		SET crack_should_execute = 1, crack_task_id = ?, status = ?
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

// findProjectDetectionCrackTaskByUuidTx 按子任务 UUID 查询石材裂缝检测子任务
func findProjectDetectionCrackTaskByUuidTx(ctx context.Context, tx *sql.Tx, taskUuid string) (*ProjectDetectionSubTaskRecord, error) {
	record, err := scanProjectDetectionSubTask(tx.QueryRowContext(ctx, `
		SELECT id, uuid, main_task_id, user_id, project_id, image_id, status, started_at, finished_at, created_at, updated_at
		FROM project_detection_crack_tasks
		WHERE uuid = ?
		FOR UPDATE
	`, taskUuid))
	if err == sql.ErrNoRows {
		return nil, nil
	}
	return record, err
}

// findProjectDetectionCrackTaskStatusByIdTx 按子任务 ID 查询石材裂缝检测状态
func findProjectDetectionCrackTaskStatusByIdTx(ctx context.Context, tx *sql.Tx, taskId uint64) (bizpb.ProjectDetectionSubTaskStatus_Value, error) {
	var status string
	if err := tx.QueryRowContext(ctx, `
		SELECT status
		FROM project_detection_crack_tasks
		WHERE id = ?
		LIMIT 1
	`, taskId).Scan(&status); err != nil {
		return bizpb.ProjectDetectionSubTaskStatus_Unknown, err
	}
	return enum.ParseProjectDetectionSubTaskStatus(status), nil
}

// updateProjectDetectionCrackTaskResultTx 写入石材裂缝检测报告并更新状态
func updateProjectDetectionCrackTaskResultTx(ctx context.Context, tx *sql.Tx, taskUuid string, status bizpb.ProjectDetectionSubTaskStatus_Value, resultJSON string) error {
	statusText := enum.ProjectDetectionSubTaskStatusString(status)
	if status != bizpb.ProjectDetectionSubTaskStatus_Succeeded {
		result, err := tx.ExecContext(ctx, `
			UPDATE project_detection_crack_tasks
			SET status = ?
			WHERE uuid = ?
		`, statusText, taskUuid)
		return checkRowsAffected(result, err)
	}

	resultJSON, err := normalizeProjectDetectionReportJSON(resultJSON)
	if err != nil {
		return err
	}
	report := projectDetectionCrackReport{}
	if err := json.Unmarshal([]byte(resultJSON), &report); err != nil {
		return err
	}
	result, err := tx.ExecContext(ctx, `
		UPDATE project_detection_crack_tasks
		SET status = ?,
			has_crack = ?,
			crack_count = ?,
			crack_pixels = ?,
			crack_ratio = ?,
			regions = ?,
			runtime_seconds = ?,
			finished_at = NOW(3),
			started_at = TIMESTAMPADD(MICROSECOND, -CAST(? * 1000000 AS SIGNED), NOW(3))
		WHERE uuid = ?
	`, statusText, report.HasCrack, report.CrackCount, report.CrackPixels, report.CrackRatio, jsonOrEmptyArray(report.Regions), report.RuntimeSeconds, report.RuntimeSeconds, taskUuid)
	return checkRowsAffected(result, err)
}
