package detection

import (
	"context"

	"icw_common/gen/core/biz"
)

// ReportReasoningResult 上报图像检测推理结果
func (s *Service) ReportReasoningResult(ctx context.Context, req *bizpb.ReportReasoningResultRequest) (*bizpb.ReportReasoningResultResponse, error) {
	resp := &bizpb.ReportReasoningResultResponse{}
	err := s.CallRPC(ctx, req, resp, func() error {
		return s.reportReasoningResult(req, resp)
	})
	return resp, err
}

func (s *Service) reportReasoningResult(_ *bizpb.ReportReasoningResultRequest, _ *bizpb.ReportReasoningResultResponse) error {
	return nil
}
