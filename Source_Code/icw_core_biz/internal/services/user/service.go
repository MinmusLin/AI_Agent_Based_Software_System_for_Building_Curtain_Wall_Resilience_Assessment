package user

import (
	"icw_core_biz/internal/services/common"
)

// Service 用户业务服务
type Service struct {
	*common.BaseService
}

func NewService(deps *common.Deps) *Service {
	if deps == nil {
		deps = &common.Deps{}
	}
	return &Service{
		BaseService: common.NewBaseService(deps),
	}
}
