package project_detection

import (
	"context"
	"database/sql"
	"fmt"

	"icw_common/enum"
	"icw_common/gen/activity"
	"icw_common/gen/core/biz"

	"icw_core_biz/repositories/mysql/model"
	"icw_core_biz/repositories/mysql/utils"
)

// ProjectDetectionSubTaskAggregateStatus 项目图像检测子任务聚合状态
type ProjectDetectionSubTaskAggregateStatus struct {
	AnyExecuted  bool
	AnyFailed    bool
	AllSucceeded bool
}

// ProjectDetectionSubTaskAggregateStatusTx 统计项目图像检测子任务聚合状态
func ProjectDetectionSubTaskAggregateStatusTx(ctx context.Context, tx *sql.Tx, mainTaskId uint64) (*ProjectDetectionSubTaskAggregateStatus, error) {
	task, err := FindProjectDetectionTaskByIdForUpdateTx(ctx, tx, mainTaskId)
	if err != nil {
		return nil, err
	}

	checks := []struct {
		shouldExecute bool
		taskId        sql.NullInt64
		statusFn      func(context.Context, *sql.Tx, uint64) (bizpb.ProjectDetectionSubTaskStatus_Value, error)
	}{{
		shouldExecute: task.CorrosionShouldExecute,
		taskId:        task.CorrosionTaskId,
		statusFn:      FindProjectDetectionCorrosionTaskStatusByIdTx,
	}, {
		shouldExecute: task.CrackShouldExecute,
		taskId:        task.CrackTaskId,
		statusFn:      FindProjectDetectionCrackTaskStatusByIdTx,
	}, {
		shouldExecute: task.StainShouldExecute,
		taskId:        task.StainTaskId,
		statusFn:      FindProjectDetectionStainTaskStatusByIdTx,
	}, {
		shouldExecute: task.FlatnessShouldExecute,
		taskId:        task.FlatnessTaskId,
		statusFn:      FindProjectDetectionFlatnessTaskStatusByIdTx,
	}, {
		shouldExecute: task.SpallingShouldExecute,
		taskId:        task.SpallingTaskId,
		statusFn:      FindProjectDetectionSpallingTaskStatusByIdTx,
	}}

	aggregate := &ProjectDetectionSubTaskAggregateStatus{
		AllSucceeded: true,
	}
	for _, check := range checks {
		if !check.shouldExecute {
			continue
		}
		aggregate.AnyExecuted = true
		if !check.taskId.Valid {
			aggregate.AllSucceeded = false
			continue
		}
		status, err := check.statusFn(ctx, tx, uint64(check.taskId.Int64))
		if err != nil {
			return nil, err
		}
		switch status {
		case bizpb.ProjectDetectionSubTaskStatus_Failed:
			aggregate.AnyFailed = true
			aggregate.AllSucceeded = false
		case bizpb.ProjectDetectionSubTaskStatus_Succeeded:
		default:
			aggregate.AllSucceeded = false
		}
	}
	if !aggregate.AnyExecuted {
		aggregate.AllSucceeded = false
	}
	return aggregate, nil
}

// FindProjectDetectionTaskByIdForUpdateTx 按主任务 ID 查询并锁定项目图像检测主任务记录
func FindProjectDetectionTaskByIdForUpdateTx(ctx context.Context, tx *sql.Tx, taskId uint64) (*model.ProjectDetectionTaskRecord, error) {
	return utils.ScanProjectDetectionTask(tx.QueryRowContext(ctx, `
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
	return utils.ScanProjectDetectionTask(tx.QueryRowContext(ctx, `
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

// FindProjectDetectionSubTaskByIdFromTableTx 按子任务代码和子任务 ID 查询项目图像检测子任务记录
func FindProjectDetectionSubTaskByIdFromTableTx(ctx context.Context, tx *sql.Tx, table string, taskId uint64) (*model.ProjectDetectionSubTaskRecord, error) {
	return utils.ScanProjectDetectionSubTask(tx.QueryRowContext(ctx, fmt.Sprintf(`
		SELECT id, uuid, main_task_id, user_id, project_id, image_id, status, started_at, finished_at, created_at, updated_at
		FROM %s
		WHERE id = ?
		LIMIT 1
	`, table), taskId))
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

// FindProjectDetectionSubTaskByIdTx 按子任务代码和子任务 ID 查询项目图像检测子任务记录
func FindProjectDetectionSubTaskByIdTx(ctx context.Context, tx *sql.Tx, taskCode string, taskId uint64) (*model.ProjectDetectionSubTaskRecord, error) {
	switch enum.ParseDetectionTaskCode(taskCode) {
	case activitypb.DetectionTaskCode_Corrosion:
		return FindProjectDetectionSubTaskByIdFromTableTx(ctx, tx, "project_detection_corrosion_tasks", taskId)
	case activitypb.DetectionTaskCode_Crack:
		return FindProjectDetectionSubTaskByIdFromTableTx(ctx, tx, "project_detection_crack_tasks", taskId)
	case activitypb.DetectionTaskCode_Stain:
		return FindProjectDetectionSubTaskByIdFromTableTx(ctx, tx, "project_detection_stain_tasks", taskId)
	case activitypb.DetectionTaskCode_Flatness:
		return FindProjectDetectionSubTaskByIdFromTableTx(ctx, tx, "project_detection_flatness_tasks", taskId)
	case activitypb.DetectionTaskCode_Spalling:
		return FindProjectDetectionSubTaskByIdFromTableTx(ctx, tx, "project_detection_spalling_tasks", taskId)
	default:
		return nil, model.ErrUnsupportedDetectionTaskCode
	}
}

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
