package services

import (
	"context"

	"google.golang.org/grpc"

	"icw_activity_reasoning/internal/services/common"
	"icw_activity_reasoning/internal/services/reasoning"
	"icw_common/gen/activity/reasoning"
	"icw_common/rpc"
)

// RegisterRPCServices 注册 gRPC 服务
func RegisterRPCServices(ctx context.Context, serviceDeps *common.Deps, grpcServer *grpc.Server) {
	rpc.RegisterServices(grpcServer, registry(ctx, serviceDeps))
}

// registry RPC 服务注册表
func registry(ctx context.Context, serviceDeps *common.Deps) []rpc.ServiceMeta {
	return []rpc.ServiceMeta{
		{
			Name:        "ReasoningService",
			Description: "推理能力服务",
			Service:     reasoning.NewService(ctx, serviceDeps),
			Register: func(server grpc.ServiceRegistrar, service interface{}) {
				reasoningpb.RegisterReasoningServiceServer(server, service.(*reasoning.Service))
			},
			Methods: []rpc.MethodMeta{
				{Name: "Start", Description: "启动推理任务"},
			},
		},
	}
}
