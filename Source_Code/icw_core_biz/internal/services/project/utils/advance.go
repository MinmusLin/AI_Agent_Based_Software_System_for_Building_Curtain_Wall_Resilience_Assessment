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
		return repo.AdvanceProjectProfileToAssets(ctx, userId, projectId, nextStatus)
	}

	// 图像资产构建阶段 -> Agent 智能检测阶段
	if fromProgress == dto.ProjectProgressProfileFinished && toProgress == dto.ProjectProgressAssetsFinished {
		return repo.AdvanceProjectAssetsToDetection(ctx, userId, projectId, nextStatus)
	}

	// Agent 智能检测阶段 -> 人工复核确认阶段
	if fromProgress == dto.ProjectProgressAssetsFinished && toProgress == dto.ProjectProgressDetectionFinished {
		return repo.AdvanceProjectDetectionToReview(ctx, userId, projectId, nextStatus)
	}

	// 人工复核确认阶段 -> 评估报告生成阶段
	if fromProgress == dto.ProjectProgressDetectionFinished && toProgress == dto.ProjectProgressReviewFinished {
		return repo.AdvanceProjectReviewToReport(ctx, userId, projectId, nextStatus)
	}

	// 评估报告生成阶段 -> 项目已完成
	if fromProgress == dto.ProjectProgressReviewFinished && toProgress == dto.ProjectProgressReportFinished {
		return repo.AdvanceProjectReportToFinished(ctx, userId, projectId, nextStatus)
	}

	return false, rpc_err.BadRequestDefault("invalid from progress and to progress")
}
