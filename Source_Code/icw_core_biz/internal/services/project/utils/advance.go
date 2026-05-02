package utils

import (
	"context"
	"errors"

	"icw_core_biz/internal/services/project/consts"
	"icw_core_biz/repositories/mysql"
)

// PreAdvanceProject 项目进度推进前扩展点
func PreAdvanceProject(ctx context.Context, repo *mysql.Repository, userId, projectId uint64, fromProgress, toProgress consts.ProjectProgress) error {
	// 项目基础信息阶段 -> 图像资产构建阶段
	if fromProgress == consts.ProjectProgressInitializationFinished && toProgress == consts.ProjectProgressProfileFinished {
		return nil
	}
	// 图像资产构建阶段 -> Agent 智能检测阶段
	if fromProgress == consts.ProjectProgressProfileFinished && toProgress == consts.ProjectProgressAssetsFinished {
		return nil
	}
	// Agent 智能检测阶段 -> 人工复核确认阶段
	if fromProgress == consts.ProjectProgressAssetsFinished && toProgress == consts.ProjectProgressDetectFinished {
		return nil
	}
	// 人工复核确认阶段 -> 评估报告生成阶段
	if fromProgress == consts.ProjectProgressDetectFinished && toProgress == consts.ProjectProgressReviewFinished {
		return nil
	}
	// 评估报告生成阶段 -> 项目已完成
	if fromProgress == consts.ProjectProgressReviewFinished && toProgress == consts.ProjectProgressReportFinished {
		return nil
	}
	return errors.New("invalid from progress and to progress")
}

// PostAdvanceProject 项目进度推进后扩展点
func PostAdvanceProject(ctx context.Context, repo *mysql.Repository, userId, projectId uint64, fromProgress, toProgress consts.ProjectProgress) error {
	// 项目基础信息阶段 -> 图像资产构建阶段
	if fromProgress == consts.ProjectProgressInitializationFinished && toProgress == consts.ProjectProgressProfileFinished {
		return nil
	}
	// 图像资产构建阶段 -> Agent 智能检测阶段
	if fromProgress == consts.ProjectProgressProfileFinished && toProgress == consts.ProjectProgressAssetsFinished {
		return nil
	}
	// Agent 智能检测阶段 -> 人工复核确认阶段
	if fromProgress == consts.ProjectProgressAssetsFinished && toProgress == consts.ProjectProgressDetectFinished {
		return nil
	}
	// 人工复核确认阶段 -> 评估报告生成阶段
	if fromProgress == consts.ProjectProgressDetectFinished && toProgress == consts.ProjectProgressReviewFinished {
		return nil
	}
	// 评估报告生成阶段 -> 项目已完成
	if fromProgress == consts.ProjectProgressReviewFinished && toProgress == consts.ProjectProgressReportFinished {
		return nil
	}
	return errors.New("invalid from progress and to progress")
}
