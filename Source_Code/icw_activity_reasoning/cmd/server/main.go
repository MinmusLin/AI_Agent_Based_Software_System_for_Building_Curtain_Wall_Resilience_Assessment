package main

import (
	"context"
	"net"

	"google.golang.org/grpc"

	"icw_common/consts"
	"icw_common/env"
	"icw_common/utils"

	"icw_activity_reasoning/configs"
	"icw_activity_reasoning/internal/detectors"
	detectorCommon "icw_activity_reasoning/internal/detectors/common"
	"icw_activity_reasoning/internal/services"
	"icw_activity_reasoning/internal/services/common"
	"icw_activity_reasoning/rpc/icw_core_biz"
)

// main icw.activity.reasoning 服务入口
func main() {
	utils.LogInfo(consts.LogScopeInit, "", "Initializing service %s...", consts.ActivityReasoningPSM)

	ctx := context.Background()

	// 加载服务配置
	if err := env.LoadDotEnv(".env"); err != nil {
		utils.LogFatal(consts.LogScopeInit, "Failed to load .env file: %v", err)
	}
	cfg, err := configs.Load()
	if err != nil {
		utils.LogFatal(consts.LogScopeInit, "Failed to load config: %v", err)
	}
	utils.LogInfo(consts.LogScopeInit, "", "Config initialized successfully:\n%s", env.FormatEnvConfig(cfg))

	// 初始化 icw.core.biz 服务
	coreBizClient, err := icw_core_biz.NewClient(cfg.CoreBizAddr)
	if err != nil {
		utils.LogFatal(consts.LogScopeInit, "Failed to initialize icw.core.biz service: %v", err)
	}
	utils.LogInfo(consts.LogScopeRPC, consts.LogColorBoldGreen, "RPC service icw.core.biz initialized successfully")
	defer func() {
		_ = coreBizClient.Close()
	}()

	// 注册 Python 原子检测能力
	detectorsRegistry, err := detectors.RegisterDetectors(cfg.ReasoningRuntimeDir, cfg.ReasoningModelDir)
	if err != nil {
		utils.LogFatal(consts.LogScopeInit, "Failed to register python detectors: %v", err)
	}
	utils.LogInfo(consts.LogScopeReasoning, consts.LogColorBoldPurple, "Python detectors registered, waiting for calls:\n%s", detectorCommon.FormatRegistryTable(detectorsRegistry))

	// 注册 gRPC 服务
	grpcServer := grpc.NewServer()
	services.RegisterRPCServices(ctx, common.NewDeps(cfg, detectorsRegistry, coreBizClient), grpcServer)

	// 运行 icw.activity.reasoning 服务
	utils.LogInfo(consts.LogScopeRPC, consts.LogColorBoldGreen, "icw.activity.reasoning service starts running on %s", cfg.ActivityReasoningAddr)
	listener, err := net.Listen("tcp", cfg.ActivityReasoningAddr)
	if err != nil {
		utils.LogFatal(consts.LogScopeInit, "Failed to run icw.activity.reasoning service: %v", err)
	}
	if err := grpcServer.Serve(listener); err != nil {
		utils.LogFatal(consts.LogScopeInit, "Failed to run icw.activity.reasoning service: %v", err)
	}
}
