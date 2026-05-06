package reasoning

import (
	"context"

	"icw_activity_reasoning/internal/services/common"
	"icw_common/gen/activity/reasoning"
)

// Service 推理能力服务
type Service struct {
	reasoningpb.UnimplementedReasoningServiceServer
	*common.BaseService
}

func NewService(ctx context.Context, deps *common.Deps) *Service {
	if deps == nil {
		deps = &common.Deps{}
	}
	return &Service{
		BaseService: common.NewBaseService(ctx, deps),
	}
}
