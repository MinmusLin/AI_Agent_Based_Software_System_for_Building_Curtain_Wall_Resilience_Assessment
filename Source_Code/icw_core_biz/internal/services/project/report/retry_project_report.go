package report

import (
	"context"

	"icw_common/gen/core/biz"
	"icw_common/rpc/error"
)

// RetryProjectReport 重试项目评估报告生成
func (s *Service) RetryProjectReport(ctx context.Context, req *bizpb.RetryProjectReportRequest) (*bizpb.RetryProjectReportResponse, error) {
	resp := &bizpb.RetryProjectReportResponse{}
	err := s.CallRPC(req, func() error {
		return s.retryProjectReport(ctx, req)
	})
	return resp, err
}

func (s *Service) retryProjectReport(ctx context.Context, req *bizpb.RetryProjectReportRequest) error {
	report, err := s.MySQL().GetProjectReport(ctx, req.UserId, req.ProjectId)
	if err != nil {
		return err
	}
	if report.Status != bizpb.ProjectReportStatus_Failed {
		return rpc_error.BadRequestDefault("project report is not failed")
	}

	if _, err := s.MySQL().CreateProjectReport(ctx, req.UserId, req.ProjectId); err != nil {
		return err
	}

	go s.DetectionWorker().StartProjectReportTask(s.Ctx(), req.UserId, req.ProjectId)

	return nil
}
