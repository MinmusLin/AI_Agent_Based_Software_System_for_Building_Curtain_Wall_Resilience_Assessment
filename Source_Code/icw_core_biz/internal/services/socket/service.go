package socket

import (
	"context"

	"icw_common/gen/core/biz"

	"icw_core_biz/internal/services/common"
)

// Service WebSocket 连接票据服务
type Service struct {
	bizpb.UnimplementedSocketServiceServer
	*common.BaseService
}

// NewService 创建 WebSocket 连接票据服务
func NewService(ctx context.Context, deps *common.Deps) *Service {
	if deps == nil {
		deps = &common.Deps{}
	}
	return &Service{
		BaseService: common.NewBaseService(ctx, deps),
	}
}
