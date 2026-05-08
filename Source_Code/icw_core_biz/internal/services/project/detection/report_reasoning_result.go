package detection

import (
	"context"
	"encoding/json"
	"strings"

	"icw_common/enum"
	"icw_common/gen/activity"
	"icw_common/gen/core/biz"
	"icw_common/rpc/error"

	"icw_core_biz/internal/services/project/events"
)

// ReportReasoningResult 上报图像检测推理结果
func (s *Service) ReportReasoningResult(ctx context.Context, req *bizpb.ReportReasoningResultRequest) (*bizpb.ReportReasoningResultResponse, error) {
	resp := &bizpb.ReportReasoningResultResponse{}
	err := s.CallRPC(ctx, req, func() error {
		return s.reportReasoningResult(ctx, req)
	})
	return resp, err
}

func (s *Service) reportReasoningResult(ctx context.Context, req *bizpb.ReportReasoningResultRequest) error {
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
	artifactSha256Map, err := artifactSha256MapJSON(req.Artifacts)
	if err != nil {
		return err
	}
	task, subTask, summaryTask, err := s.MySQL().UpdateProjectDetectionReasoningTaskResult(ctx, req.TaskCode, req.TaskUuid, taskStatus, req.ResultJson, artifactSha256Map)
	if err != nil || task == nil || subTask == nil {
		return err
	}
	events.PublishProjectDetectionNodeStatusChangedEvent(
		ctx,
		s.RocketMQ(),
		task.UserId,
		task.ProjectId,
		task.ImageUuid,
		events.ReasoningNodeCode(req.TaskCode),
		task.Uuid,
		enum.ProjectDetectionTaskStatusString(task.Status),
		subTask.Uuid,
		enum.ProjectDetectionSubTaskStatusString(taskStatus),
	)
	if summaryTask != nil {
		s.DetectionWorker().StartDetectionSummaryTask(ctx, task, summaryTask)
	}
	return nil
}

func artifactSha256MapJSON(artifacts []*bizpb.ReasoningArtifactUploadResult) (string, error) {
	artifactSha256Map := make(map[string]string)
	for _, artifact := range artifacts {
		if artifact == nil || !artifact.Uploaded {
			continue
		}
		name := strings.TrimSpace(artifact.Name)
		sha256 := strings.TrimSpace(artifact.Sha256)
		if name == "" || sha256 == "" {
			continue
		}
		artifactSha256Map[name] = sha256
	}
	bytes, err := json.Marshal(artifactSha256Map)
	if err != nil {
		return "", err
	}
	return string(bytes), nil
}
