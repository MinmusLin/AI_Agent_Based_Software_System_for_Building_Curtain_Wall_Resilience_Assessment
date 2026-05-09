package summary

import (
	"context"
	"io"
	"net/http"
	"strconv"
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
	// ProjectSummaryType 项目总结任务
	ProjectSummaryType = "project"
)

// StartProjectSummary 启动项目总结任务
func (s *Service) StartProjectSummary(ctx context.Context, req *summarypb.StartProjectSummaryRequest) (*summarypb.StartProjectSummaryResponse, error) {
	resp := &summarypb.StartProjectSummaryResponse{}
	err := s.CallRPC(req, func() error {
		if err := summaryUtils.ValidateStartProjectSummaryRequest(req); err != nil {
			return err
		}
		requestId := rpc.RequestIdFromIncomingContext(ctx)
		projectReq := proto.Clone(req).(*summarypb.StartProjectSummaryRequest)
		go s.asyncExecuteProjectSummary(requestId, projectReq)
		return nil
	})
	return resp, err
}

// asyncExecuteProjectSummary 异步执行项目总结任务并回调 icw.core.biz
func (s *Service) asyncExecuteProjectSummary(requestId string, req *summarypb.StartProjectSummaryRequest) {
	s.AcquireProjectSummary()
	defer s.ReleaseProjectSummary()

	projectId, _ := strconv.ParseUint(req.ProjectId, 10, 64)
	callbackReq := &bizpb.ReportProjectSummaryResultRequest{
		ProjectId: projectId,
	}
	ctx, cancel := context.WithTimeout(s.Ctx(), time.Duration(s.Config().AgentRequestTimeoutSeconds)*time.Second)
	defer cancel()

	start := time.Now()
	result, err := executeProjectSummary(ctx, s.ProjectSummaryAgentClient(), req)
	cost := time.Since(start)

	if utils.IsEmptyError(err) {
		callbackReq.Status = commonpb.TaskStatus_Succeeded
		callbackReq.ResultJson = result
		callbackReq.ErrorMessage = ""
		common.SummaryInfo(requestId, ProjectSummaryType, req.ProjectId, cost, result)
	} else {
		callbackReq.Status = commonpb.TaskStatus_Failed
		callbackReq.ResultJson = ""
		callbackReq.ErrorMessage = err.Error()
		common.SummaryError(requestId, ProjectSummaryType, req.ProjectId, cost, result, err)
	}

	// 上报项目总结结果
	callbackCtx := rpc.WithRequestIdToOutgoingContext(context.Background(), requestId)
	callbackResp := &bizpb.ReportProjectSummaryResultResponse{}
	callbackStart := time.Now()
	if err := icw_core_biz.ReportProjectSummaryResult(callbackCtx, s.CoreBizClient(), callbackReq, callbackResp); utils.IsEmptyError(err) {
		common.CallbackInfo(requestId, ProjectSummaryType, req.ProjectId, enum.TaskStatusString(callbackReq.Status), callbackStart)
		return
	}
	common.CallbackError(requestId, ProjectSummaryType, req.ProjectId, enum.TaskStatusString(callbackReq.Status), callbackStart, err)
}

// executeProjectSummary 执行项目总结任务
func executeProjectSummary(ctx context.Context, client *agent.Client, req *summarypb.StartProjectSummaryRequest) (string, error) {
	sourceJSON, err := downloadSummarySource(ctx, req.SourceUrl)
	if err != nil {
		return "", err
	}
	return client.Chat(ctx, agent.Message{
		Text: sourceJSON,
	})
}

// downloadSummarySource 下载项目总结原始数据 JSON
func downloadSummarySource(ctx context.Context, sourceURL string) (string, error) {
	request, err := http.NewRequestWithContext(ctx, http.MethodGet, sourceURL, nil)
	if err != nil {
		return "", err
	}
	response, err := http.DefaultClient.Do(request)
	if err != nil {
		return "", err
	}
	defer func() {
		_ = response.Body.Close()
	}()
	if response.StatusCode < http.StatusOK || response.StatusCode >= http.StatusMultipleChoices {
		return "", io.ErrUnexpectedEOF
	}
	body, err := io.ReadAll(response.Body)
	if err != nil {
		return "", err
	}
	return string(body), nil
}
