package detection

import (
	"icw_core_biz/pkg/dto/project"
)

// ReportClassificationResult 上报分类结果
func (s *Service) ReportClassificationResult(req *project.ReportClassificationResultRequest, resp *project.ReportClassificationResultResponse) error {
	return s.CallRPC("ProjectDetectionService.ReportClassificationResult", req, resp, func() error {
		return s.reportClassificationResult(req, resp)
	})
}

func (s *Service) reportClassificationResult(_ *project.ReportClassificationResultRequest, _ *project.ReportClassificationResultResponse) error {
	return nil
}
