package services

import (
	"context"

	"google.golang.org/grpc"

	"icw_common/gen/activity/classification"
	"icw_common/rpc"

	"icw_activity_classification/internal/services/classification"
	"icw_activity_classification/internal/services/common"
)

// RegisterRPCServices 注册 gRPC 服务
func RegisterRPCServices(ctx context.Context, serviceDeps *common.Deps, grpcServer *grpc.Server) {
	rpc.RegisterServices(grpcServer, registry(ctx, serviceDeps))
}

// registry RPC 服务注册表
func registry(ctx context.Context, serviceDeps *common.Deps) []rpc.ServiceMeta {
	return []rpc.ServiceMeta{
		{
			Name:        "ClassificationService",
			Description: "分类能力服务",
			Service:     classification.NewService(ctx, serviceDeps),
			Register: func(server grpc.ServiceRegistrar, service interface{}) {
				classificationpb.RegisterClassificationServiceServer(server, service.(*classification.Service))
			},
			Methods: []rpc.MethodMeta{
				{Name: "Start", Description: "启动分类任务"},
			},
		},
	}
}
