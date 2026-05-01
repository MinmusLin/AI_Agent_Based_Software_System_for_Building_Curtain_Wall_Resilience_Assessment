package report

import (
	"icw_core_biz/internal/services/common"
)

// Service 评估报告 Service
type Service struct {
	*common.BaseService
}

func NewService(deps *common.Deps) *Service {
	return &Service{
		BaseService: common.NewBaseService(deps),
	}
}
