package summary

import (
	"context"
	"time"

	"google.golang.org/protobuf/proto"

	"icw_common/enum"
	"icw_common/gen/activity/summary"
	"icw_common/gen/core/biz"
	"icw_common/gen/core/common"
	"icw_common/rpc"
	"icw_common/utils"

	"icw_activity_summary/internal/agent"
	"icw_activity_summary/internal/services/common"
	summaryUtils "icw_activity_summary/internal/services/summary/utils"
	"icw_activity_summary/rpc/icw_core_biz"
)

const (
	// DetectionSummaryType 图像检测总结任务
	DetectionSummaryType = "detection"
)

// StartDetectionSummary 启动图像检测总结任务
func (s *Service) StartDetectionSummary(ctx context.Context, req *summarypb.StartDetectionSummaryRequest) (*summarypb.StartDetectionSummaryResponse, error) {
	resp := &summarypb.StartDetectionSummaryResponse{}
	err := s.CallRPC(req, func() error {
		if err := summaryUtils.ValidateStartDetectionSummaryRequest(req); err != nil {
			return err
		}
		requestId := rpc.RequestIdFromIncomingContext(ctx)
		taskReq := proto.Clone(req).(*summarypb.StartDetectionSummaryRequest)
		go s.asyncExecuteDetectionSummary(requestId, taskReq)
		return nil
	})
	return resp, err
}

// asyncExecuteDetectionSummary 异步执行图像检测总结任务并回调 icw.core.biz
func (s *Service) asyncExecuteDetectionSummary(requestId string, req *summarypb.StartDetectionSummaryRequest) {
	s.AcquireDetectionSummary()
	defer s.ReleaseDetectionSummary()

	callbackReq := &bizpb.ReportDetectionSummaryResultRequest{
		TaskUuid:  req.TaskUuid,
		ImageUuid: req.ImageUuid,
	}
	ctx, cancel := context.WithTimeout(s.Ctx(), time.Duration(s.Config().AgentRequestTimeoutSeconds)*time.Second)
	defer cancel()

	start := time.Now()
	result, err := executeDetectionSummary(ctx, s.DetectionSummaryAgentClient(), req)
	cost := time.Since(start)

	if utils.IsEmptyError(err) {
		callbackReq.Status = commonpb.TaskStatus_Succeeded
		callbackReq.Result = result
		callbackReq.ErrorMessage = ""
		common.SummaryInfo(requestId, DetectionSummaryType, req.TaskUuid, cost, result)
	} else {
		callbackReq.Status = commonpb.TaskStatus_Failed
		callbackReq.Result = ""
		callbackReq.ErrorMessage = err.Error()
		common.SummaryError(requestId, DetectionSummaryType, req.TaskUuid, cost, result, err)
	}

	// 上报图像检测总结结果
	callbackCtx := rpc.WithRequestIdToOutgoingContext(context.Background(), requestId)
	callbackResp := &bizpb.ReportDetectionSummaryResultResponse{}
	callbackStart := time.Now()
	if err := icw_core_biz.ReportDetectionSummaryResult(callbackCtx, s.CoreBizClient(), callbackReq, callbackResp); utils.IsEmptyError(err) {
		common.CallbackInfo(requestId, DetectionSummaryType, req.TaskUuid, enum.TaskStatusString(callbackReq.Status), callbackStart)
		return
	}
	common.CallbackError(requestId, DetectionSummaryType, req.TaskUuid, enum.TaskStatusString(callbackReq.Status), callbackStart, err)
}

// executeDetectionSummary 执行图像检测总结任务
func executeDetectionSummary(ctx context.Context, client *agent.Client, req *summarypb.StartDetectionSummaryRequest) (string, error) {
	return client.Chat(ctx, req.SourceJson)
}
