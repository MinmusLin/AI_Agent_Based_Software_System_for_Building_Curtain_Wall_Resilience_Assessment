package project_detection

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"

	"github.com/google/uuid"

	"icw_common/enum"
	"icw_common/gen/core/biz"

	"icw_core_biz/repositories/mysql/model"
	"icw_core_biz/repositories/mysql/utils"
)

// FindProjectDetectionCorrosionTaskStatusByIdTx 按子任务 ID 查询项目图像金属锈蚀检测子任务状态
func FindProjectDetectionCorrosionTaskStatusByIdTx(ctx context.Context, tx *sql.Tx, taskId uint64) (bizpb.ProjectDetectionSubTaskStatus_Value, error) {
	var status string
	if err := tx.QueryRowContext(ctx, `
		SELECT status
		FROM project_detection_corrosion_tasks
		WHERE id = ?
		LIMIT 1
	`, taskId).Scan(&status); err != nil {
		return bizpb.ProjectDetectionSubTaskStatus_Unknown, err
	}
	return enum.ParseProjectDetectionSubTaskStatus(status), nil
}

// findProjectDetectionCorrosionTaskByUuidTx 按子任务 UUID 查询项目图像金属锈蚀检测子任务记录
func findProjectDetectionCorrosionTaskByUuidTx(ctx context.Context, tx *sql.Tx, taskUuid string) (*model.ProjectDetectionSubTaskRecord, error) {
	record, err := utils.ScanProjectDetectionSubTask(tx.QueryRowContext(ctx, `
		SELECT id, uuid, main_task_id, user_id, project_id, image_id, status, started_at, finished_at, created_at, updated_at
		FROM project_detection_corrosion_tasks
		WHERE uuid = ?
		FOR UPDATE
	`, taskUuid))
	if errors.Is(err, sql.ErrNoRows) {
		return nil, nil
	}
	return record, err
}

// createProjectDetectionCorrosionTaskTx 创建项目图像金属锈蚀检测子任务记录
func createProjectDetectionCorrosionTaskTx(ctx context.Context, tx *sql.Tx, task *model.ProjectDetectionTaskRecord) (*model.ProjectDetectionSubTaskRecord, error) {
	taskUuid := uuid.NewString()

	result, err := tx.ExecContext(ctx, `
		INSERT INTO project_detection_corrosion_tasks(uuid, main_task_id, user_id, project_id, image_id, status)
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
		SET corrosion_should_execute = 1, corrosion_task_id = ?, status = ?
		WHERE id = ?
	`, subTaskId, enum.ProjectDetectionTaskStatusString(bizpb.ProjectDetectionTaskStatus_Detecting), task.Id)
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

// updateProjectDetectionCorrosionTaskResultTx 更新项目图像金属锈蚀检测子任务报告与状态
func updateProjectDetectionCorrosionTaskResultTx(ctx context.Context, tx *sql.Tx, taskUuid string, status bizpb.ProjectDetectionSubTaskStatus_Value, resultJSON string) error {
	statusText := enum.ProjectDetectionSubTaskStatusString(status)

	// 检测失败时，只更新任务状态
	if status != bizpb.ProjectDetectionSubTaskStatus_Succeeded {
		result, err := tx.ExecContext(ctx, `
			UPDATE project_detection_corrosion_tasks
			SET status = ?
			WHERE uuid = ?
		`, statusText, taskUuid)

		return utils.CheckRowsAffected(result, err)
	}

	// 检测成功时，更新任务状态、开始时间、完成时间和检测报告
	report := struct {
		HasCorrosion      bool            `json:"has_corrosion"`
		CorrosionCount    uint64          `json:"corrosion_count"`
		MaxConfidence     float64         `json:"max_confidence"`
		AverageConfidence float64         `json:"average_confidence"`
		CorrosionPixels   uint64          `json:"corrosion_pixels"`
		CorrosionRatio    float64         `json:"corrosion_ratio"`
		Regions           json.RawMessage `json:"regions"`
		RuntimeSeconds    float64         `json:"runtime_seconds"`
	}{}
	if err := json.Unmarshal([]byte(resultJSON), &report); err != nil {
		return err
	}

	result, err := tx.ExecContext(ctx, `
		UPDATE project_detection_corrosion_tasks
		SET status = ?,
			has_corrosion = ?,
			corrosion_count = ?,
			max_confidence = ?,
			average_confidence = ?,
			corrosion_pixels = ?,
			corrosion_ratio = ?,
			regions = ?,
			runtime_seconds = ?,
			started_at = TIMESTAMPADD(MICROSECOND, -CAST(? * 1000000 AS SIGNED), NOW(3)),
			finished_at = NOW(3)
		WHERE uuid = ?
	`, statusText, report.HasCorrosion, report.CorrosionCount, report.MaxConfidence, report.AverageConfidence, report.CorrosionPixels, report.CorrosionRatio, utils.JsonOrEmptyArray(report.Regions), report.RuntimeSeconds, report.RuntimeSeconds, taskUuid)

	return utils.CheckRowsAffected(result, err)
}
