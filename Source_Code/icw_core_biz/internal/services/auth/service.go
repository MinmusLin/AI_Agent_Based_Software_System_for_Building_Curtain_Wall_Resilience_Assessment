package auth

import (
	"context"

	"icw_core_biz/internal/services/auth/utils"
	"icw_core_biz/internal/services/common"
)

// Service 登录鉴权服务
type Service struct {
	*common.BaseService
	tokens *utils.TokenManager
}

func NewService(ctx context.Context, deps *common.Deps) *Service {
	if deps == nil {
		deps = &common.Deps{}
	}
	return &Service{
		BaseService: common.NewBaseService(ctx, deps),
		tokens:      utils.NewTokenManager(deps.Config.JWTSecret, deps.Config.AccessTokenTTL),
	}
}
