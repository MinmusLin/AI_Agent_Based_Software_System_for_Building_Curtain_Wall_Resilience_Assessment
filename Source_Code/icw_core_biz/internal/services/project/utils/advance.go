package utils

import (
	"context"
	"database/sql"

	"icw_common/gen/core/biz"
	"icw_common/rpc/error"

	"icw_core_biz/internal/services/project/consts"
	"icw_core_biz/repositories/mysql"
	"icw_core_biz/repositories/mysql/utils"
)

// ProjectAlreadyAdvanced 判断项目进度是否被并发请求流转
func ProjectAlreadyAdvanced(ctx context.Context, repo *mysql.Repository, userId, projectId uint64, toProgress bizpb.ProjectProgress_Value, nextStatus bizpb.ProjectStatus_Value) (bool, error) {
	projectRecord, err := repo.FindProjectByIdAndUserId(ctx, userId, projectId)
	if err != nil || projectRecord == nil {
		return false, err
	}
	return projectRecord.Progress == toProgress && projectRecord.Status == nextStatus, nil
}

// BeforeAdvanceProject 项目进度流转前置扩展点
func BeforeAdvanceProject(ctx context.Context, repo *mysql.Repository, userId, projectId uint64, fromProgress, toProgress bizpb.ProjectProgress_Value) error {
	// 项目基础信息阶段 -> 图像资产构建阶段
	if fromProgress == bizpb.ProjectProgress_InitializationFinished && toProgress == bizpb.ProjectProgress_ProfileFinished {
		return nil
	}

	// 图像资产构建阶段 -> Agent 智能检测阶段
	if fromProgress == bizpb.ProjectProgress_ProfileFinished && toProgress == bizpb.ProjectProgress_AssetsFinished {
		groupCount, err := repo.CountProjectGroups(ctx, userId, projectId)
		if err != nil {
			return err
		}
		if groupCount == 0 {
			return rpc_error.BadRequest(rpc_error.DetailProjectAtLeastOneGroupRequired, "project must keep at least one group")
		}

		stats, err := repo.GetProjectAssetsReadyStats(ctx, userId, projectId)
		if err != nil {
			return err
		}
		if stats == nil {
			return rpc_error.BadRequest(rpc_error.DetailProjectNotAccessible, "project group is not accessible")
		}
		if stats.UploadedImageCount == 0 {
			return rpc_error.BadRequest(rpc_error.DetailProjectUploadedImageCountRequired, "project must keep at least one uploaded image")
		}
		if stats.EmptyGroupCount > 0 {
			return rpc_error.BadRequest(rpc_error.DetailProjectEmptyGroupCountInvalid, "project must not have empty groups")
		}
		if stats.PendingImageCount > 0 || stats.FailedImageCount > 0 {
			return rpc_error.BadRequest(rpc_error.DetailProjectPendingOrFailedImageCountInvalid, "project must not have pending or failed images")
		}
		return nil
	}

	// Agent 智能检测阶段 -> 人工复核确认阶段
	if fromProgress == bizpb.ProjectProgress_AssetsFinished && toProgress == bizpb.ProjectProgress_DetectionFinished {
		allSucceeded, err := repo.ProjectDetectionTasksAllSucceeded(ctx, userId, projectId)
		if err != nil {
			return err
		}
		if !allSucceeded {
			return rpc_error.BadRequestDefault("project detection tasks must all succeed")
		}
		return nil
	}

	// 人工复核确认阶段 -> 评估报告生成阶段
	if fromProgress == bizpb.ProjectProgress_DetectionFinished && toProgress == bizpb.ProjectProgress_ReviewFinished {
		return nil
	}

	// 评估报告生成阶段 -> 项目已完成
	if fromProgress == bizpb.ProjectProgress_ReviewFinished && toProgress == bizpb.ProjectProgress_ReportFinished {
		return nil
	}

	return rpc_error.BadRequestDefault("invalid from progress and to progress")
}

// AdvanceProject 项目进度流转扩展点
func AdvanceProject(ctx context.Context, repo *mysql.Repository, userId, projectId uint64, fromProgress, toProgress bizpb.ProjectProgress_Value, nextStatus bizpb.ProjectStatus_Value) (bool, error) {
	// 项目基础信息阶段 -> 图像资产构建阶段
	if fromProgress == bizpb.ProjectProgress_InitializationFinished && toProgress == bizpb.ProjectProgress_ProfileFinished {
		return repo.AdvanceProject(ctx, userId, projectId, fromProgress, toProgress, nextStatus, postAdvanceProjectProfileToAssets)
	}

	// 图像资产构建阶段 -> Agent 智能检测阶段
	if fromProgress == bizpb.ProjectProgress_ProfileFinished && toProgress == bizpb.ProjectProgress_AssetsFinished {
		return repo.AdvanceProject(ctx, userId, projectId, fromProgress, toProgress, nextStatus, nil)
	}

	// Agent 智能检测阶段 -> 人工复核确认阶段
	if fromProgress == bizpb.ProjectProgress_AssetsFinished && toProgress == bizpb.ProjectProgress_DetectionFinished {
		return repo.AdvanceProject(ctx, userId, projectId, fromProgress, toProgress, nextStatus, nil)
	}

	// 人工复核确认阶段 -> 评估报告生成阶段
	if fromProgress == bizpb.ProjectProgress_DetectionFinished && toProgress == bizpb.ProjectProgress_ReviewFinished {
		return repo.AdvanceProject(ctx, userId, projectId, fromProgress, toProgress, nextStatus, nil)
	}

	// 评估报告生成阶段 -> 项目已完成
	if fromProgress == bizpb.ProjectProgress_ReviewFinished && toProgress == bizpb.ProjectProgress_ReportFinished {
		return repo.AdvanceProject(ctx, userId, projectId, fromProgress, toProgress, nextStatus, nil)
	}

	return false, rpc_error.BadRequestDefault("invalid from progress and to progress")
}

// postAdvanceProjectProfileToAssets 项目进度流转后置扩展点：项目基础信息阶段 -> 图像资产构建阶段
func postAdvanceProjectProfileToAssets(ctx context.Context, tx *sql.Tx, userId, projectId uint64) error {
	_, err := tx.ExecContext(ctx, `
		INSERT INTO project_groups(project_id, user_id, name, sort_order)
		SELECT ?, ?, ?, COALESCE(MAX(sort_order), -1) + 1
		FROM project_groups
		WHERE user_id = ? AND project_id = ?
	`, projectId, userId, consts.DefaultProjectGroupName, userId, projectId)

	if utils.IsDuplicateEntryError(err) {
		return nil
	}

	return err
}
