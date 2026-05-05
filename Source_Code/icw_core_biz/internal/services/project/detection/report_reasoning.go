package detection

import (
	"icw_core_biz/pkg/dto/project"
)

// ReportReasoningResult 上报推理结果
func (s *Service) ReportReasoningResult(req *project.ReportReasoningResultRequest, resp *project.ReportReasoningResultResponse) error {
	return s.CallRPC("ProjectDetectionService.ReportReasoningResult", req, resp, func() error {
		return s.reportReasoningResult(req, resp)
	})
}

func (s *Service) reportReasoningResult(_ *project.ReportReasoningResultRequest, _ *project.ReportReasoningResultResponse) error {
	return nil
}
