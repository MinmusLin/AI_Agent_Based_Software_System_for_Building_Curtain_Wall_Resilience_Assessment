package reasoning

import (
	"context"

	"icw_activity_reasoning/internal/services/common"
)

// Service 推理能力服务
type Service struct {
	*common.BaseService
}

func NewService(ctx context.Context, deps *common.Deps) *Service {
	return &Service{
		BaseService: common.NewBaseService(ctx, deps),
	}
}
