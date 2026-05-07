package user

import (
	"context"

	"icw_common/gen/core/biz"

	"icw_core_biz/internal/services/common"
)

// Service 用户业务服务
type Service struct {
	bizpb.UnimplementedUserServiceServer
	*common.BaseService
}

// NewService 创建用户业务服务
func NewService(ctx context.Context, deps *common.Deps) *Service {
	if deps == nil {
		deps = &common.Deps{}
	}
	return &Service{
		BaseService: common.NewBaseService(ctx, deps),
	}
}
