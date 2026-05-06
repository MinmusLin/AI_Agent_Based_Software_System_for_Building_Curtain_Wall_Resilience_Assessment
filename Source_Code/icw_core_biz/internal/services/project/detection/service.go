package detection

import (
	"context"

	"icw_common/gen/core/biz"
	"icw_core_biz/internal/services/common"
)

// Service 智能检测服务
type Service struct {
	bizpb.UnimplementedProjectDetectionServiceServer
	*common.BaseService
}

func NewService(ctx context.Context, deps *common.Deps) *Service {
	return &Service{
		BaseService: common.NewBaseService(ctx, deps),
	}
}
