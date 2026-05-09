package report

import (
	"context"

	"icw_common/gen/activity"
	"icw_common/gen/core/biz"
	"icw_common/rpc/error"

	"icw_core_biz/internal/services/project/events"
)

// ReportProjectSummaryResult 上报项目总结结果
func (s *Service) ReportProjectSummaryResult(ctx context.Context, req *bizpb.ReportProjectSummaryResultRequest) (*bizpb.ReportProjectSummaryResultResponse, error) {
	resp := &bizpb.ReportProjectSummaryResultResponse{}
	err := s.CallRPC(req, func() error {
		return s.reportProjectSummaryResult(req)
	})
	return resp, err
}

func (s *Service) reportProjectSummaryResult(req *bizpb.ReportProjectSummaryResultRequest) error {
	var reportStatus bizpb.ProjectReportStatus_Value
	switch req.Status {
	case activitypb.DetectionStatus_Succeeded:
		reportStatus = bizpb.ProjectReportStatus_Succeeded
	case activitypb.DetectionStatus_Failed:
		reportStatus = bizpb.ProjectReportStatus_Failed
	default:
		return rpc_error.BadRequestDefault("project summary status is invalid")
	}

	// 按项目 ID 更新项目评估报告结果
	report, err := s.MySQL().UpdateProjectReportResult(s.Ctx(), req.ProjectId, reportStatus, req.ResultJson)
	if err != nil {
		return err
	}

	if report != nil {
		// 发布项目评估报告状态变化事件
		events.PublishProjectReportStatusChangedEvent(
			s.Ctx(),
			s.RocketMQ(),
			report.UserId,
			report.ProjectId,
			report.Uuid,
			report.Status,
		)
	}

	return err
}
