package model

import (
	"errors"
)

var (
	// ErrRefreshTokenNotReplaceable 旧 Refresh Token 不存在或已吊销，不能继续完成轮换
	ErrRefreshTokenNotReplaceable = errors.New("refresh token not replaceable")
	// ErrProjectGroupCannotDeleteLast 项目应至少存在一个图像组
	ErrProjectGroupCannotDeleteLast = errors.New("project must keep at least one group")
	// ErrUnsupportedDetectionTaskCode 不支持的原子检测能力代码
	ErrUnsupportedDetectionTaskCode = errors.New("unsupported detection task code")
	// ErrMoveProjectGroupDirectionConflict 移动图像组时不能同时置顶和置底
	ErrMoveProjectGroupDirectionConflict = errors.New("move to first and move to last cannot both be true")
	// ErrProjectDetectionTaskStatusInvalid 项目图像检测主任务状态非法
	ErrProjectDetectionTaskStatusInvalid = errors.New("project detection task status is invalid")
	// ErrProjectDetectionSubTaskStatusInvalid 项目图像检测子任务状态非法
	ErrProjectDetectionSubTaskStatusInvalid = errors.New("project detection sub task status is invalid")
	// ErrProjectDetectionSummaryTaskStatusInvalid 项目图像检测总结任务状态非法
	ErrProjectDetectionSummaryTaskStatusInvalid = errors.New("project detection summary task status is invalid")
	// ErrProjectDetectionReviewVerdictInvalid 项目图像检测人工复核结论非法
	ErrProjectDetectionReviewVerdictInvalid = errors.New("project detection review verdict is invalid")
	// ErrProjectReportStatusInvalid 项目评估报告状态非法
	ErrProjectReportStatusInvalid = errors.New("project report status is invalid")
)
