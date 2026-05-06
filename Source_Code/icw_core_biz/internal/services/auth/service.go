package auth

import (
	"context"

	"icw_common/gen/core/biz"
	"icw_core_biz/internal/services/auth/utils"
	"icw_core_biz/internal/services/common"
)

// Service 登录鉴权服务
type Service struct {
	bizpb.UnimplementedAuthServiceServer
	*common.BaseService
	tokens *utils.TokenManager
}

// NewService 创建登录鉴权服务
func NewService(ctx context.Context, deps *common.Deps) *Service {
	if deps == nil {
		deps = &common.Deps{}
	}
	return &Service{
		BaseService: common.NewBaseService(ctx, deps),
		tokens:      utils.NewTokenManager(deps.Config.JWTSecret, deps.Config.AccessTokenTTL),
	}
}
