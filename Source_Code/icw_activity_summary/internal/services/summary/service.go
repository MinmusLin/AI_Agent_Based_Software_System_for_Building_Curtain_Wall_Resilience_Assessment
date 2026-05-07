package summary

import (
	"context"

	"icw_common/gen/activity/summary"

	"icw_activity_summary/internal/services/common"
)

// Service 总结能力服务
type Service struct {
	summarypb.UnimplementedSummaryServiceServer
	*common.BaseService
}

// NewService 创建总结能力服务
func NewService(ctx context.Context, deps *common.Deps) *Service {
	if deps == nil {
		deps = &common.Deps{}
	}
	return &Service{
		BaseService: common.NewBaseService(ctx, deps),
	}
}
