package core

import (
	"icw_core_biz/internal/services/common"
)

// Service 项目核心 Service
type Service struct {
	*common.BaseService
}

func NewService(deps *common.Deps) *Service {
	return &Service{
		BaseService: common.NewBaseService(deps),
	}
}
