package project_detection

import (
	"context"

	"icw_common/gen/core/biz"
	"icw_common/rpc"

	"icw_core_api/rpc/icw_core_biz"
	"icw_core_api/utils"
)

// GetProjectDetectionTasks 获取项目检测任务列表
func GetProjectDetectionTasks(ctx context.Context, client *icw_core_biz.Client, req *bizpb.GetProjectDetectionTasksRequest, resp *bizpb.GetProjectDetectionTasksResponse) error {
	return rpc.CallGRPC[bizpb.GetProjectDetectionTasksRequest, bizpb.GetProjectDetectionTasksResponse](ctx, client, req, resp, client.ProjectDetection().GetProjectDetectionTasks, rpc.WithRequestIdResolver(utils.GetXRequestId))
}

// RetryProjectDetection 重试项目智能检测
func RetryProjectDetection(ctx context.Context, client *icw_core_biz.Client, req *bizpb.RetryProjectDetectionRequest, resp *bizpb.RetryProjectDetectionResponse) error {
	return rpc.CallGRPC[bizpb.RetryProjectDetectionRequest, bizpb.RetryProjectDetectionResponse](ctx, client, req, resp, client.ProjectDetection().RetryProjectDetection, rpc.WithRequestIdResolver(utils.GetXRequestId))
}

// StartProjectDetection 启动项目智能检测
func StartProjectDetection(ctx context.Context, client *icw_core_biz.Client, req *bizpb.StartProjectDetectionRequest, resp *bizpb.StartProjectDetectionResponse) error {
	return rpc.CallGRPC[bizpb.StartProjectDetectionRequest, bizpb.StartProjectDetectionResponse](ctx, client, req, resp, client.ProjectDetection().StartProjectDetection, rpc.WithRequestIdResolver(utils.GetXRequestId))
}
