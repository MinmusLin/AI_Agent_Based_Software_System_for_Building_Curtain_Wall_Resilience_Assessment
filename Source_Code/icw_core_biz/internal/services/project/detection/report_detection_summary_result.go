package detection

import (
	"context"
	"strings"

	"icw_common/enum"
	"icw_common/gen/activity"
	"icw_common/gen/core/biz"
	"icw_common/rpc/error"

	"icw_core_biz/internal/services/project/events"
)

// ReportDetectionSummaryResult 上报图像检测总结结果
func (s *Service) ReportDetectionSummaryResult(ctx context.Context, req *bizpb.ReportDetectionSummaryResultRequest) (*bizpb.ReportDetectionSummaryResultResponse, error) {
	resp := &bizpb.ReportDetectionSummaryResultResponse{}
	err := s.CallRPC(ctx, req, func() error {
		return s.reportDetectionSummaryResult(ctx, req)
	})
	return resp, err
}

func (s *Service) reportDetectionSummaryResult(ctx context.Context, req *bizpb.ReportDetectionSummaryResultRequest) error {
	req.TaskUuid = strings.TrimSpace(req.TaskUuid)
	req.ImageUuid = strings.TrimSpace(req.ImageUuid)
	if req.TaskUuid == "" {
		return rpc_error.BadRequestDefault("task uuid is required")
	}
	if req.ImageUuid == "" {
		return rpc_error.BadRequestDefault("image uuid is required")
	}
	status := enum.ParseDetectionStatus(req.Status)
	var taskStatus bizpb.ProjectDetectionSubTaskStatus_Value
	switch status {
	case activitypb.DetectionStatus_Succeeded:
		taskStatus = bizpb.ProjectDetectionSubTaskStatus_Succeeded
	case activitypb.DetectionStatus_Failed:
		taskStatus = bizpb.ProjectDetectionSubTaskStatus_Failed
	default:
		return rpc_error.BadRequestDefault("detection status is invalid")
	}
	task, summaryTask, err := s.MySQL().UpdateProjectDetectionSummaryResult(ctx, req.TaskUuid, taskStatus, req.ResultJson)
	if err != nil || task == nil || summaryTask == nil {
		return err
	}
	events.PublishProjectDetectionNodeStatusChangedEvent(
		ctx,
		s.RocketMQ(),
		task.UserId,
		task.ProjectId,
		task.ImageUuid,
		events.DetectionNodeCodeSummary,
		task.Uuid,
		enum.ProjectDetectionTaskStatusString(task.Status),
		summaryTask.Uuid,
		enum.ProjectDetectionSubTaskStatusString(taskStatus),
	)
	return nil
}
