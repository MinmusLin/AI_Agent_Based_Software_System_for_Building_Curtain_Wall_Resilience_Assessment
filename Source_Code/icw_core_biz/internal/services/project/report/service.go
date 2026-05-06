package report

import (
	"context"

	"icw_common/gen/core/biz"
	"icw_core_biz/internal/services/common"
)

// Service 评估报告服务
type Service struct {
	bizpb.UnimplementedProjectReportServiceServer
	*common.BaseService
}

// NewService 创建评估报告服务
func NewService(ctx context.Context, deps *common.Deps) *Service {
	return &Service{
		BaseService: common.NewBaseService(ctx, deps),
	}
}

// Ping .
func (s *Service) Ping(ctx context.Context, req *bizpb.PingReportRequest) (*bizpb.PingReportResponse, error) {
	resp := &bizpb.PingReportResponse{}
	err := s.CallRPC(ctx, req, resp, func() error {
		return s.ping(req, resp)
	})
	return resp, err
}

func (s *Service) ping(_ *bizpb.PingReportRequest, _ *bizpb.PingReportResponse) error {
	return nil
}
