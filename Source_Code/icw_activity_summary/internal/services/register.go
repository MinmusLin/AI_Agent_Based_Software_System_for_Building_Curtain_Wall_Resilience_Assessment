package services

import (
	"context"

	"google.golang.org/grpc"

	"icw_common/gen/activity/summary"
	"icw_common/rpc"

	"icw_activity_summary/internal/services/common"
	"icw_activity_summary/internal/services/summary"
)

// RegisterRPCServices 注册 gRPC 服务
func RegisterRPCServices(ctx context.Context, serviceDeps *common.Deps, grpcServer *grpc.Server) {
	rpc.RegisterServices(grpcServer, registry(ctx, serviceDeps))
}

// registry RPC 服务注册表
func registry(ctx context.Context, serviceDeps *common.Deps) []rpc.ServiceMeta {
	return []rpc.ServiceMeta{
		{
			Name:        "SummaryService",
			Description: "总结能力服务",
			Service:     summary.NewService(ctx, serviceDeps),
			Register: func(server grpc.ServiceRegistrar, service interface{}) {
				summarypb.RegisterSummaryServiceServer(server, service.(*summary.Service))
			},
			Methods: []rpc.MethodMeta{
				{Name: "StartDetectionSummary", Description: "启动图像检测总结任务"},
				{Name: "StartProjectSummary", Description: "启动项目总结任务"},
			},
		},
	}
}
