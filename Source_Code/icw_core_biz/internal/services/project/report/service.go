package report

import (
	"context"

	"icw_core_biz/internal/services/common"
)

// Service 评估报告服务
type Service struct {
	*common.BaseService
}

type PingRequest struct{}

type PingResponse struct{}

func NewService(ctx context.Context, deps *common.Deps) *Service {
	return &Service{
		BaseService: common.NewBaseService(ctx, deps),
	}
}

// Ping .
func (s *Service) Ping(req *PingRequest, resp *PingResponse) error {
	return s.CallRPC(req, resp, func() error {
		return s.ping(req, resp)
	})
}

func (s *Service) ping(_ *PingRequest, _ *PingResponse) error {
	return nil
}
