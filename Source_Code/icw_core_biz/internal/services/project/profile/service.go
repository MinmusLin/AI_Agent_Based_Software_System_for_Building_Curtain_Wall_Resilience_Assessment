package profile

import (
	"icw_core_biz/internal/services/common"
)

// Service 基础信息 Service
type Service struct {
	*common.BaseService
}

func NewService(deps *common.Deps) *Service {
	return &Service{
		BaseService: common.NewBaseService(deps),
	}
}
