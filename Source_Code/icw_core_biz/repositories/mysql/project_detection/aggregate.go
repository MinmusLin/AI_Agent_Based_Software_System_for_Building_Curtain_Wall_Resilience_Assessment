package project_detection

import (
	"context"
	"database/sql"

	"icw_common/enum"
	"icw_common/gen/activity"
	"icw_common/gen/core/biz"
)

// ProjectDetectionSubTaskAggregateStatusTx 统计项目图像检测子任务聚合状态
func ProjectDetectionSubTaskAggregateStatusTx(ctx context.Context, tx *sql.Tx, mainTaskId uint64) (bool, bool, error) {
	task, err := FindProjectDetectionTaskByIdForUpdateTx(ctx, tx, mainTaskId)
	if err != nil {
		return false, false, err
	}

	checks := []struct {
		shouldExecute bool
		taskId        sql.NullInt64
		taskCode      string
	}{{
		shouldExecute: task.CorrosionShouldExecute,
		taskId:        task.CorrosionTaskId,
		taskCode:      enum.DetectionTaskCodeString(activitypb.DetectionTaskCode_Corrosion),
	}, {
		shouldExecute: task.CrackShouldExecute,
		taskId:        task.CrackTaskId,
		taskCode:      enum.DetectionTaskCodeString(activitypb.DetectionTaskCode_Crack),
	}, {
		shouldExecute: task.StainShouldExecute,
		taskId:        task.StainTaskId,
		taskCode:      enum.DetectionTaskCodeString(activitypb.DetectionTaskCode_Stain),
	}, {
		shouldExecute: task.FlatnessShouldExecute,
		taskId:        task.FlatnessTaskId,
		taskCode:      enum.DetectionTaskCodeString(activitypb.DetectionTaskCode_Flatness),
	}, {
		shouldExecute: task.SpallingShouldExecute,
		taskId:        task.SpallingTaskId,
		taskCode:      enum.DetectionTaskCodeString(activitypb.DetectionTaskCode_Spalling),
	}}

	var (
		anyFailed      bool
		allSucceeded   bool
		executedCount  int
		succeededCount int
	)

	for _, check := range checks {
		if !check.shouldExecute {
			continue
		}
		executedCount++
		if !check.taskId.Valid {
			continue
		}
		status, err := findProjectDetectionSubTaskStatusByIdTx(ctx, tx, check.taskCode, uint64(check.taskId.Int64))
		if err != nil {
			return false, false, err
		}
		switch status {
		case bizpb.ProjectDetectionSubTaskStatus_Succeeded:
			succeededCount++
		default:
			anyFailed = true

		}
	}

	allSucceeded = executedCount > 0 && succeededCount == executedCount

	return anyFailed, allSucceeded, nil
}
