package mysql

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"

	"github.com/google/uuid"

	"icw_common/enum"
	"icw_common/gen/core/biz"
	"icw_common/gen/core/common"

	"icw_core_biz/repositories/mysql/model"
	"icw_core_biz/repositories/mysql/project_detection"
	"icw_core_biz/repositories/mysql/utils"
)

// CreateProjectDetectionTasks 按用户 ID 和项目 ID 创建项目图像检测主任务
func (r *Repository) CreateProjectDetectionTasks(ctx context.Context, userId, projectId uint64) ([]*model.ProjectDetectionTaskRecord, error) {
	images, err := r.ListProjectImages(ctx, userId, projectId)
	if err != nil {
		return nil, err
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

// RetryProjectDetectionTasks 按用户 ID 和项目 ID 重试项目图像检测主任务
func (r *Repository) RetryProjectDetectionTasks(ctx context.Context, userId, projectId uint64) ([]string, []*model.ProjectDetectionTaskRecord, error) {
	tx, err := r.mysql.BeginTx(ctx, nil)
	if err != nil {
		return nil, nil, err
	}
	defer func() {
		_ = tx.Rollback()
	}()

	rows, err := tx.QueryContext(ctx, `
		SELECT
			id,
			image_id,
			image_uuid,
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
			summary_task_id
		FROM project_detection_tasks
		WHERE user_id = ? AND project_id = ? AND status = ?
		ORDER BY created_at ASC, id ASC
		FOR UPDATE
	`, userId, projectId, enum.ProjectDetectionTaskStatusString(bizpb.ProjectDetectionTaskStatus_Failed))
	if err != nil {
		return nil, nil, err
	}

	type failedTask struct {
		id                     uint64
		imageId                uint64
		imageUuid              string
		corrosionShouldExecute bool
		corrosionTaskId        sql.NullInt64
		crackShouldExecute     bool
		crackTaskId            sql.NullInt64
		stainShouldExecute     bool
		stainTaskId            sql.NullInt64
		flatnessShouldExecute  bool
		flatnessTaskId         sql.NullInt64
		spallingShouldExecute  bool
		spallingTaskId         sql.NullInt64
		summaryShouldExecute   bool
		summaryTaskId          sql.NullInt64
	}

	failedTasks := make([]*failedTask, 0)
	for rows.Next() {
		task := &failedTask{}
		if err := rows.Scan(
			&task.id,
			&task.imageId,
			&task.imageUuid,
			&task.corrosionShouldExecute,
			&task.corrosionTaskId,
			&task.crackShouldExecute,
			&task.crackTaskId,
			&task.stainShouldExecute,
			&task.stainTaskId,
			&task.flatnessShouldExecute,
			&task.flatnessTaskId,
			&task.spallingShouldExecute,
			&task.spallingTaskId,
			&task.summaryShouldExecute,
			&task.summaryTaskId,
		); err != nil {
			_ = rows.Close()
			return nil, nil, err
		}
		failedTasks = append(failedTasks, task)
	}

	if err := rows.Err(); err != nil {
		_ = rows.Close()
		return nil, nil, err
	}
	if err := rows.Close(); err != nil {
		return nil, nil, err
	}
	if len(failedTasks) == 0 {
		if err := tx.Commit(); err != nil {
			return nil, nil, err
		}
		return nil, nil, nil
	}

	imageUuids := make([]string, 0, len(failedTasks))
	taskIds := make([]uint64, 0, len(failedTasks))
	deleteSubTask := func(table string, shouldExecute bool, taskId sql.NullInt64, mainTaskId uint64) error {
		if !shouldExecute || !taskId.Valid {
			return nil
		}
		_, err := tx.ExecContext(ctx, fmt.Sprintf(`
			DELETE FROM %s
			WHERE id = ? AND main_task_id = ? AND user_id = ? AND project_id = ?
		`, table), taskId.Int64, mainTaskId, userId, projectId)
		return err
	}

	for _, task := range failedTasks {
		imageUuids = append(imageUuids, task.imageUuid)

		if err := deleteSubTask("project_detection_corrosion_tasks", task.corrosionShouldExecute, task.corrosionTaskId, task.id); err != nil {
			return nil, nil, err
		}
		if err := deleteSubTask("project_detection_crack_tasks", task.crackShouldExecute, task.crackTaskId, task.id); err != nil {
			return nil, nil, err
		}
		if err := deleteSubTask("project_detection_stain_tasks", task.stainShouldExecute, task.stainTaskId, task.id); err != nil {
			return nil, nil, err
		}
		if err := deleteSubTask("project_detection_flatness_tasks", task.flatnessShouldExecute, task.flatnessTaskId, task.id); err != nil {
			return nil, nil, err
		}
		if err := deleteSubTask("project_detection_spalling_tasks", task.spallingShouldExecute, task.spallingTaskId, task.id); err != nil {
			return nil, nil, err
		}
		if err := deleteSubTask("project_detection_summary_tasks", task.summaryShouldExecute, task.summaryTaskId, task.id); err != nil {
			return nil, nil, err
		}

		result, err := tx.ExecContext(ctx, `
			DELETE FROM project_detection_tasks
			WHERE id = ? AND user_id = ? AND project_id = ? AND status = ?
		`, task.id, userId, projectId, enum.ProjectDetectionTaskStatusString(bizpb.ProjectDetectionTaskStatus_Failed))
		if err := utils.CheckRowsAffected(result, err); err != nil {
			return nil, nil, err
		}

		result, err = tx.ExecContext(ctx, `
			INSERT INTO project_detection_tasks(uuid, user_id, project_id, image_id, image_uuid, status)
			SELECT ?, user_id, project_id, id, uuid, ?
			FROM project_group_images
			WHERE id = ? AND user_id = ? AND project_id = ? AND status = ?
			LIMIT 1
		`, uuid.NewString(), enum.ProjectDetectionTaskStatusString(bizpb.ProjectDetectionTaskStatus_Pending), task.imageId, userId, projectId, enum.ProjectImageStatusString(bizpb.ProjectImageStatus_Uploaded))
		if err != nil {
			return nil, nil, err
		}

		affected, err := result.RowsAffected()
		if err != nil {
			return nil, nil, err
		}
		if affected == 0 {
			continue
		}

		taskId, err := result.LastInsertId()
		if err != nil {
			return nil, nil, err
		}

		taskIds = append(taskIds, uint64(taskId))
	}

	if err := tx.Commit(); err != nil {
		return nil, nil, err
	}

	tasks := make([]*model.ProjectDetectionTaskRecord, 0, len(taskIds))
	for _, taskId := range taskIds {
		task, err := r.FindProjectDetectionTaskById(ctx, taskId)
		if err != nil {
			return nil, nil, err
		}
		if task != nil {
			tasks = append(tasks, task)
		}
	}

	return imageUuids, tasks, nil
}

// FindProjectDetectionTaskById 按主任务 ID 查询项目图像检测主任务
func (r *Repository) FindProjectDetectionTaskById(ctx context.Context, taskId uint64) (*model.ProjectDetectionTaskRecord, error) {
	task, err := model.ScanProjectDetectionTask(r.mysql.QueryRowContext(ctx, `
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

// FindProjectDetectionTaskByImageUuid 按用户 ID、项目 ID 和图像 UUID 查询项目图像检测主任务
func (r *Repository) FindProjectDetectionTaskByImageUuid(ctx context.Context, userId, projectId uint64, imageUuid string) (*model.ProjectDetectionTaskRecord, error) {
	task, err := model.ScanProjectDetectionTask(r.mysql.QueryRowContext(ctx, `
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
		WHERE user_id = ? AND project_id = ? AND image_uuid = ?
		LIMIT 1
	`, userId, projectId, strings.TrimSpace(imageUuid)))
	if errors.Is(err, sql.ErrNoRows) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return task, nil
}

// GetProjectDetectionTasksStatus 按用户 ID 和项目 ID 查询项目图像检测任务状态
func (r *Repository) GetProjectDetectionTasksStatus(ctx context.Context, userId, projectId uint64) ([]*commonpb.ProjectDetectionStatus, error) {
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
		task, err := model.ScanProjectDetectionTask(rows)
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
		item, err := r.projectDetectionTaskToStatusDTO(ctx, tx, task)
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

// GetProjectDetectionTaskStatus 按用户 ID、项目 ID 和图像 UUID 查询项目图像检测任务状态
func (r *Repository) GetProjectDetectionTaskStatus(ctx context.Context, userId, projectId uint64, imageUuid string) (*commonpb.ProjectDetectionStatus, error) {
	task, err := r.FindProjectDetectionTaskByImageUuid(ctx, userId, projectId, imageUuid)
	if err != nil || task == nil {
		return nil, err
	}

	tx, err := r.mysql.BeginTx(ctx, nil)
	if err != nil {
		return nil, err
	}
	defer func() {
		_ = tx.Rollback()
	}()

	item, err := r.projectDetectionTaskToStatusDTO(ctx, tx, task)
	if err != nil {
		return nil, err
	}

	if err := tx.Commit(); err != nil {
		return nil, err
	}

	return item, nil
}

// projectDetectionTaskToStatusDTO 将项目图像检测主任务记录转换为检测状态数据结构
func (r *Repository) projectDetectionTaskToStatusDTO(ctx context.Context, tx *sql.Tx, task *model.ProjectDetectionTaskRecord) (*commonpb.ProjectDetectionStatus, error) {
	item := &commonpb.ProjectDetectionStatus{
		ImageUuid:       task.ImageUuid,
		MainTaskUuid:    task.Uuid,
		MainStatus:      task.Status,
		DetectionStatus: make([]*commonpb.ProjectDetectionNodeStatus, 0),
		StartedAt:       utils.NullTimeString(task.StartedAt),
		FinishedAt:      utils.NullTimeString(task.FinishedAt),
	}

	if subStatus := project_detection.ClassificationNodeStatus(task); subStatus != bizpb.ProjectDetectionSubTaskStatus_Unknown {
		item.ClassificationStatus = &commonpb.ProjectDetectionNodeStatus{
			NodeCode:  "classification",
			SubStatus: subStatus,
		}
	}

	reasoningNodes := []struct {
		taskCode      string
		shouldExecute bool
		taskId        sql.NullInt64
	}{{
		taskCode:      enum.DetectionTaskCodeString(commonpb.DetectionTaskCode_Corrosion),
		shouldExecute: task.CorrosionShouldExecute,
		taskId:        task.CorrosionTaskId,
	}, {
		taskCode:      enum.DetectionTaskCodeString(commonpb.DetectionTaskCode_Crack),
		shouldExecute: task.CrackShouldExecute,
		taskId:        task.CrackTaskId,
	}, {
		taskCode:      enum.DetectionTaskCodeString(commonpb.DetectionTaskCode_Stain),
		shouldExecute: task.StainShouldExecute,
		taskId:        task.StainTaskId,
	}, {
		taskCode:      enum.DetectionTaskCodeString(commonpb.DetectionTaskCode_Flatness),
		shouldExecute: task.FlatnessShouldExecute,
		taskId:        task.FlatnessTaskId,
	}, {
		taskCode:      enum.DetectionTaskCodeString(commonpb.DetectionTaskCode_Spalling),
		shouldExecute: task.SpallingShouldExecute,
		taskId:        task.SpallingTaskId,
	}}

	for _, node := range reasoningNodes {
		if !node.shouldExecute || !node.taskId.Valid {
			continue
		}
		subTask, err := project_detection.FindProjectDetectionSubTaskByIdTx(ctx, tx, node.taskCode, uint64(node.taskId.Int64))
		if err != nil {
			return nil, err
		}
		item.DetectionStatus = append(item.DetectionStatus, &commonpb.ProjectDetectionNodeStatus{
			NodeCode:    "reasoning:" + node.taskCode,
			SubTaskUuid: subTask.Uuid,
			SubStatus:   subTask.Status,
		})
	}

	if task.SummaryShouldExecute && task.SummaryTaskId.Valid {
		subTask, err := project_detection.FindProjectDetectionSubTaskByIdFromTableTx(ctx, tx, "project_detection_summary_tasks", uint64(task.SummaryTaskId.Int64))
		if err != nil {
			return nil, err
		}
		item.SummaryStatus = &commonpb.ProjectDetectionNodeStatus{
			NodeCode:    "summary",
			SubTaskUuid: subTask.Uuid,
			SubStatus:   subTask.Status,
		}
	}

	return item, nil
}

// ProjectDetectionTasksAllSucceeded 按用户 ID 和项目 ID 判断项目图像检测任务是否全部成功
func (r *Repository) ProjectDetectionTasksAllSucceeded(ctx context.Context, userId, projectId uint64) (bool, error) {
	var uploadedImageCount, succeededTaskCount uint64
	if err := r.mysql.QueryRowContext(ctx, `
		SELECT
			COUNT(i.id),
			COALESCE(SUM(CASE WHEN t.status = ? THEN 1 ELSE 0 END), 0)
		FROM project_group_images i
		LEFT JOIN project_detection_tasks t
			ON t.image_id = i.id
			AND t.user_id = i.user_id
			AND t.project_id = i.project_id
		WHERE i.user_id = ? AND i.project_id = ? AND i.status = ?
	`, enum.ProjectDetectionTaskStatusString(bizpb.ProjectDetectionTaskStatus_Succeeded),
		userId,
		projectId,
		enum.ProjectImageStatusString(bizpb.ProjectImageStatus_Uploaded)).Scan(&uploadedImageCount, &succeededTaskCount); err != nil {
		return false, err
	}
	return uploadedImageCount > 0 && uploadedImageCount == succeededTaskCount, nil
}

// GetProjectDetectionCorrosionResult 按项目图像检测子任务 ID 查询金属锈蚀检测结果
func (r *Repository) GetProjectDetectionCorrosionResult(ctx context.Context, taskId uint64) (*commonpb.ProjectDetectionCorrosionResult, error) {
	result := &commonpb.ProjectDetectionCorrosionResult{}

	var (
		status            string
		hasCorrosion      sql.NullBool
		corrosionCount    sql.NullInt64
		maxConfidence     sql.NullFloat64
		averageConfidence sql.NullFloat64
		corrosionPixels   sql.NullInt64
		corrosionRatio    sql.NullFloat64
		regions           sql.NullString
		runtimeSeconds    sql.NullFloat64
		startedAt         sql.NullTime
		finishedAt        sql.NullTime
	)

	if err := r.mysql.QueryRowContext(ctx, `
		SELECT
			uuid,
			status,
			has_corrosion,
			corrosion_count,
			max_confidence,
			average_confidence,
			corrosion_pixels,
			corrosion_ratio,
			CAST(regions AS CHAR),
			runtime_seconds,
			started_at,
			finished_at
		FROM project_detection_corrosion_tasks
		WHERE id = ?
		LIMIT 1
	`, taskId).Scan(
		&result.TaskUuid,
		&status,
		&hasCorrosion,
		&corrosionCount,
		&maxConfidence,
		&averageConfidence,
		&corrosionPixels,
		&corrosionRatio,
		&regions,
		&runtimeSeconds,
		&startedAt,
		&finishedAt,
	); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}

	result.Status = enum.ParseProjectDetectionSubTaskStatus(status)
	result.HasCorrosion = utils.NullBool(hasCorrosion)
	result.CorrosionCount = utils.NullUint32(corrosionCount)
	result.MaxConfidence = utils.NullFloat64(maxConfidence)
	result.AverageConfidence = utils.NullFloat64(averageConfidence)
	result.CorrosionPixels = utils.NullUint64(corrosionPixels)
	result.CorrosionRatio = utils.NullFloat64(corrosionRatio)
	if err := utils.UnmarshalRegions(regions, &result.Regions); err != nil {
		return nil, err
	}
	result.RuntimeSeconds = utils.NullFloat64(runtimeSeconds)
	result.StartedAt = utils.NullTimeString(startedAt)
	result.FinishedAt = utils.NullTimeString(finishedAt)

	return result, nil
}

// GetProjectDetectionCrackResult 按项目图像检测子任务 ID 查询石材裂缝检测结果
func (r *Repository) GetProjectDetectionCrackResult(ctx context.Context, taskId uint64) (*commonpb.ProjectDetectionCrackResult, error) {
	result := &commonpb.ProjectDetectionCrackResult{}

	var (
		status         string
		hasCrack       sql.NullBool
		crackCount     sql.NullInt64
		crackPixels    sql.NullInt64
		crackRatio     sql.NullFloat64
		regions        sql.NullString
		runtimeSeconds sql.NullFloat64
		startedAt      sql.NullTime
		finishedAt     sql.NullTime
	)

	if err := r.mysql.QueryRowContext(ctx, `
		SELECT
			uuid,
			status,
			has_crack,
			crack_count,
			crack_pixels,
			crack_ratio,
			CAST(regions AS CHAR),
			runtime_seconds,
			started_at,
			finished_at
		FROM project_detection_crack_tasks
		WHERE id = ?
		LIMIT 1
	`, taskId).Scan(
		&result.TaskUuid,
		&status,
		&hasCrack,
		&crackCount,
		&crackPixels,
		&crackRatio,
		&regions,
		&runtimeSeconds,
		&startedAt,
		&finishedAt,
	); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}

	result.Status = enum.ParseProjectDetectionSubTaskStatus(status)
	result.HasCrack = utils.NullBool(hasCrack)
	result.CrackCount = utils.NullUint32(crackCount)
	result.CrackPixels = utils.NullUint64(crackPixels)
	result.CrackRatio = utils.NullFloat64(crackRatio)
	if err := utils.UnmarshalRegions(regions, &result.Regions); err != nil {
		return nil, err
	}
	result.RuntimeSeconds = utils.NullFloat64(runtimeSeconds)
	result.StartedAt = utils.NullTimeString(startedAt)
	result.FinishedAt = utils.NullTimeString(finishedAt)

	return result, nil
}

// GetProjectDetectionStainResult 按项目图像检测子任务 ID 查询石材污渍检测结果
func (r *Repository) GetProjectDetectionStainResult(ctx context.Context, taskId uint64) (*commonpb.ProjectDetectionStainResult, error) {
	result := &commonpb.ProjectDetectionStainResult{}

	var (
		status            string
		hasStain          sql.NullBool
		stainCount        sql.NullInt64
		averageStainRatio sql.NullFloat64
		maxStainRatio     sql.NullFloat64
		regions           sql.NullString
		runtimeSeconds    sql.NullFloat64
		startedAt         sql.NullTime
		finishedAt        sql.NullTime
	)

	if err := r.mysql.QueryRowContext(ctx, `
		SELECT
			uuid,
			status,
			has_stain,
			stain_count,
			average_stain_ratio,
			max_stain_ratio,
			CAST(regions AS CHAR),
			runtime_seconds,
			started_at,
			finished_at
		FROM project_detection_stain_tasks
		WHERE id = ?
		LIMIT 1
	`, taskId).Scan(
		&result.TaskUuid,
		&status,
		&hasStain,
		&stainCount,
		&averageStainRatio,
		&maxStainRatio,
		&regions,
		&runtimeSeconds,
		&startedAt,
		&finishedAt,
	); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}

	result.Status = enum.ParseProjectDetectionSubTaskStatus(status)
	result.HasStain = utils.NullBool(hasStain)
	result.StainCount = utils.NullUint32(stainCount)
	result.AverageStainRatio = utils.NullFloat64(averageStainRatio)
	result.MaxStainRatio = utils.NullFloat64(maxStainRatio)
	if err := utils.UnmarshalRegions(regions, &result.Regions); err != nil {
		return nil, err
	}
	result.RuntimeSeconds = utils.NullFloat64(runtimeSeconds)
	result.StartedAt = utils.NullTimeString(startedAt)
	result.FinishedAt = utils.NullTimeString(finishedAt)

	return result, nil
}

// GetProjectDetectionFlatnessResult 按项目图像检测子任务 ID 查询玻璃平整度检测结果
func (r *Repository) GetProjectDetectionFlatnessResult(ctx context.Context, taskId uint64) (*commonpb.ProjectDetectionFlatnessResult, error) {
	result := &commonpb.ProjectDetectionFlatnessResult{}

	var (
		status         string
		flatnessResult sql.NullString
		unevenCount    sql.NullInt64
		regions        sql.NullString
		runtimeSeconds sql.NullFloat64
		startedAt      sql.NullTime
		finishedAt     sql.NullTime
	)

	if err := r.mysql.QueryRowContext(ctx, `
		SELECT
			uuid,
			status,
			result,
			uneven_count,
			CAST(regions AS CHAR),
			runtime_seconds,
			started_at,
			finished_at
		FROM project_detection_flatness_tasks
		WHERE id = ?
		LIMIT 1
	`, taskId).Scan(
		&result.TaskUuid,
		&status,
		&flatnessResult,
		&unevenCount,
		&regions,
		&runtimeSeconds,
		&startedAt,
		&finishedAt,
	); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}

	result.Status = enum.ParseProjectDetectionSubTaskStatus(status)
	result.Result = utils.NullString(flatnessResult)
	result.UnevenCount = utils.NullUint32(unevenCount)
	if err := utils.UnmarshalRegions(regions, &result.Regions); err != nil {
		return nil, err
	}
	result.RuntimeSeconds = utils.NullFloat64(runtimeSeconds)
	result.StartedAt = utils.NullTimeString(startedAt)
	result.FinishedAt = utils.NullTimeString(finishedAt)

	return result, nil
}

// GetProjectDetectionSpallingResult 按项目图像检测子任务 ID 查询玻璃爆裂检测结果
func (r *Repository) GetProjectDetectionSpallingResult(ctx context.Context, taskId uint64) (*commonpb.ProjectDetectionSpallingResult, error) {
	result := &commonpb.ProjectDetectionSpallingResult{}

	var (
		status         string
		hasSpalling    sql.NullBool
		confidence     sql.NullFloat64
		runtimeSeconds sql.NullFloat64
		startedAt      sql.NullTime
		finishedAt     sql.NullTime
	)

	if err := r.mysql.QueryRowContext(ctx, `
		SELECT
			uuid,
			status,
			has_spalling,
			confidence,
			runtime_seconds,
			started_at,
			finished_at
		FROM project_detection_spalling_tasks
		WHERE id = ?
		LIMIT 1
	`, taskId).Scan(
		&result.TaskUuid,
		&status,
		&hasSpalling,
		&confidence,
		&runtimeSeconds,
		&startedAt,
		&finishedAt,
	); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}

	result.Status = enum.ParseProjectDetectionSubTaskStatus(status)
	result.HasSpalling = utils.NullBool(hasSpalling)
	result.Confidence = utils.NullFloat64(confidence)
	result.RuntimeSeconds = utils.NullFloat64(runtimeSeconds)
	result.StartedAt = utils.NullTimeString(startedAt)
	result.FinishedAt = utils.NullTimeString(finishedAt)

	return result, nil
}

// StartProjectDetectionClassificationTask 执行分类任务时，按项目图像检测主任务 ID 更新开始时间
func (r *Repository) StartProjectDetectionClassificationTask(ctx context.Context, taskId uint64) (*model.ProjectDetectionTaskRecord, error) {
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

// StartProjectDetectionReasoningTask 执行推理任务时，按项目图像检测子任务代码和子任务 UUID 更新开始时间
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

// StartProjectDetectionSummaryTask 执行总结任务时，按项目图像检测主任务 UUID 更新开始时间
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

// UpdateProjectDetectionTaskStatus 按项目图像检测主任务 ID 更新任务状态
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

// UpdateProjectDetectionClassificationResult 按项目图像检测主任务 UUID 更新分类结果
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

	taskCodes, err = project_detection.NormalizeDetectionTaskCodes(taskCodes)
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

// UpdateProjectDetectionReasoningTaskResult 按项目图像检测子任务代码和子任务 UUID 更新推理结果
func (r *Repository) UpdateProjectDetectionReasoningTaskResult(ctx context.Context, taskCode, taskUuid string, status bizpb.ProjectDetectionSubTaskStatus_Value, resultJSON, artifactSha256Map string) (*model.ProjectDetectionTaskRecord, *model.ProjectDetectionSubTaskRecord, *model.ProjectDetectionSummaryTaskRecord, error) {
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

	if err := project_detection.UpdateProjectDetectionSubTaskResultTx(ctx, tx, taskCode, taskUuid, status, resultJSON, artifactSha256Map); err != nil {
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

// UpdateProjectDetectionSummaryResult 按项目图像检测主任务 UUID 更新总结结果
func (r *Repository) UpdateProjectDetectionSummaryResult(ctx context.Context, taskUuid string, status bizpb.ProjectDetectionSubTaskStatus_Value, resultJSON string) (*model.ProjectDetectionTaskRecord, *model.ProjectDetectionSummaryTaskRecord, error) {
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

	if err := project_detection.UpdateProjectDetectionSummaryTaskResultTx(ctx, tx, taskUuid, status, resultJSON); err != nil {
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

// GetProjectDetectionSubReportJSON 按项目图像检测子任务代码和子任务 ID 查询检测报告 JSON
func (r *Repository) GetProjectDetectionSubReportJSON(ctx context.Context, taskCode string, taskId uint64) (string, error) {
	switch enum.ParseDetectionTaskCode(taskCode) {
	case commonpb.DetectionTaskCode_Corrosion:
		return r.queryProjectDetectionSubReportJSON(ctx, taskId, `
			SELECT CAST(JSON_OBJECT(
				'has_corrosion', COALESCE(has_corrosion, false),
				'corrosion_count', COALESCE(corrosion_count, 0),
				'max_confidence', COALESCE(max_confidence, 0),
				'average_confidence', COALESCE(average_confidence, 0),
				'corrosion_pixels', COALESCE(corrosion_pixels, 0),
				'corrosion_ratio', COALESCE(corrosion_ratio, 0),
				'regions', COALESCE(regions, JSON_ARRAY()),
				'runtime_seconds', COALESCE(runtime_seconds, 0)
			) AS CHAR)
			FROM project_detection_corrosion_tasks
			WHERE id = ?
			LIMIT 1
		`)
	case commonpb.DetectionTaskCode_Crack:
		return r.queryProjectDetectionSubReportJSON(ctx, taskId, `
			SELECT CAST(JSON_OBJECT(
				'has_crack', COALESCE(has_crack, false),
				'crack_count', COALESCE(crack_count, 0),
				'crack_pixels', COALESCE(crack_pixels, 0),
				'crack_ratio', COALESCE(crack_ratio, 0),
				'regions', COALESCE(regions, JSON_ARRAY()),
				'runtime_seconds', COALESCE(runtime_seconds, 0)
			) AS CHAR)
			FROM project_detection_crack_tasks
			WHERE id = ?
			LIMIT 1
		`)
	case commonpb.DetectionTaskCode_Stain:
		return r.queryProjectDetectionSubReportJSON(ctx, taskId, `
			SELECT CAST(JSON_OBJECT(
				'has_stain', COALESCE(has_stain, false),
				'stain_count', COALESCE(stain_count, 0),
				'average_stain_ratio', COALESCE(average_stain_ratio, 0),
				'max_stain_ratio', COALESCE(max_stain_ratio, 0),
				'regions', COALESCE(regions, JSON_ARRAY()),
				'runtime_seconds', COALESCE(runtime_seconds, 0)
			) AS CHAR)
			FROM project_detection_stain_tasks
			WHERE id = ?
			LIMIT 1
		`)
	case commonpb.DetectionTaskCode_Flatness:
		return r.queryProjectDetectionSubReportJSON(ctx, taskId, `
			SELECT CAST(JSON_OBJECT(
				'result', COALESCE(result, ''),
				'uneven_count', COALESCE(uneven_count, 0),
				'regions', COALESCE(regions, JSON_ARRAY()),
				'runtime_seconds', COALESCE(runtime_seconds, 0)
			) AS CHAR)
			FROM project_detection_flatness_tasks
			WHERE id = ?
			LIMIT 1
		`)
	case commonpb.DetectionTaskCode_Spalling:
		return r.queryProjectDetectionSubReportJSON(ctx, taskId, `
			SELECT CAST(JSON_OBJECT(
				'has_spalling', COALESCE(has_spalling, false),
				'confidence', COALESCE(confidence, 0),
				'runtime_seconds', COALESCE(runtime_seconds, 0)
			) AS CHAR)
			FROM project_detection_spalling_tasks
			WHERE id = ?
			LIMIT 1
		`)
	default:
		return "", model.ErrUnsupportedDetectionTaskCode
	}
}

// queryProjectDetectionSubReportJSON 按项目图像检测子任务 ID 和查询语句查询检测报告 JSON
func (r *Repository) queryProjectDetectionSubReportJSON(ctx context.Context, taskId uint64, query string) (string, error) {
	reportJSON := ""
	if err := r.mysql.QueryRowContext(ctx, query, taskId).Scan(&reportJSON); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return "", nil
		}
		return "", err
	}
	return reportJSON, nil
}

// GetProjectDetectionSummaryTypedResult 按总结任务 ID 查询图像检测总结结果
func (r *Repository) GetProjectDetectionSummaryTypedResult(ctx context.Context, taskId uint64) (*commonpb.ProjectDetectionSummaryResult, error) {
	result := &commonpb.ProjectDetectionSummaryResult{}
	var (
		status     string
		startedAt  sql.NullTime
		finishedAt sql.NullTime
	)
	if err := r.mysql.QueryRowContext(ctx, `
		SELECT uuid, status, CAST(COALESCE(result_json, JSON_OBJECT()) AS CHAR), started_at, finished_at
		FROM project_detection_summary_tasks
		WHERE id = ?
		LIMIT 1
	`, taskId).Scan(&result.TaskUuid, &status, &result.ResultJson, &startedAt, &finishedAt); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	result.Status = enum.ParseProjectDetectionSubTaskStatus(status)
	result.StartedAt = utils.NullTimeString(startedAt)
	result.FinishedAt = utils.NullTimeString(finishedAt)
	return result, nil
}
