package report

import (
	"context"

	"icw_common/enum"
	"icw_common/gen/activity"
	"icw_common/gen/core/biz"
	"icw_common/rpc/error"
)

// ReportProjectSummaryResult 上报项目总结结果
func (s *Service) ReportProjectSummaryResult(ctx context.Context, req *bizpb.ReportProjectSummaryResultRequest) (*bizpb.ReportProjectSummaryResultResponse, error) {
	resp := &bizpb.ReportProjectSummaryResultResponse{}
	err := s.CallRPC(ctx, req, func() error {
		return s.reportProjectSummaryResult(req)
	})
	return resp, err
}

func (s *Service) reportProjectSummaryResult(req *bizpb.ReportProjectSummaryResultRequest) error {
	switch enum.ParseDetectionStatus(req.Status) {
	case activitypb.DetectionStatus_Succeeded:
		return nil
	case activitypb.DetectionStatus_Failed:
		return nil
	default:
		return rpc_error.BadRequestDefault("summary status is invalid")
	}
}
