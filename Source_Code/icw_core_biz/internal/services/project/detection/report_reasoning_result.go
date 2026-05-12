package detection

import (
	"context"

	"icw_common/enum"
	"icw_common/gen/core/biz"
	"icw_common/gen/core/common"
	"icw_common/rpc/error"

	"icw_core_biz/internal/services/project/events"
	"icw_core_biz/internal/services/project/utils"
)

// ReportReasoningResult 上报图像检测推理结果
func (s *Service) ReportReasoningResult(ctx context.Context, req *bizpb.ReportReasoningResultRequest) (*bizpb.ReportReasoningResultResponse, error) {
	resp := &bizpb.ReportReasoningResultResponse{}
	err := s.CallRPC(req, func() error {
		return s.reportReasoningResult(ctx, req)
	})
	return resp, err
}

func (s *Service) reportReasoningResult(ctx context.Context, req *bizpb.ReportReasoningResultRequest) error {
	var taskStatus bizpb.ProjectDetectionSubTaskStatus_Value
	switch req.Status {
	case commonpb.TaskStatus_Succeeded:
		taskStatus = bizpb.ProjectDetectionSubTaskStatus_Succeeded
	case commonpb.TaskStatus_Failed:
		taskStatus = bizpb.ProjectDetectionSubTaskStatus_Failed
	default:
		return rpc_error.BadRequestDefault("reasoning status is invalid")
	}

	taskCode := enum.DetectionTaskCodeString(req.TaskCode)
	if taskCode == "" {
		return rpc_error.BadRequestDefault("detection task code is invalid")
	}

	artifactSha256Map := ""
	if req.Status == commonpb.TaskStatus_Succeeded {
		if err := utils.ValidateReasoningArtifactUploads(req.Artifacts); err != nil {
			return err
		}

		// 将图像检测推理产物上传结果转换为 Sha256 Map JSON
		var err error
		artifactSha256Map, err = utils.ArtifactSha256MapJSON(req.Artifacts)
		if err != nil {
			return err
		}
	}

	// 按推理任务 UUID 更新项目图像检测推理子任务结果
	task, subTask, summaryTask, err := s.MySQL().UpdateProjectDetectionReasoningTaskResult(ctx, taskCode, req.TaskUuid, taskStatus, req.ResultJson, artifactSha256Map)
	if err != nil || task == nil || subTask == nil {
		return err
	}

	// 发布项目图像检测任务状态变化事件
	events.PublishProjectDetectionNodeStatusChangedEvent(
		ctx,
		s.RocketMQ(),
		task.UserId,
		task.ProjectId,
		task.ImageUuid,
		events.ReasoningNodeCode(taskCode),
		task.Uuid,
		task.Status,
		subTask.Uuid,
		taskStatus,
	)

	if summaryTask != nil {
		// 启动项目图像检测总结任务
		s.DetectionWorker().StartDetectionSummaryTask(ctx, task, summaryTask)
	}

	return nil
}
