package reasoning

import (
	"context"
	"os"
	"path/filepath"
	"time"

	"google.golang.org/protobuf/proto"

	"icw_activity_reasoning/consts"
	"icw_activity_reasoning/internal/services/common"
	reasoningUtils "icw_activity_reasoning/internal/services/reasoning/utils"
	"icw_activity_reasoning/rpc/icw_core_biz"
	"icw_common/gen/activity/reasoning"
	"icw_common/gen/core/biz"
	"icw_common/utils"
)

// Start 启动原子检测任务
func (s *Service) Start(ctx context.Context, req *reasoningpb.StartRequest) (*reasoningpb.StartResponse, error) {
	resp := &reasoningpb.StartResponse{}
	err := s.CallRPC(ctx, req, func() error {
		if err := reasoningUtils.ValidateRequest(req, s.Registry()); err != nil {
			return err
		}
		taskReq := proto.Clone(req).(*reasoningpb.StartRequest)
		requestId := utils.RequestIdFromIncomingContext(ctx)
		go s.runModuleDetection(requestId, taskReq)
		return nil
	})
	return resp, err
}

// runModuleDetection 异步执行原子检测任务
func (s *Service) runModuleDetection(requestId string, req *reasoningpb.StartRequest) {
	s.Acquire()
	defer s.Release()

	ctx, cancel := context.WithTimeout(s.Ctx(), s.Config().ReasoningTaskTimeout)
	defer cancel()
	start := time.Now()

	callbackReq := &bizpb.ReportReasoningResultRequest{
		TaskUuid:  req.TaskUuid,
		ImageUuid: req.ImageUuid,
		TaskCode:  req.TaskCode,
	}

	artifactCount, err := s.executeModuleDetection(ctx, req, callbackReq)

	if utils.IsEmptyError(err) {
		callbackReq.Status = consts.DetectionStatusSucceeded
		callbackReq.ErrorMessage = ""
		common.ReasoningInfo(requestId, req.TaskCode, req.TaskUuid, req.ImageUuid, artifactCount, time.Since(start))
	} else {
		callbackReq.Status = consts.DetectionStatusFailed
		callbackReq.ErrorMessage = err.Error()
		common.ReasoningError(requestId, req.TaskCode, req.TaskUuid, req.ImageUuid, artifactCount, time.Since(start), err)
	}

	callbackCtx := utils.AppendRequestIdToOutgoingContext(context.Background(), requestId)
	cost := time.Since(start)
	callbackResp := &bizpb.ReportReasoningResultResponse{}
	if err := icw_core_biz.ReportReasoningResult(callbackCtx, s.CoreBizClient(), callbackReq, callbackResp); err != nil {
		common.CallbackError(requestId, req.TaskCode, req.TaskUuid, req.ImageUuid, callbackReq.Status, cost, err)
		return
	}
	common.CallbackInfo(requestId, req.TaskCode, req.TaskUuid, req.ImageUuid, callbackReq.Status, cost)
}

// executeModuleDetection 执行原子检测任务
func (s *Service) executeModuleDetection(ctx context.Context, req *reasoningpb.StartRequest, callbackReq *bizpb.ReportReasoningResultRequest) (int, error) {
	taskDir := filepath.Join(s.Config().ReasoningWorkDir, req.TaskCode, req.ImageUuid)
	if err := os.MkdirAll(taskDir, 0o755); err != nil {
		return 0, err
	}
	defer func() {
		_ = os.RemoveAll(taskDir)
	}()

	if err := reasoningUtils.DownloadOriginalImage(ctx, req, taskDir, s.Config().ArtifactDownloadTimeout); err != nil {
		return 0, err
	}

	detector, err := s.Registry().Get(req.TaskCode)
	if err != nil {
		return 0, err
	}
	if err := detector.Detect(ctx, req.ImageUuid); err != nil {
		return 0, err
	}
	reportJSON, err := reasoningUtils.ReadCompactReportJSON(taskDir)
	if err != nil {
		return 0, err
	}
	callbackReq.ResultJson = reportJSON
	if len(req.Artifacts) > 0 {
		uploadedArtifacts := reasoningUtils.UploadArtifacts(ctx, req.Artifacts, taskDir, s.Config().ArtifactUploadTimeout)
		callbackReq.Artifacts = uploadedArtifacts
	}
	callbackReq.Status = consts.DetectionStatusSucceeded
	return len(callbackReq.Artifacts), nil
}
