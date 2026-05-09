package report

import (
	"context"
	"database/sql"
	"errors"

	"icw_common/gen/core/biz"

	"icw_core_biz/repositories/mysql/model"
)

// GetProjectReport 获取项目评估报告
func (s *Service) GetProjectReport(ctx context.Context, req *bizpb.GetProjectReportRequest) (*bizpb.GetProjectReportResponse, error) {
	resp := &bizpb.GetProjectReportResponse{}
	err := s.CallRPC(req, func() error {
		return s.getProjectReport(ctx, req, resp)
	})
	return resp, err
}

func (s *Service) getProjectReport(ctx context.Context, req *bizpb.GetProjectReportRequest, resp *bizpb.GetProjectReportResponse) error {
	report, err := s.MySQL().GetProjectReport(ctx, req.UserId, req.ProjectId)
	if errors.Is(err, sql.ErrNoRows) {
		return nil
	}
	if err != nil {
		return err
	}
	resp.Report = model.ProjectReportRecordToDTO(report)
	return nil
}
