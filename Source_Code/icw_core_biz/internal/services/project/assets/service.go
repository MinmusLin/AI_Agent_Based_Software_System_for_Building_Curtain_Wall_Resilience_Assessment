package assets

import (
	"context"

	"icw_common/gen/core/biz"
	"icw_core_biz/internal/services/common"
)

// Service 图像资产服务
type Service struct {
	bizpb.UnimplementedProjectAssetsServiceServer
	*common.BaseService
}

func NewService(ctx context.Context, deps *common.Deps) *Service {
	return &Service{
		BaseService: common.NewBaseService(ctx, deps),
	}
}
