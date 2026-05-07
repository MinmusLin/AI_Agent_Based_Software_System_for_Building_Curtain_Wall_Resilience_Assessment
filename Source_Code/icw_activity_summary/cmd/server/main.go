package main

import (
	"context"
	"net"

	"google.golang.org/grpc"

	"icw_common/consts"
	"icw_common/env"
	"icw_common/utils"

	"icw_activity_summary/configs"
	"icw_activity_summary/internal/services"
	"icw_activity_summary/internal/services/common"
	"icw_activity_summary/rpc/icw_core_biz"
)

// main icw.activity.summary 服务入口
func main() {
	utils.LogInfo(consts.LogScopeInit, "", "Initializing service %s...", consts.ActivitySummaryPSM)

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

	// 注册 gRPC 服务
	grpcServer := grpc.NewServer()
	services.RegisterRPCServices(ctx, common.NewDeps(cfg, coreBizClient), grpcServer)

	// 运行 icw.activity.summary 服务
	utils.LogInfo(consts.LogScopeRPC, consts.LogColorBoldGreen, "icw.activity.summary service starts running on %s", cfg.ActivitySummaryAddr)
	listener, err := net.Listen("tcp", cfg.ActivitySummaryAddr)
	if err != nil {
		utils.LogFatal(consts.LogScopeInit, "Failed to run icw.activity.summary service: %v", err)
	}
	if err := grpcServer.Serve(listener); err != nil {
		utils.LogFatal(consts.LogScopeInit, "Failed to run icw.activity.summary service: %v", err)
	}
}
