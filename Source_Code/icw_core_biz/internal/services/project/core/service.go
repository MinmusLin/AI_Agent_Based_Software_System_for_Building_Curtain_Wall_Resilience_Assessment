package core

import (
	"context"

	"icw_common/gen/core/biz"
	"icw_core_biz/internal/services/common"
)

// Service 项目核心服务
type Service struct {
	bizpb.UnimplementedProjectCoreServiceServer
	*common.BaseService
}

// NewService 创建项目核心服务
func NewService(ctx context.Context, deps *common.Deps) *Service {
	return &Service{
		BaseService: common.NewBaseService(ctx, deps),
	}
}
