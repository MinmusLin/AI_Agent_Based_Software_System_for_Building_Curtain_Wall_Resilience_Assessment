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
	if deps == nil {
		deps = &common.Deps{}
	}
	return &Service{
		BaseService: common.NewBaseService(ctx, deps),
	}
}
