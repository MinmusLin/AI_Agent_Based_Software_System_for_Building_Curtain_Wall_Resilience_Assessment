package main

import (
	"context"
	"net"
	"net/rpc"

	"icw_activity_reasoning/configs"
	"icw_activity_reasoning/consts"
	"icw_activity_reasoning/internal/detectors"
	"icw_activity_reasoning/internal/detectors/utils"
	"icw_activity_reasoning/internal/services"
	"icw_activity_reasoning/internal/services/common"
	"icw_activity_reasoning/rpc/icw_core_biz"
	bizConfigs "icw_core_biz/configs"
	bizConsts "icw_core_biz/consts"
	bizUtils "icw_core_biz/utils"
)

// main icw.activity.reasoning 服务入口
func main() {
	bizUtils.LogInfo(consts.LogScopeInit, "", "Initializing service %s...", consts.ActivityReasoningPSM)

	ctx := context.Background()

	// 加载服务配置
	bizConfigs.LoadDotEnv(".env")
	cfg, err := configs.Load()
	if err != nil {
		bizUtils.LogFatal(consts.LogScopeInit, "Failed to load config: %v", err)
	}
	bizUtils.LogInfo(consts.LogScopeInit, "", "Config initialized successfully:\n%s", bizUtils.FormatEnvConfig(cfg))

	// 初始化 icw.core.biz 服务
	coreBizClient := icw_core_biz.NewClient(cfg.CoreBizAddr)
	bizUtils.LogInfo(bizConsts.LogScopeRPC, bizConsts.LogColorBoldGreen, "RPC service icw.core.biz initialized successfully")

	// 注册 Python 检测能力
	detectorsRegistry := detectors.RegisterDetectors(cfg.PythonBin, cfg.ReasoningWorkDir)
	bizUtils.LogInfo(consts.LogScopeInit, "", "Python detectors registered, waiting for calls:\n%s", utils.FormatRegistryTable(detectorsRegistry))

	// 注册 RPC 服务
	services.RegisterRPCServices(ctx, common.NewDeps(cfg, detectorsRegistry, coreBizClient))

	// 运行 icw.activity.reasoning 服务
	bizUtils.LogInfo(bizConsts.LogScopeRPC, bizConsts.LogColorBoldGreen, "icw.activity.reasoning service starts running on %s", cfg.ActivityReasoningAddr)
	listener, err := net.Listen("tcp", cfg.ActivityReasoningAddr)
	if err != nil {
		bizUtils.LogFatal(consts.LogScopeInit, "Failed to run icw.activity.reasoning service: %v", err)
	}
	rpc.Accept(listener)
}
