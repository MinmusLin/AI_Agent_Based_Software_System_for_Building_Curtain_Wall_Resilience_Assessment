package services

import (
	"context"

	"google.golang.org/grpc"

	"icw_common/gen/core/biz"
	"icw_core_biz/internal/services/auth"
	"icw_core_biz/internal/services/common"
	"icw_core_biz/internal/services/project/assets"
	"icw_core_biz/internal/services/project/core"
	"icw_core_biz/internal/services/project/detection"
	"icw_core_biz/internal/services/project/profile"
	"icw_core_biz/internal/services/project/report"
	"icw_core_biz/internal/services/project/review"
	"icw_core_biz/internal/services/socket"
	"icw_core_biz/internal/services/user"
)

// RegisterRPCServices 注册 RPC 服务
func RegisterRPCServices(ctx context.Context, serviceDeps *common.Deps, grpcServer *grpc.Server) {
	registeredMethods := make([]common.RegisteredRPCMethodMeta, 0)
	for _, meta := range registry(ctx, serviceDeps) {
		methods, err := common.ResolveRPCMethods(meta)
		if err != nil {
			common.RpcFatal("Failed to register RPC service %s: %v", meta.Name, err)
		}
		if meta.Register == nil {
			common.RpcFatal("Failed to register RPC service %s: register function is nil", meta.Name)
		}
		meta.Register(grpcServer, meta.Service)
		registeredMethods = append(registeredMethods, methods...)
	}
	common.RpcInfo("RPC methods registered, waiting for requests:\n%s", common.FormatRegistryTable(registeredMethods))
}

// registry RPC 服务注册表
func registry(ctx context.Context, serviceDeps *common.Deps) []common.RPCServiceMeta {
	return []common.RPCServiceMeta{
		{
			Name:        "SocketService",
			Description: "WebSocket 连接票据服务",
			Service:     socket.NewService(ctx, serviceDeps),
			Register: func(server grpc.ServiceRegistrar, service interface{}) {
				bizpb.RegisterSocketServiceServer(server, service.(*socket.Service))
			},
			Methods: []common.RPCMethodMeta{
				{Name: "CreateSocketTicket", Description: "创建 WebSocket 连接票据"},
				{Name: "ValidateSocketTicket", Description: "校验 WebSocket 连接票据"},
			},
		},
		{
			Name:        "UserService",
			Description: "用户业务服务",
			Service:     user.NewService(ctx, serviceDeps),
			Register: func(server grpc.ServiceRegistrar, service interface{}) {
				bizpb.RegisterUserServiceServer(server, service.(*user.Service))
			},
			Methods: []common.RPCMethodMeta{
				{Name: "GetAvatar", Description: "获取用户头像"},
				{Name: "UploadAvatar", Description: "上传用户自定义头像"},
				{Name: "DeleteAvatar", Description: "删除用户自定义头像"},
			},
		},
		{
			Name:        "AuthService",
			Description: "登录鉴权服务",
			Service:     auth.NewService(ctx, serviceDeps),
			Register: func(server grpc.ServiceRegistrar, service interface{}) {
				bizpb.RegisterAuthServiceServer(server, service.(*auth.Service))
			},
			Methods: []common.RPCMethodMeta{
				{Name: "Login", Description: "登录"},
				{Name: "Logout", Description: "登出"},
				{Name: "Me", Description: "获取用户信息"},
				{Name: "Refresh", Description: "刷新 Token"},
				{Name: "Register", Description: "注册"},
				{Name: "ResetPassword", Description: "重置密码"},
				{Name: "SendEmailCode", Description: "发送邮箱验证码"},
			},
		},
		{
			Name:        "ProjectCoreService",
			Description: "项目核心服务",
			Service:     core.NewService(ctx, serviceDeps),
			Register: func(server grpc.ServiceRegistrar, service interface{}) {
				bizpb.RegisterProjectCoreServiceServer(server, service.(*core.Service))
			},
			Methods: []common.RPCMethodMeta{
				{Name: "AdvanceProject", Description: "项目进度流转"},
				{Name: "CreateProject", Description: "创建项目"},
				{Name: "DeleteProject", Description: "删除项目"},
				{Name: "ListProjects", Description: "获取项目列表"},
				{Name: "CheckProjectAccess", Description: "校验项目访问权限"},
			},
		},
		{
			Name:        "ProjectProfileService",
			Description: "基础信息服务",
			Service:     profile.NewService(ctx, serviceDeps),
			Register: func(server grpc.ServiceRegistrar, service interface{}) {
				bizpb.RegisterProjectProfileServiceServer(server, service.(*profile.Service))
			},
			Methods: []common.RPCMethodMeta{
				{Name: "GetProjectProfile", Description: "获取项目基础信息"},
				{Name: "GetProjectThumbnail", Description: "获取项目缩略图"},
				{Name: "UploadProjectThumbnail", Description: "上传项目缩略图"},
				{Name: "DeleteProjectThumbnail", Description: "删除项目缩略图"},
				{Name: "UpdateProjectProfile", Description: "更新项目基础信息"},
			},
		},
		{
			Name:        "ProjectAssetsService",
			Description: "图像资产服务",
			Service:     assets.NewService(ctx, serviceDeps),
			Register: func(server grpc.ServiceRegistrar, service interface{}) {
				bizpb.RegisterProjectAssetsServiceServer(server, service.(*assets.Service))
			},
			Methods: []common.RPCMethodMeta{
				{Name: "GetProjectAssets", Description: "获取项目图像列表"},
				{Name: "CreateProjectGroup", Description: "创建图像组"},
				{Name: "DeleteProjectGroup", Description: "删除图像组"},
				{Name: "MoveProjectGroup", Description: "移动图像组"},
				{Name: "UpdateProjectGroup", Description: "更新图像组"},
				{Name: "DeleteProjectImage", Description: "删除图像"},
				{Name: "GetProjectImageOriginal", Description: "获取原图"},
				{Name: "MoveProjectImage", Description: "移动图像"},
				{Name: "ReportProjectImage", Description: "上报图像"},
				{Name: "UploadProjectImage", Description: "上传图像"},
			},
		},
		{
			Name:        "ProjectDetectionService",
			Description: "智能检测服务",
			Service:     detection.NewService(ctx, serviceDeps),
			Register: func(server grpc.ServiceRegistrar, service interface{}) {
				bizpb.RegisterProjectDetectionServiceServer(server, service.(*detection.Service))
			},
			Methods: []common.RPCMethodMeta{
				{Name: "ReportClassificationResult", Description: "上报图像检测分类结果"},
				{Name: "ReportReasoningResult", Description: "上报图像检测推理结果"},
			},
		},
		{
			Name:        "ProjectReviewService",
			Description: "人工复核服务",
			Service:     review.NewService(ctx, serviceDeps),
			Register: func(server grpc.ServiceRegistrar, service interface{}) {
				bizpb.RegisterProjectReviewServiceServer(server, service.(*review.Service))
			},
			Methods: []common.RPCMethodMeta{
				{Name: "Ping", Description: "人工复核服务探活"},
			},
		},
		{
			Name:        "ProjectReportService",
			Description: "评估报告服务",
			Service:     report.NewService(ctx, serviceDeps),
			Register: func(server grpc.ServiceRegistrar, service interface{}) {
				bizpb.RegisterProjectReportServiceServer(server, service.(*report.Service))
			},
			Methods: []common.RPCMethodMeta{
				{Name: "Ping", Description: "评估报告服务探活"},
			},
		},
	}
}
