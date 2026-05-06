package classification

import (
	"context"

	"icw_activity_classification/internal/services/common"
	"icw_common/gen/activity/classification"
)

// Service 分类能力服务
type Service struct {
	classificationpb.UnimplementedClassificationServiceServer
	*common.BaseService
}

// NewService 创建分类能力服务
func NewService(ctx context.Context, deps *common.Deps) *Service {
	if deps == nil {
		deps = &common.Deps{}
	}
	return &Service{
		BaseService: common.NewBaseService(ctx, deps),
	}
}
