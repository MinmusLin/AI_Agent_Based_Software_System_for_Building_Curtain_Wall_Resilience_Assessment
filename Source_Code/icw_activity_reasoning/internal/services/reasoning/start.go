package reasoning

import (
	"context"
	"errors"
	"os"
	"path/filepath"
	"strings"
	"time"

	"google.golang.org/protobuf/proto"

	"icw_common/enum"
	"icw_common/gen/activity/reasoning"
	"icw_common/gen/core/biz"
	"icw_common/gen/core/common"
	"icw_common/rpc"
	"icw_common/utils"

	"icw_activity_reasoning/internal/services/common"
	reasoningUtils "icw_activity_reasoning/internal/services/reasoning/utils"
	"icw_activity_reasoning/rpc/icw_core_biz"
)

// Start 启动原子检测任务
func (s *Service) Start(ctx context.Context, req *reasoningpb.StartRequest) (*reasoningpb.StartResponse, error) {
	resp := &reasoningpb.StartResponse{}
	err := s.CallRPC(req, func() error {
		if err := reasoningUtils.ValidateRequest(req, s.Registry()); err != nil {
			return err
		}
		requestId := rpc.RequestIdFromIncomingContext(ctx)
		taskReq := proto.Clone(req).(*reasoningpb.StartRequest)
		go s.asyncExecuteDetection(requestId, taskReq)
		return nil
	})
	return resp, err
}

// asyncExecuteDetection 异步执行原子检测任务并回调 icw.core.biz
func (s *Service) asyncExecuteDetection(requestId string, req *reasoningpb.StartRequest) {
	s.Acquire()
	defer s.Release()

	ctx, cancel := context.WithTimeout(s.Ctx(), s.Config().ReasoningTaskTimeout)
	defer cancel()
	ctx = rpc.WithRequestIdToOutgoingContext(ctx, requestId)

	callbackReq := &bizpb.ReportReasoningResultRequest{
		TaskCode:  req.TaskCode,
		TaskUuid:  req.TaskUuid,
		ImageUuid: req.ImageUuid,
	}
	taskCode := enum.DetectionTaskCodeString(req.TaskCode)

	// 执行原子检测任务
	artifactCount, detectorCost, err := s.executeDetection(ctx, req, callbackReq)
	if utils.IsEmptyError(err) {
		callbackReq.Status = commonpb.TaskStatus_Succeeded
		callbackReq.ErrorMessage = ""
		common.ReasoningInfo(requestId, taskCode, req.TaskUuid, req.ImageUuid, artifactCount, detectorCost)
	} else {
		callbackReq.Status = commonpb.TaskStatus_Failed
		callbackReq.ErrorMessage = err.Error()
		common.ReasoningError(requestId, taskCode, req.TaskUuid, req.ImageUuid, artifactCount, detectorCost, err)
	}

	// 上报图像检测推理结果
	callbackCtx := rpc.WithRequestIdToOutgoingContext(context.Background(), requestId)
	callbackResp := &bizpb.ReportReasoningResultResponse{}
	callbackStart := time.Now()
	if err := icw_core_biz.ReportReasoningResult(callbackCtx, s.CoreBizClient(), callbackReq, callbackResp); utils.IsEmptyError(err) {
		common.CallbackInfo(requestId, taskCode, req.TaskUuid, req.ImageUuid, enum.TaskStatusString(callbackReq.Status), callbackStart)
		return
	}
	common.CallbackError(requestId, taskCode, req.TaskUuid, req.ImageUuid, enum.TaskStatusString(callbackReq.Status), callbackStart, err)
}

// executeDetection 执行原子检测任务
func (s *Service) executeDetection(ctx context.Context, req *reasoningpb.StartRequest, callbackReq *bizpb.ReportReasoningResultRequest) (int, time.Duration, error) {
	taskCode := enum.DetectionTaskCodeString(req.TaskCode)
	taskDir := filepath.Join(s.Config().ReasoningRuntimeDir, taskCode, req.ImageUuid)
	if err := os.RemoveAll(taskDir); err != nil {
		return 0, 0, err
	}
	if err := os.MkdirAll(taskDir, 0o755); err != nil {
		return 0, 0, err
	}
	defer func() {
		_ = os.RemoveAll(taskDir)
	}()

	originalURL, err := s.getProjectImageOriginalURL(ctx, req)
	if err != nil {
		return 0, 0, err
	}
	if err := reasoningUtils.DownloadOriginalImage(ctx, originalURL, taskDir, s.Config().ArtifactDownloadTimeout); err != nil {
		return 0, 0, err
	}

	detector, err := s.Registry().Get(taskCode)
	if err != nil {
		return 0, 0, err
	}

	detectorStart := time.Now()
	err = detector.Detect(ctx, req.ImageUuid)
	detectorCost := time.Since(detectorStart)
	if err != nil {
		return 0, time.Since(detectorStart), err
	}

	reportJSON, err := reasoningUtils.ReadCompactReportJSON(taskDir)
	if err != nil {
		return 0, detectorCost, err
	}

	callbackReq.ResultJson = reportJSON
	callbackReq.Artifacts = reasoningUtils.UploadArtifacts(ctx, req.ArtifactPolicy, taskDir, s.Config().ArtifactUploadTimeout)

	artifactCount := 0
	for _, artifact := range callbackReq.Artifacts {
		if artifact != nil && artifact.Uploaded {
			artifactCount++
		}
	}

	return artifactCount, detectorCost, nil
}

// getProjectImageOriginalURL 获取项目图像原图下载地址
func (s *Service) getProjectImageOriginalURL(ctx context.Context, req *reasoningpb.StartRequest) (string, error) {
	resp := &bizpb.GetProjectImageOriginalResponse{}
	err := icw_core_biz.GetProjectImageOriginal(ctx, s.CoreBizClient(), &bizpb.GetProjectImageOriginalRequest{
		UserId:    req.UserId,
		ProjectId: req.ProjectId,
		ImageUuid: req.ImageUuid,
	}, resp)
	if err != nil {
		return "", err
	}
	originalURL := strings.TrimSpace(resp.OriginalUrl)
	if originalURL == "" {
		return "", errors.New("project image original url is empty")
	}
	return originalURL, nil
}
