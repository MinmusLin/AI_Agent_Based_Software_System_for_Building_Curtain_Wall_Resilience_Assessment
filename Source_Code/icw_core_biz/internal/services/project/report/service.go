package report

import (
	"icw_core_biz/internal/services/common"
)

// Service 评估报告 Service
type Service struct {
	*common.BaseService
}

type PingRequest struct{}

type PingResponse struct{}

func NewService(deps *common.Deps) *Service {
	return &Service{
		BaseService: common.NewBaseService(deps),
	}
}

// Ping .
func (s *Service) Ping(req *PingRequest, resp *PingResponse) error {
	return s.CallRPC("ProjectReportService.Ping", req, resp, func() error {
		return s.ping(req, resp)
	})
}

func (s *Service) ping(_ *PingRequest, _ *PingResponse) error {
	return nil
}
