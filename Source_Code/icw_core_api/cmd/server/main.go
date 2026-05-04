package main

import (
	"icw_core_api/configs"
	"icw_core_api/consts"
	"icw_core_api/internal/handlers"
	"icw_core_api/internal/handlers/common"
	"icw_core_api/internal/rocketmq"
	"icw_core_api/internal/socket"
	bizConsts "icw_core_biz/consts"
	"icw_core_biz/utils"
)

// main icw.core.api 服务入口
func main() {
	utils.LogInfo(bizConsts.LogScopeInit, "", "Initializing service %s...", consts.CoreApiPSM)

	// 加载服务配置
	configs.LoadDotEnv(".env")
	cfg, err := configs.Load()
	if err != nil {
		utils.LogFatal(bizConsts.LogScopeInit, "Failed to load config: %v", err)
	}
	utils.LogInfo(bizConsts.LogScopeInit, "", "Config initialized successfully")

	// 初始化 WebSocket Hub
	webSocketHub := socket.NewHub()
	socket.WSInfo("WebSocket Hub initialized successfully")

	// 初始化 RocketMQ 事件消费者
	eventConsumer, err := rocketmq.NewConsumer(cfg, webSocketHub)
	if err != nil {
		utils.LogFatal(bizConsts.LogScopeInit, "Failed to create RocketMQ event consumer: %v", err)
	}
	if err := eventConsumer.Start(); err != nil {
		utils.LogFatal(bizConsts.LogScopeInit, "Failed to start RocketMQ event consumer: %v", err)
	}
	defer func() {
		_ = eventConsumer.Close()
	}()
	rocketmq.MQInfo("RocketMQ consumer starts running")

	// 初始化 icw.core.biz 服务
	coreBizClient, err := common.NewRPCClient(bizConsts.CoreBizPSM, cfg.CoreBizAddr)
	if err != nil {
		utils.LogFatal(bizConsts.LogScopeInit, "Failed to connect to icw.core.biz service: %v", err)
	}
	defer func() {
		_ = coreBizClient.Close()
	}()
	utils.LogInfo(bizConsts.LogScopeRPC, bizConsts.LogColorBoldGreen, "Connected to RPC service icw.core.biz successfully")

	// 初始化路由
	router := handlers.RegisterRoutes(cfg, coreBizClient, webSocketHub)

	// 运行 icw.core.api 服务
	utils.LogInfo(bizConsts.LogScopeHTTP, "", "icw.core.api service starts running on %s", cfg.CoreApiAddr)
	if err := router.Run(cfg.CoreApiAddr); err != nil {
		utils.LogFatal(bizConsts.LogScopeInit, "Failed to run icw.core.api service: %v", err)
	}
}
