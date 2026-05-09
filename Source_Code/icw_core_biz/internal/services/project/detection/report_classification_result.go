package detection

import (
	"context"

	"icw_common/enum"
	"icw_common/gen/core/biz"
	"icw_common/gen/core/common"
	"icw_common/rpc/error"

	"icw_core_biz/internal/services/project/events"
)

// ReportClassificationResult 上报图像检测分类结果
func (s *Service) ReportClassificationResult(ctx context.Context, req *bizpb.ReportClassificationResultRequest) (*bizpb.ReportClassificationResultResponse, error) {
	resp := &bizpb.ReportClassificationResultResponse{}
	err := s.CallRPC(req, func() error {
		return s.reportClassificationResult(ctx, req)
	})
	return resp, err
}

func (s *Service) reportClassificationResult(ctx context.Context, req *bizpb.ReportClassificationResultRequest) error {
	var taskStatus bizpb.ProjectDetectionSubTaskStatus_Value
	switch req.Status {
	case commonpb.TaskStatus_Succeeded:
		taskStatus = bizpb.ProjectDetectionSubTaskStatus_Succeeded
	case commonpb.TaskStatus_Failed:
		taskStatus = bizpb.ProjectDetectionSubTaskStatus_Failed
	default:
		return rpc_error.BadRequestDefault("classification status is invalid")
	}
	taskCodes := make([]string, 0, len(req.TaskCodes))
	for _, taskCode := range req.TaskCodes {
		code := enum.DetectionTaskCodeString(taskCode)
		if code != "" {
			taskCodes = append(taskCodes, code)
		}
	}

	// 按主任务 UUID 更新项目图像检测分类结果
	task, subTasks, err := s.MySQL().UpdateProjectDetectionClassificationResult(ctx, req.TaskUuid, taskStatus, taskCodes)
	if err != nil || task == nil {
		return err
	}

	// 发布项目图像检测任务状态变化事件
	events.PublishProjectDetectionNodeStatusChangedEvent(
		ctx,
		s.RocketMQ(),
		task.UserId,
		task.ProjectId,
		task.ImageUuid,
		events.DetectionNodeCodeClassification,
		task.Uuid,
		task.Status,
		"",
		taskStatus,
	)

	if task.Status == bizpb.ProjectDetectionTaskStatus_Failed {
		return nil
	}

	// 启动项目图像检测推理子任务
	s.DetectionWorker().StartReasoningTasks(ctx, subTasks)

	return nil
}
