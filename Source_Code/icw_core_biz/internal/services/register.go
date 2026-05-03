package services

import (
	"context"
	"log"
	"net/rpc"

	"icw_core_biz/internal/services/auth"
	"icw_core_biz/internal/services/common"
	"icw_core_biz/internal/services/project/assets"
	"icw_core_biz/internal/services/project/core"
	"icw_core_biz/internal/services/project/detection"
	"icw_core_biz/internal/services/project/profile"
	"icw_core_biz/internal/services/project/report"
	"icw_core_biz/internal/services/project/review"
	"icw_core_biz/internal/services/user"
)

// RegisterRPCServices 注册 RPC 服务
func RegisterRPCServices(ctx context.Context, serviceDeps *common.Deps) {
	// 注册用户业务 Service
	userService := user.NewService(ctx, serviceDeps)
	if err := rpc.RegisterName("UserService", userService); err != nil {
		log.Fatalf("Failed to register user rpc service: %v", err)
	}

	// 注册登录鉴权 Service
	authService := auth.NewService(ctx, serviceDeps)
	if err := rpc.RegisterName("AuthService", authService); err != nil {
		log.Fatalf("Failed to register auth rpc service: %v", err)
	}

	// 注册项目核心 Service
	projectCoreService := core.NewService(ctx, serviceDeps)
	if err := rpc.RegisterName("ProjectCoreService", projectCoreService); err != nil {
		log.Fatalf("Failed to register project core rpc service: %v", err)
	}

	// 注册基础信息 Service
	projectProfileService := profile.NewService(ctx, serviceDeps)
	if err := rpc.RegisterName("ProjectProfileService", projectProfileService); err != nil {
		log.Fatalf("Failed to register project profile rpc service: %v", err)
	}

	// 注册图像资产 Service
	projectAssetsService := assets.NewService(ctx, serviceDeps)
	if err := rpc.RegisterName("ProjectAssetsService", projectAssetsService); err != nil {
		log.Fatalf("Failed to register project assets rpc service: %v", err)
	}

	// 注册智能检测 Service
	projectDetectionService := detection.NewService(ctx, serviceDeps)
	if err := rpc.RegisterName("ProjectDetectionService", projectDetectionService); err != nil {
		log.Fatalf("Failed to register project detection rpc service: %v", err)
	}

	// 注册人工复核 Service
	projectReviewService := review.NewService(ctx, serviceDeps)
	if err := rpc.RegisterName("ProjectReviewService", projectReviewService); err != nil {
		log.Fatalf("Failed to register project review rpc service: %v", err)
	}

	// 注册评估报告 Service
	projectReportService := report.NewService(ctx, serviceDeps)
	if err := rpc.RegisterName("ProjectReportService", projectReportService); err != nil {
		log.Fatalf("Failed to register project report rpc service: %v", err)
	}
}
