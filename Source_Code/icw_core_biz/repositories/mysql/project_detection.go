package mysql

import (
	"context"
	"database/sql"
	"errors"

	"github.com/google/uuid"

	"icw_common/enum"
	"icw_common/gen/activity"
	"icw_common/gen/core/biz"
	"icw_common/gen/core/common"

	"icw_core_biz/repositories/mysql/model"
	"icw_core_biz/repositories/mysql/project_detection"
	"icw_core_biz/repositories/mysql/utils"
)

// CreateProjectDetectionTasks 按用户 ID 和项目 ID 创建项目图像检测主任务
func (r *Repository) CreateProjectDetectionTasks(ctx context.Context, userId, projectId uint64) ([]*model.ProjectDetectionTaskRecord, error) {
	return r.createProjectDetectionTasks(ctx, userId, projectId, nil)
}

// CreateProjectDetectionTasksByImageUuids 按用户 ID、项目 ID 和图像 UUID 列表创建项目图像检测主任务
func (r *Repository) CreateProjectDetectionTasksByImageUuids(ctx context.Context, userId, projectId uint64, imageUuids []string) ([]*model.ProjectDetectionTaskRecord, error) {
	return r.createProjectDetectionTasks(ctx, userId, projectId, imageUuids)
}

// createProjectDetectionTasks 按用户 ID、项目 ID 和可选图像 UUID 列表创建项目图像检测主任务
func (r *Repository) createProjectDetectionTasks(ctx context.Context, userId, projectId uint64, imageUuids []string) ([]*model.ProjectDetectionTaskRecord, error) {
	images, err := r.ListProjectImages(ctx, userId, projectId)
	if err != nil {
		return nil, err
	}
	imageUuidSet := make(map[string]bool, len(imageUuids))
	for _, imageUuid := range imageUuids {
		if imageUuid != "" {
			imageUuidSet[imageUuid] = true
		}
	}

	tx, err := r.mysql.BeginTx(ctx, nil)
	if err != nil {
		return nil, err
	}
	defer func() {
		_ = tx.Rollback()
	}()

	taskIds := make([]uint64, 0, len(images))
	for _, image := range images {
		if image == nil || image.Status != bizpb.ProjectImageStatus_Uploaded {
			continue
		}
		if len(imageUuidSet) > 0 && !imageUuidSet[image.Uuid] {
			continue
		}
		result, err := tx.ExecContext(ctx, `
			INSERT INTO project_detection_tasks(uuid, user_id, project_id, image_id, image_uuid, status)
			VALUES (?, ?, ?, ?, ?, ?)
		`, uuid.NewString(), userId, projectId, image.Id, image.Uuid, enum.ProjectDetectionTaskStatusString(bizpb.ProjectDetectionTaskStatus_Pending))
		if err != nil {
			return nil, err
		}
		taskId, err := result.LastInsertId()
		if err != nil {
			return nil, err
		}
		taskIds = append(taskIds, uint64(taskId))
	}

	if err := tx.Commit(); err != nil {
		return nil, err
	}

	tasks := make([]*model.ProjectDetectionTaskRecord, 0, len(taskIds))
	for _, taskId := range taskIds {
		task, err := r.FindProjectDetectionTaskById(ctx, taskId)
		if err != nil {
			return nil, err
		}
		if task != nil {
			tasks = append(tasks, task)
		}
	}
	return tasks, nil
}

// ListProjectDetectionImageUuids 按用户 ID 和项目 ID 查询已创建检测任务的图像 UUID 列表
func (r *Repository) ListProjectDetectionImageUuids(ctx context.Context, userId, projectId uint64) ([]string, error) {
	rows, err := r.mysql.QueryContext(ctx, `
		SELECT image_uuid
		FROM project_detection_tasks
		WHERE user_id = ? AND project_id = ?
		ORDER BY created_at ASC, id ASC
	`, userId, projectId)

	if err != nil {
		return nil, err
	}
	defer func() {
		_ = rows.Close()
	}()

	imageUuids := make([]string, 0)
	for rows.Next() {
		var imageUuid string
		if err := rows.Scan(&imageUuid); err != nil {
			return nil, err
		}
		imageUuids = append(imageUuids, imageUuid)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return imageUuids, nil
}

// ListProjectDetectionImageUuidsByStatus 按用户 ID、项目 ID 和状态查询已创建检测任务的图像 UUID 列表
func (r *Repository) ListProjectDetectionImageUuidsByStatus(ctx context.Context, userId, projectId uint64, status bizpb.ProjectDetectionTaskStatus_Value) ([]string, error) {
	statusText := enum.ProjectDetectionTaskStatusString(status)
	if statusText == "" {
		return nil, model.ErrProjectDetectionTaskStatusInvalid
	}
	rows, err := r.mysql.QueryContext(ctx, `
		SELECT image_uuid
		FROM project_detection_tasks
		WHERE user_id = ? AND project_id = ? AND status = ?
		ORDER BY created_at ASC, id ASC
	`, userId, projectId, statusText)

	if err != nil {
		return nil, err
	}
	defer func() {
		_ = rows.Close()
	}()

	imageUuids := make([]string, 0)
	for rows.Next() {
		var imageUuid string
		if err := rows.Scan(&imageUuid); err != nil {
			return nil, err
		}
		imageUuids = append(imageUuids, imageUuid)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return imageUuids, nil
}

// DeleteProjectDetectionTasks 按用户 ID 和项目 ID 删除项目图像检测任务
func (r *Repository) DeleteProjectDetectionTasks(ctx context.Context, userId, projectId uint64) error {
	_, err := r.mysql.ExecContext(ctx, `
		DELETE FROM project_detection_tasks
		WHERE user_id = ? AND project_id = ?
	`, userId, projectId)
	return err
}

// DeleteProjectDetectionTasksByStatus 按用户 ID、项目 ID 和状态删除项目图像检测任务
func (r *Repository) DeleteProjectDetectionTasksByStatus(ctx context.Context, userId, projectId uint64, status bizpb.ProjectDetectionTaskStatus_Value) error {
	statusText := enum.ProjectDetectionTaskStatusString(status)
	if statusText == "" {
		return model.ErrProjectDetectionTaskStatusInvalid
	}
	_, err := r.mysql.ExecContext(ctx, `
		DELETE FROM project_detection_tasks
		WHERE user_id = ? AND project_id = ? AND status = ?
	`, userId, projectId, statusText)
	return err
}

// ListProjectDetectionTasksByStatus 按用户 ID、项目 ID 和状态查询项目图像检测主任务
func (r *Repository) ListProjectDetectionTasksByStatus(ctx context.Context, userId, projectId uint64, status bizpb.ProjectDetectionTaskStatus_Value) ([]*model.ProjectDetectionTaskRecord, error) {
	rows, err := r.mysql.QueryContext(ctx, `
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
		WHERE user_id = ? AND project_id = ? AND status = ?
		ORDER BY created_at ASC, id ASC
	`, userId, projectId, enum.ProjectDetectionTaskStatusString(status))

	if err != nil {
		return nil, err
	}
	defer func() {
		_ = rows.Close()
	}()

	tasks := make([]*model.ProjectDetectionTaskRecord, 0)
	for rows.Next() {
		task, err := utils.ScanProjectDetectionTask(rows)
		if err != nil {
			return nil, err
		}
		tasks = append(tasks, task)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return tasks, nil
}

// FindProjectDetectionTaskById 按主任务 ID 查询项目图像检测主任务
func (r *Repository) FindProjectDetectionTaskById(ctx context.Context, taskId uint64) (*model.ProjectDetectionTaskRecord, error) {
	task, err := utils.ScanProjectDetectionTask(r.mysql.QueryRowContext(ctx, `
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
		LIMIT 1
	`, taskId))
	if errors.Is(err, sql.ErrNoRows) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return task, nil
}

// GetProjectDetectionTasks 按用户 ID 和项目 ID 查询项目图像检测任务状态
func (r *Repository) GetProjectDetectionTasks(ctx context.Context, userId, projectId uint64) ([]*commonpb.ProjectDetectionStatus, error) {
	rows, err := r.mysql.QueryContext(ctx, `
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
		WHERE user_id = ? AND project_id = ?
		ORDER BY created_at ASC, id ASC
	`, userId, projectId)

	if err != nil {
		return nil, err
	}
	defer func() {
		_ = rows.Close()
	}()

	tasks := make([]*model.ProjectDetectionTaskRecord, 0)
	for rows.Next() {
		task, err := utils.ScanProjectDetectionTask(rows)
		if err != nil {
			return nil, err
		}
		tasks = append(tasks, task)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	tx, err := r.mysql.BeginTx(ctx, nil)
	if err != nil {
		return nil, err
	}
	defer func() {
		_ = tx.Rollback()
	}()

	items := make([]*commonpb.ProjectDetectionStatus, 0, len(tasks))
	for _, task := range tasks {
		item, err := projectDetectionTaskToStatusDTO(ctx, tx, task)
		if err != nil {
			return nil, err
		}
		items = append(items, item)
	}
	if err := tx.Commit(); err != nil {
		return nil, err
	}
	return items, nil
}

// UpdateProjectDetectionTaskStatus 按主任务 ID 更新项目图像检测主任务状态
func (r *Repository) UpdateProjectDetectionTaskStatus(ctx context.Context, taskId uint64, status bizpb.ProjectDetectionTaskStatus_Value) (*model.ProjectDetectionTaskRecord, error) {
	statusText := enum.ProjectDetectionTaskStatusString(status)
	if statusText == "" {
		return nil, model.ErrProjectDetectionTaskStatusInvalid
	}
	result, err := r.mysql.ExecContext(ctx, `
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
	if err := utils.CheckRowsAffected(result, err); err != nil {
		return nil, err
	}
	return r.FindProjectDetectionTaskById(ctx, taskId)
}

// StartProjectDetectionClassification 按主任务 ID 抢占项目图像检测主任务并推进到分类阶段
func (r *Repository) StartProjectDetectionClassification(ctx context.Context, taskId uint64) (*model.ProjectDetectionTaskRecord, error) {
	result, err := r.mysql.ExecContext(ctx, `
		UPDATE project_detection_tasks
		SET status = ?, started_at = NOW(3)
		WHERE id = ? AND status = ?
	`, enum.ProjectDetectionTaskStatusString(bizpb.ProjectDetectionTaskStatus_Classifying),
		taskId,
		enum.ProjectDetectionTaskStatusString(bizpb.ProjectDetectionTaskStatus_Pending))
	if err != nil {
		return nil, err
	}
	affected, err := result.RowsAffected()
	if err != nil {
		return nil, err
	}
	if affected == 0 {
		return nil, nil
	}
	return r.FindProjectDetectionTaskById(ctx, taskId)
}

// UpdateProjectDetectionClassificationResult 按主任务 UUID 更新项目图像检测分类结果
func (r *Repository) UpdateProjectDetectionClassificationResult(ctx context.Context, taskUuid string, status bizpb.ProjectDetectionSubTaskStatus_Value, taskCodes []string) (*model.ProjectDetectionTaskRecord, map[string]*model.ProjectDetectionSubTaskRecord, error) {
	tx, err := r.mysql.BeginTx(ctx, nil)
	if err != nil {
		return nil, nil, err
	}
	defer func() {
		_ = tx.Rollback()
	}()

	mainTask, err := project_detection.FindProjectDetectionTaskByUuidTx(ctx, tx, taskUuid)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, nil, nil
	}
	if err != nil {
		return nil, nil, err
	}
	if mainTask.Status != bizpb.ProjectDetectionTaskStatus_Classifying {
		return mainTask, nil, tx.Commit()
	}

	if status == bizpb.ProjectDetectionSubTaskStatus_Failed {
		if err := project_detection.UpdateProjectDetectionTaskStatusTx(ctx, tx, mainTask.Id, bizpb.ProjectDetectionTaskStatus_Failed); err != nil {
			return nil, nil, err
		}
		if err := tx.Commit(); err != nil {
			return nil, nil, err
		}
		task, err := r.FindProjectDetectionTaskById(ctx, mainTask.Id)
		return task, nil, err
	}

	taskCodes, err = utils.NormalizeDetectionTaskCodes(taskCodes)
	if err != nil {
		return nil, nil, err
	}
	if len(taskCodes) == 0 {
		if err := project_detection.UpdateProjectDetectionTaskStatusTx(ctx, tx, mainTask.Id, bizpb.ProjectDetectionTaskStatus_Succeeded); err != nil {
			return nil, nil, err
		}
		if err := tx.Commit(); err != nil {
			return nil, nil, err
		}
		task, err := r.FindProjectDetectionTaskById(ctx, mainTask.Id)
		return task, nil, err
	}

	subTasks := make(map[string]*model.ProjectDetectionSubTaskRecord, len(taskCodes))
	for _, taskCode := range taskCodes {
		// 创建 检测任务 表
		subTask, err := project_detection.CreateProjectDetectionSubTaskTx(ctx, tx, mainTask, taskCode)
		if err != nil {
			return nil, nil, err
		}
		subTasks[taskCode] = subTask
	}

	if err := tx.Commit(); err != nil {
		return nil, nil, err
	}
	task, err := r.FindProjectDetectionTaskById(ctx, mainTask.Id)
	return task, subTasks, err
}

// UpdateProjectDetectionReasoningTaskResult 按推理任务 UUID 更新项目图像检测推理子任务结果
func (r *Repository) UpdateProjectDetectionReasoningTaskResult(ctx context.Context, taskCode, taskUuid string, status bizpb.ProjectDetectionSubTaskStatus_Value, resultJSON string) (*model.ProjectDetectionTaskRecord, *model.ProjectDetectionSubTaskRecord, *model.ProjectDetectionSummaryTaskRecord, error) {
	statusText := enum.ProjectDetectionSubTaskStatusString(status)
	if statusText == "" {
		return nil, nil, nil, model.ErrProjectDetectionSubTaskStatusInvalid
	}

	tx, err := r.mysql.BeginTx(ctx, nil)
	if err != nil {
		return nil, nil, nil, err
	}
	defer func() {
		_ = tx.Rollback()
	}()

	// 根据uuid查询 MainTask ID
	subTask, err := project_detection.FindProjectDetectionSubTaskByUuidTx(ctx, tx, taskCode, taskUuid)
	if err != nil || subTask == nil {
		return nil, nil, nil, err
	}
	if subTask.Status != bizpb.ProjectDetectionSubTaskStatus_Pending {
		if err := tx.Commit(); err != nil {
			return nil, nil, nil, err
		}
		task, err := r.FindProjectDetectionTaskById(ctx, subTask.MainTaskId)
		return task, subTask, nil, err
	}
	// 写入检测任务 更新检测任务的状态
	if err := project_detection.UpdateProjectDetectionSubTaskResultTx(ctx, tx, taskCode, taskUuid, status, resultJSON); err != nil {
		return nil, nil, nil, err
	}

	var summaryTask *model.ProjectDetectionSummaryTaskRecord
	if status == bizpb.ProjectDetectionSubTaskStatus_Failed {
		if err := project_detection.UpdateProjectDetectionTaskStatusTx(ctx, tx, subTask.MainTaskId, bizpb.ProjectDetectionTaskStatus_Failed); err != nil {
			return nil, nil, nil, err
		}
	} else {
		anyFailed, allSucceeded, err := project_detection.ProjectDetectionSubTaskAggregateStatusTx(ctx, tx, subTask.MainTaskId)
		if err != nil {
			return nil, nil, nil, err
		}
		if anyFailed {
			if err := project_detection.UpdateProjectDetectionTaskStatusTx(ctx, tx, subTask.MainTaskId, bizpb.ProjectDetectionTaskStatus_Failed); err != nil {
				return nil, nil, nil, err
			}
		} else if allSucceeded {
			mainTask, err := project_detection.FindProjectDetectionTaskByIdForUpdateTx(ctx, tx, subTask.MainTaskId)
			if err != nil {
				return nil, nil, nil, err
			}
			if !mainTask.SummaryShouldExecute && !mainTask.SummaryTaskId.Valid {
				summaryTask, err = project_detection.CreateProjectDetectionSummaryTaskTx(ctx, tx, mainTask)
				if err != nil {
					return nil, nil, nil, err
				}
			}
		}
	}

	if err := tx.Commit(); err != nil {
		return nil, nil, nil, err
	}
	task, err := r.FindProjectDetectionTaskById(ctx, subTask.MainTaskId)
	return task, subTask, summaryTask, err
}

// StartProjectDetectionReasoningTask 按推理任务 UUID 标记项目图像检测推理子任务开始执行
func (r *Repository) StartProjectDetectionReasoningTask(ctx context.Context, taskCode, taskUuid string) (*model.ProjectDetectionTaskRecord, *model.ProjectDetectionSubTaskRecord, error) {
	tx, err := r.mysql.BeginTx(ctx, nil)
	if err != nil {
		return nil, nil, err
	}
	defer func() {
		_ = tx.Rollback()
	}()

	subTask, err := project_detection.FindProjectDetectionSubTaskByUuidTx(ctx, tx, taskCode, taskUuid)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, nil, nil
	}
	if err != nil {
		return nil, nil, err
	}
	if subTask.Status == bizpb.ProjectDetectionSubTaskStatus_Pending {
		if err := project_detection.UpdateProjectDetectionSubTaskStatusTx(ctx, tx, taskCode, taskUuid, bizpb.ProjectDetectionSubTaskStatus_Pending, true, false); err != nil {
			return nil, nil, err
		}
	}
	if err := tx.Commit(); err != nil {
		return nil, nil, err
	}

	task, err := r.FindProjectDetectionTaskById(ctx, subTask.MainTaskId)
	return task, subTask, err
}

// StartProjectDetectionSummaryTask 按总结任务 UUID 标记图像检测总结任务开始执行
func (r *Repository) StartProjectDetectionSummaryTask(ctx context.Context, taskUuid string) (*model.ProjectDetectionTaskRecord, *model.ProjectDetectionSummaryTaskRecord, error) {
	tx, err := r.mysql.BeginTx(ctx, nil)
	if err != nil {
		return nil, nil, err
	}
	defer func() {
		_ = tx.Rollback()
	}()

	subTask, err := project_detection.FindProjectDetectionSummaryTaskByUuidTx(ctx, tx, taskUuid)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, nil, nil
	}
	if err != nil {
		return nil, nil, err
	}
	if subTask.Status == bizpb.ProjectDetectionSubTaskStatus_Pending {
		if err := project_detection.UpdateProjectDetectionSummaryTaskStatusTx(ctx, tx, taskUuid, bizpb.ProjectDetectionSubTaskStatus_Pending, true, false); err != nil {
			return nil, nil, err
		}
	}
	if err := tx.Commit(); err != nil {
		return nil, nil, err
	}

	task, err := r.FindProjectDetectionTaskById(ctx, subTask.MainTaskId)
	return task, subTask, err
}

// UpdateProjectDetectionSummaryResult 按总结任务 UUID 更新项目图像检测总结结果
func (r *Repository) UpdateProjectDetectionSummaryResult(ctx context.Context, taskUuid string, status bizpb.ProjectDetectionSubTaskStatus_Value) (*model.ProjectDetectionTaskRecord, *model.ProjectDetectionSummaryTaskRecord, error) {
	statusText := enum.ProjectDetectionSubTaskStatusString(status)
	if statusText == "" {
		return nil, nil, model.ErrProjectDetectionSummaryTaskStatusInvalid
	}

	tx, err := r.mysql.BeginTx(ctx, nil)
	if err != nil {
		return nil, nil, err
	}
	defer func() {
		_ = tx.Rollback()
	}()

	subTask, err := project_detection.FindProjectDetectionSummaryTaskByUuidTx(ctx, tx, taskUuid)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, nil, nil
	}
	if err != nil {
		return nil, nil, err
	}
	if subTask.Status != bizpb.ProjectDetectionSubTaskStatus_Pending {
		if err := tx.Commit(); err != nil {
			return nil, nil, err
		}
		task, err := r.FindProjectDetectionTaskById(ctx, subTask.MainTaskId)
		return task, subTask, err
	}
	if err := project_detection.UpdateProjectDetectionSummaryTaskStatusTx(ctx, tx, taskUuid, status, false, true); err != nil {
		return nil, nil, err
	}

	mainStatus := bizpb.ProjectDetectionTaskStatus_Succeeded
	if status == bizpb.ProjectDetectionSubTaskStatus_Failed {
		mainStatus = bizpb.ProjectDetectionTaskStatus_Failed
	}
	if err := project_detection.UpdateProjectDetectionTaskStatusTx(ctx, tx, subTask.MainTaskId, mainStatus); err != nil {
		return nil, nil, err
	}

	if err := tx.Commit(); err != nil {
		return nil, nil, err
	}
	task, err := r.FindProjectDetectionTaskById(ctx, subTask.MainTaskId)
	return task, subTask, err
}

// projectDetectionTaskToStatusDTO 将项目图像检测主任务记录转换为检测状态 DTO
func projectDetectionTaskToStatusDTO(ctx context.Context, tx *sql.Tx, task *model.ProjectDetectionTaskRecord) (*commonpb.ProjectDetectionStatus, error) {
	item := &commonpb.ProjectDetectionStatus{
		ImageUuid:  task.ImageUuid,
		MainTaskId: task.Uuid,
		MainStatus: enum.ProjectDetectionTaskStatusString(task.Status),
		Nodes:      make([]*commonpb.ProjectDetectionNodeStatus, 0, 7),
	}

	if subStatus := classificationNodeStatus(task); subStatus != "" {
		item.Nodes = append(item.Nodes, &commonpb.ProjectDetectionNodeStatus{
			NodeCode:  "classification",
			SubStatus: subStatus,
		})
	}

	reasoningNodes := []struct {
		taskCode      string
		shouldExecute bool
		taskId        sql.NullInt64
	}{
		{taskCode: enum.DetectionTaskCodeString(activitypb.DetectionTaskCode_Corrosion), shouldExecute: task.CorrosionShouldExecute, taskId: task.CorrosionTaskId},
		{taskCode: enum.DetectionTaskCodeString(activitypb.DetectionTaskCode_Crack), shouldExecute: task.CrackShouldExecute, taskId: task.CrackTaskId},
		{taskCode: enum.DetectionTaskCodeString(activitypb.DetectionTaskCode_Stain), shouldExecute: task.StainShouldExecute, taskId: task.StainTaskId},
		{taskCode: enum.DetectionTaskCodeString(activitypb.DetectionTaskCode_Flatness), shouldExecute: task.FlatnessShouldExecute, taskId: task.FlatnessTaskId},
		{taskCode: enum.DetectionTaskCodeString(activitypb.DetectionTaskCode_Spalling), shouldExecute: task.SpallingShouldExecute, taskId: task.SpallingTaskId},
	}
	for _, node := range reasoningNodes {
		if !node.shouldExecute || !node.taskId.Valid {
			continue
		}
		subTask, err := project_detection.FindProjectDetectionSubTaskByIdTx(ctx, tx, node.taskCode, uint64(node.taskId.Int64))
		if err != nil {
			return nil, err
		}
		item.Nodes = append(item.Nodes, &commonpb.ProjectDetectionNodeStatus{
			NodeCode:  "reasoning:" + node.taskCode,
			SubTaskId: subTask.Uuid,
			SubStatus: enum.ProjectDetectionSubTaskStatusString(subTask.Status),
		})
	}

	if task.SummaryShouldExecute && task.SummaryTaskId.Valid {
		subTask, err := project_detection.FindProjectDetectionSubTaskByIdFromTableTx(ctx, tx, "project_detection_summary_tasks", uint64(task.SummaryTaskId.Int64))
		if err != nil {
			return nil, err
		}
		item.Nodes = append(item.Nodes, &commonpb.ProjectDetectionNodeStatus{
			NodeCode:  "summary",
			SubTaskId: subTask.Uuid,
			SubStatus: enum.ProjectDetectionSubTaskStatusString(subTask.Status),
		})
	}

	return item, nil
}

// classificationNodeStatus 根据主任务状态推导分类节点状态
func classificationNodeStatus(task *model.ProjectDetectionTaskRecord) string {
	switch task.Status {
	case bizpb.ProjectDetectionTaskStatus_Classifying:
		return enum.ProjectDetectionSubTaskStatusString(bizpb.ProjectDetectionSubTaskStatus_Pending)
	case bizpb.ProjectDetectionTaskStatus_Detecting,
		bizpb.ProjectDetectionTaskStatus_Summarizing,
		bizpb.ProjectDetectionTaskStatus_Succeeded:
		return enum.ProjectDetectionSubTaskStatusString(bizpb.ProjectDetectionSubTaskStatus_Succeeded)
	case bizpb.ProjectDetectionTaskStatus_Failed:
		if !task.CorrosionShouldExecute &&
			!task.CrackShouldExecute &&
			!task.StainShouldExecute &&
			!task.FlatnessShouldExecute &&
			!task.SpallingShouldExecute &&
			!task.SummaryShouldExecute {
			return enum.ProjectDetectionSubTaskStatusString(bizpb.ProjectDetectionSubTaskStatus_Failed)
		}
		return enum.ProjectDetectionSubTaskStatusString(bizpb.ProjectDetectionSubTaskStatus_Succeeded)
	default:
		return ""
	}
}
