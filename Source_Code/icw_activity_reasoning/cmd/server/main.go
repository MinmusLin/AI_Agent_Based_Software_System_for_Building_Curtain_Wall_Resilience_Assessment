package main

import (
	"context"
	"net"
	"net/rpc"

	"icw_activity_reasoning/configs"
	"icw_activity_reasoning/consts"
	"icw_activity_reasoning/internal/detectors"
	"icw_activity_reasoning/internal/services"
	"icw_activity_reasoning/internal/services/common"
	"icw_activity_reasoning/rpc/icw_core_biz"
	bizConfigs "icw_core_biz/configs"
	"icw_core_biz/utils"
)

// main icw.activity.reasoning 服务入口
func main() {
	utils.LogInfo(consts.LogScopeInit, "", "Initializing service %s...", consts.ActivityReasoningPSM)

	ctx := context.Background()

	// 加载服务配置
	bizConfigs.LoadDotEnv(".env")
	cfg, err := configs.Load()
	if err != nil {
		utils.LogFatal(consts.LogScopeInit, "Failed to load config: %v", err)
	}
	utils.LogInfo(consts.LogScopeInit, "", "Config initialized successfully:\n%s", utils.FormatEnvConfig(cfg))

	// 初始化 Python 原子检测能力注册表
	registry := detectors.NewPythonRegistry(cfg.PythonBin, cfg.ReasoningWorkDir, cfg.ReasoningTaskCodes)
	utils.LogInfo(consts.LogScopeInit, "", "Python detector registry initialized successfully:\n%s", detectors.FormatRegistryTable(registry))

	// 初始化 BIZ RPC 回调仓库
	coreBizClient := icw_core_biz.NewClient(cfg.CoreBizAddr)
	utils.LogInfo(consts.LogScopeInit, "", "icw.core.biz RPC client initialized successfully")

	// 注册 RPC 服务
	services.RegisterRPCServices(ctx, common.NewDeps(cfg, registry, coreBizClient))

	// 运行 icw.activity.reasoning RPC 服务
	utils.LogInfo(consts.LogScopeInit, "", "icw.activity.reasoning rpc service starts running on %s", cfg.ActivityReasoningAddr)
	listener, err := net.Listen("tcp", cfg.ActivityReasoningAddr)
	if err != nil {
		utils.LogFatal(consts.LogScopeInit, "Failed to run icw.activity.reasoning service: %v", err)
	}
	rpc.Accept(listener)
}
