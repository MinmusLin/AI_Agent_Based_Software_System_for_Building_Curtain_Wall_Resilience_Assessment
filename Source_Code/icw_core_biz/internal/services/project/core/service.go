package core

import (
	"context"

	"icw_core_biz/internal/services/common"
)

// Service 项目核心服务
type Service struct {
	*common.BaseService
}

func NewService(ctx context.Context, deps *common.Deps) *Service {
	return &Service{
		BaseService: common.NewBaseService(ctx, deps),
	}
}
