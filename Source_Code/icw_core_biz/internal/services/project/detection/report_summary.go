package detection

import (
	"icw_core_biz/pkg/dto/project"
)

// ReportSummaryResult 上报图像检测总结结果
func (s *Service) ReportSummaryResult(req *project.ReportSummaryResultRequest, resp *project.ReportSummaryResultResponse) error {
	return s.CallRPC("ProjectDetectionService.ReportSummaryResult", req, resp, func() error {
		return s.reportSummaryResult(req, resp)
	})
}

func (s *Service) reportSummaryResult(_ *project.ReportSummaryResultRequest, _ *project.ReportSummaryResultResponse) error {
	return nil
}
