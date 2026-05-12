package summary

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/google/uuid"
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
	// ProjectSummarySourceFileName 项目总结原始数据附件名称
	ProjectSummarySourceFileName = "source.txt"
	// ProjectSummarySourceContentType 项目总结原始数据附件类型
	ProjectSummarySourceContentType = "text/plain; charset=utf-8"
)

// projectSummarySourceMeta 项目总结源数据基础元信息
type projectSummarySourceMeta struct {
	GroupCount int `json:"group_count"`
	ImageCount int `json:"image_count"`
	Project    struct {
		ProjectName  string `json:"project_name"`
		BuildingName string `json:"building_name"`
	} `json:"project"`
}

// StartProjectSummary 启动项目总结任务
func (s *Service) StartProjectSummary(ctx context.Context, req *summarypb.StartProjectSummaryRequest) (*summarypb.StartProjectSummaryResponse, error) {
	resp := &summarypb.StartProjectSummaryResponse{}
	err := s.CallRPC(req, func() error {
		if err := summaryUtils.ValidateStartProjectSummaryRequest(req); err != nil {
			return err
		}
		requestId := rpc.RequestIdFromIncomingContext(ctx)
		if requestId == "" {
			requestId = uuid.NewString()
		}
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
	ctx, cancel := context.WithTimeout(s.Ctx(), s.Config().AgentRequestTimeout)
	defer cancel()

	// 执行项目总结任务
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
	sourceFile, err := downloadSummarySourceFile(ctx, req.SourceUrl)
	if err != nil {
		return "", err
	}
	sourceMeta, err := parseProjectSummarySourceMeta(sourceFile)
	if err != nil {
		return "", err
	}
	output, err := client.ChatWithFile(ctx, fmt.Sprintf(
		"附件已上传，请根据指令进行符合要求的输出。输出 JSON 顶层 group_count 必须为 %d，image_count 必须为 %d。",
		sourceMeta.GroupCount,
		sourceMeta.ImageCount,
	), agent.File{
		Name:        ProjectSummarySourceFileName,
		Data:        sourceFile,
		ContentType: ProjectSummarySourceContentType,
	})
	if err != nil {
		return "", err
	}
	compacted, err := utils.CompactAgentJSONObjectString(output)
	if err != nil {
		return "", err
	}
	if err := validateProjectSummaryOutput(compacted, sourceMeta); err != nil {
		return "", err
	}
	return compacted, nil
}

// downloadSummarySourceFile 下载项目总结原始数据文本附件
func downloadSummarySourceFile(ctx context.Context, sourceURL string) ([]byte, error) {
	request, err := http.NewRequestWithContext(ctx, http.MethodGet, sourceURL, nil)
	if err != nil {
		return nil, err
	}
	response, err := http.DefaultClient.Do(request)
	if err != nil {
		return nil, err
	}
	defer func() {
		_ = response.Body.Close()
	}()
	if response.StatusCode < http.StatusOK || response.StatusCode >= http.StatusMultipleChoices {
		return nil, io.ErrUnexpectedEOF
	}
	body, err := io.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}
	if strings.TrimSpace(string(body)) == "" {
		return nil, io.ErrUnexpectedEOF
	}
	return body, nil
}

// parseProjectSummarySourceMeta 从源数据附件中解析项目总结的基础元信息
func parseProjectSummarySourceMeta(sourceFile []byte) (*projectSummarySourceMeta, error) {
	meta := &projectSummarySourceMeta{}
	if err := json.Unmarshal(sourceFile, meta); err != nil {
		return nil, err
	}
	if meta.GroupCount <= 0 {
		return nil, errors.New("source group count is required")
	}
	if meta.ImageCount <= 0 {
		return nil, errors.New("source image count is required")
	}
	return meta, nil
}

// validateProjectSummaryOutput 校验输出报告中的基础元信息是否与源数据一致
func validateProjectSummaryOutput(output string, meta *projectSummarySourceMeta) error {
	result := &projectSummarySourceMeta{}
	if err := json.Unmarshal([]byte(output), result); err != nil {
		return err
	}
	if result.GroupCount != meta.GroupCount {
		return fmt.Errorf("project summary group count mismatch: got %d, want %d", result.GroupCount, meta.GroupCount)
	}
	if result.ImageCount != meta.ImageCount {
		return fmt.Errorf("project summary image count mismatch: got %d, want %d", result.ImageCount, meta.ImageCount)
	}
	return nil
}
