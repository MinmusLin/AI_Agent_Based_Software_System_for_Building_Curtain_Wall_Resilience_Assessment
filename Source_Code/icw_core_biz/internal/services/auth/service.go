package auth

import (
	"icw_core_biz/internal/services/auth/utils"
	"icw_core_biz/internal/services/common"
)

// Service 登录鉴权 Service
type Service struct {
	*common.BaseService
	tokens *utils.TokenManager
}

func NewService(deps *common.Deps) *Service {
	if deps == nil {
		deps = &common.Deps{}
	}
	return &Service{
		BaseService: common.NewBaseService(deps),
		tokens:      utils.NewTokenManager(deps.Config.JWTSecret, deps.Config.AccessTokenTTL),
	}
}
