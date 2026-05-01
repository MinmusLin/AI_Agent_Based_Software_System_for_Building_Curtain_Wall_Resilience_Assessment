package detection

import (
	"icw_core_biz/internal/services/common"
)

// Service 智能检测 Service
type Service struct {
	*common.BaseService
}

func NewService(deps *common.Deps) *Service {
	return &Service{
		BaseService: common.NewBaseService(deps),
	}
}
