package services

import (
	"context"

	"google.golang.org/grpc"

	"icw_activity_reasoning/internal/services/common"
	"icw_activity_reasoning/internal/services/reasoning"
	"icw_common/consts"
	"icw_common/gen/activity/reasoning"
	"icw_common/utils"
)

// RegisterRPCServices 注册 gRPC 服务
func RegisterRPCServices(ctx context.Context, serviceDeps *common.Deps, grpcServer *grpc.Server) {
	registeredMethods := make([]common.RegisteredRPCMethodMeta, 0)
	for _, meta := range registry(ctx, serviceDeps) {
		methods, err := common.ResolveRPCMethods(meta)
		if err != nil {
			utils.LogFatal(consts.LogScopeRPC, "Failed to register RPC service %s: %v", meta.Name, err)
		}
		if meta.Register == nil {
			utils.LogFatal(consts.LogScopeRPC, "Failed to register RPC service %s: register function is nil", meta.Name)
		}
		meta.Register(grpcServer, meta.Service)
		registeredMethods = append(registeredMethods, methods...)
	}
	utils.LogFatal(consts.LogScopeRPC, "RPC methods registered, waiting for requests:\n%s", common.FormatRegistryTable(registeredMethods))
}

// registry RPC 服务注册表
func registry(ctx context.Context, serviceDeps *common.Deps) []common.RPCServiceMeta {
	return []common.RPCServiceMeta{
		{
			Name:        "ReasoningService",
			Description: "推理能力服务",
			Service:     reasoning.NewService(ctx, serviceDeps),
			Register: func(server grpc.ServiceRegistrar, service interface{}) {
				reasoningpb.RegisterReasoningServiceServer(server, service.(*reasoning.Service))
			},
			Methods: []common.RPCMethodMeta{
				{Name: "Start", Description: "启动推理任务"},
			},
		},
	}
}
