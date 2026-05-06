package review

import (
	"context"

	"icw_common/gen/core/biz"
	"icw_core_biz/internal/services/common"
)

// Service 人工复核服务
type Service struct {
	bizpb.UnimplementedProjectReviewServiceServer
	*common.BaseService
}

// NewService 创建人工复核服务
func NewService(ctx context.Context, deps *common.Deps) *Service {
	if deps == nil {
		deps = &common.Deps{}
	}
	return &Service{
		BaseService: common.NewBaseService(ctx, deps),
	}
}

// Ping .
func (s *Service) Ping(ctx context.Context, req *bizpb.PingReviewRequest) (*bizpb.PingReviewResponse, error) {
	resp := &bizpb.PingReviewResponse{}
	err := s.CallRPC(ctx, req, func() error {
		return s.ping(req, resp)
	})
	return resp, err
}

func (s *Service) ping(_ *bizpb.PingReviewRequest, _ *bizpb.PingReviewResponse) error {
	return nil
}
