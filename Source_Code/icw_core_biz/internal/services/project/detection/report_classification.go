package detection

import (
	"context"

	"icw_common/gen/core/biz"
)

// ReportClassificationResult 上报图像检测分类结果
func (s *Service) ReportClassificationResult(ctx context.Context, req *bizpb.ReportClassificationResultRequest) (*bizpb.ReportClassificationResultResponse, error) {
	resp := &bizpb.ReportClassificationResultResponse{}
	err := s.CallRPC(ctx, req, resp, func() error {
		return s.reportClassificationResult(req, resp)
	})
	return resp, err
}

func (s *Service) reportClassificationResult(_ *bizpb.ReportClassificationResultRequest, _ *bizpb.ReportClassificationResultResponse) error {
	return nil
}
