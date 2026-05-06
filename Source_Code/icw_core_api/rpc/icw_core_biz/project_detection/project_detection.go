package project_detection

import (
	"context"

	"icw_common/gen/core/biz"
	"icw_common/rpc"
	"icw_core_api/rpc/icw_core_biz"
	"icw_core_api/utils"
)

// StartProjectDetection 启动项目智能检测
func StartProjectDetection(ctx context.Context, client *icw_core_biz.Client, req *bizpb.StartProjectDetectionRequest, resp *bizpb.StartProjectDetectionResponse) error {
	return rpc.CallGRPC[bizpb.StartProjectDetectionRequest, bizpb.StartProjectDetectionResponse](ctx, client, req, resp, client.ProjectDetection().StartProjectDetection, rpc.WithRequestIdResolver(utils.GetXRequestId))
}
