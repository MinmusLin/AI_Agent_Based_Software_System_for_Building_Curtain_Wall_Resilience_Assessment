package user

import (
	"context"

	"icw_core_biz/internal/services/common"
)

// Service 用户业务 Service
type Service struct {
	*common.BaseService
}

func NewService(ctx context.Context, deps *common.Deps) *Service {
	if deps == nil {
		deps = &common.Deps{}
	}
	return &Service{
		BaseService: common.NewBaseService(ctx, deps),
	}
}
