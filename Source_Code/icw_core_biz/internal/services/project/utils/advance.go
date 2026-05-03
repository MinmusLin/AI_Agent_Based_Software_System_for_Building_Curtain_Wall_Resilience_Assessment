package utils

import (
	"context"

	"icw_core_biz/pkg/dto"
	"icw_core_biz/pkg/rpc_err"
	"icw_core_biz/repositories/mysql"
)

// ProjectAlreadyAdvanced 判断项目进度是否被并发请求流转
func ProjectAlreadyAdvanced(ctx context.Context, repo *mysql.Repository, userId, projectId uint64, toProgress dto.ProjectProgress, nextStatus dto.ProjectStatus) (bool, error) {
	projectRecord, err := repo.FindProjectByIdAndUserId(ctx, userId, projectId)
	if err != nil || projectRecord == nil {
		return false, err
	}
	return projectRecord.Progress == toProgress && projectRecord.Status == nextStatus, nil
}

// BeforeAdvanceProject 项目进度流转前置扩展点
func BeforeAdvanceProject(ctx context.Context, repo *mysql.Repository, userId, projectId uint64, fromProgress, toProgress dto.ProjectProgress) error {
	// 项目基础信息阶段 -> 图像资产构建阶段
	if fromProgress == dto.ProjectProgressInitializationFinished && toProgress == dto.ProjectProgressProfileFinished {
		return nil
	}

	// 图像资产构建阶段 -> Agent 智能检测阶段
	if fromProgress == dto.ProjectProgressProfileFinished && toProgress == dto.ProjectProgressAssetsFinished {
		groupCount, err := repo.CountProjectGroups(ctx, userId, projectId)
		if err != nil {
			return err
		}
		if groupCount == 0 {
			return rpc_err.BadRequest(rpc_err.DetailProjectGroupCannotDeleteLast, "project must keep at least one group")
		}

		stats, err := repo.GetProjectAssetsReadyStats(ctx, userId, projectId)
		if err != nil {
			return err
		}
		if stats == nil {
			return rpc_err.BadRequest(rpc_err.DetailProjectNotAccessible, "project group is not accessible")
		}
		if stats.UploadedImageCount == 0 {
			return rpc_err.BadRequestDefault("project must keep at least one uploaded image")
		}
		if stats.EmptyGroupCount > 0 {
			return rpc_err.BadRequestDefault("project must not have empty groups")
		}
		if stats.PendingImageCount > 0 || stats.FailedImageCount > 0 {
			return rpc_err.BadRequestDefault("project must not have pending or failed images")
		}
		return nil
	}

	// Agent 智能检测阶段 -> 人工复核确认阶段
	if fromProgress == dto.ProjectProgressAssetsFinished && toProgress == dto.ProjectProgressDetectionFinished {
		return nil
	}

	// 人工复核确认阶段 -> 评估报告生成阶段
	if fromProgress == dto.ProjectProgressDetectionFinished && toProgress == dto.ProjectProgressReviewFinished {
		return nil
	}

	// 评估报告生成阶段 -> 项目已完成
	if fromProgress == dto.ProjectProgressReviewFinished && toProgress == dto.ProjectProgressReportFinished {
		return nil
	}

	return rpc_err.BadRequestDefault("invalid from progress and to progress")
}

// AdvanceProject 项目进度流转扩展点
func AdvanceProject(ctx context.Context, repo *mysql.Repository, userId, projectId uint64, fromProgress, toProgress dto.ProjectProgress, nextStatus dto.ProjectStatus) (bool, error) {
	// 项目基础信息阶段 -> 图像资产构建阶段
	if fromProgress == dto.ProjectProgressInitializationFinished && toProgress == dto.ProjectProgressProfileFinished {
		return repo.AdvanceProject(ctx, userId, projectId, fromProgress, toProgress, nextStatus, mysql.PostAdvanceProjectProfileToAssets)
	}

	// 图像资产构建阶段 -> Agent 智能检测阶段
	if fromProgress == dto.ProjectProgressProfileFinished && toProgress == dto.ProjectProgressAssetsFinished {
		return repo.AdvanceProject(ctx, userId, projectId, fromProgress, toProgress, nextStatus, nil)
	}

	// Agent 智能检测阶段 -> 人工复核确认阶段
	if fromProgress == dto.ProjectProgressAssetsFinished && toProgress == dto.ProjectProgressDetectionFinished {
		return repo.AdvanceProject(ctx, userId, projectId, fromProgress, toProgress, nextStatus, nil)
	}

	// 人工复核确认阶段 -> 评估报告生成阶段
	if fromProgress == dto.ProjectProgressDetectionFinished && toProgress == dto.ProjectProgressReviewFinished {
		return repo.AdvanceProject(ctx, userId, projectId, fromProgress, toProgress, nextStatus, nil)
	}

	// 评估报告生成阶段 -> 项目已完成
	if fromProgress == dto.ProjectProgressReviewFinished && toProgress == dto.ProjectProgressReportFinished {
		return repo.AdvanceProject(ctx, userId, projectId, fromProgress, toProgress, nextStatus, nil)
	}

	return false, rpc_err.BadRequestDefault("invalid from progress and to progress")
}
