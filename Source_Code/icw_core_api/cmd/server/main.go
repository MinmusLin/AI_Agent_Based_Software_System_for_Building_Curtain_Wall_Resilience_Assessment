package main

import (
	"icw_common/consts"
	"icw_common/env"
	"icw_common/utils"
	"icw_core_api/configs"
	"icw_core_api/internal/handlers"
	"icw_core_api/internal/rocketmq"
	"icw_core_api/internal/socket"
	"icw_core_api/rpc/icw_core_biz"
)

// main icw.core.api 服务入口
func main() {
	utils.LogInfo(consts.LogScopeInit, "", "Initializing service %s...", consts.CoreApiPSM)

	// 加载服务配置
	if err := env.LoadDotEnv(".env"); err != nil {
		utils.LogFatal(consts.LogScopeInit, "Failed to load .env file: %v", err)
	}
	cfg, err := configs.Load()
	if err != nil {
		utils.LogFatal(consts.LogScopeInit, "Failed to load config: %v", err)
	}
	utils.LogInfo(consts.LogScopeInit, "", "Config initialized successfully:\n%s", env.FormatEnvConfig(cfg))

	// 初始化 WebSocket Hub
	webSocketHub := socket.NewHub()
	socket.WSInfo("WebSocket Hub initialized successfully")

	// 初始化 RocketMQ 事件消费者
	eventConsumer, err := rocketmq.NewConsumer(cfg, webSocketHub)
	if err != nil {
		utils.LogFatal(consts.LogScopeInit, "Failed to create RocketMQ event consumer: %v", err)
	}
	if err := eventConsumer.Start(); err != nil {
		utils.LogFatal(consts.LogScopeInit, "Failed to start RocketMQ event consumer: %v", err)
	}
	defer func() {
		_ = eventConsumer.Close()
	}()
	rocketmq.MQInfo("RocketMQ consumer starts running")

	// 初始化 icw.core.biz 服务
	coreBizClient, err := icw_core_biz.NewClient(cfg.CoreBizAddr)
	if err != nil {
		utils.LogFatal(consts.LogScopeInit, "Failed to initialize icw.core.biz service: %v", err)
	}
	utils.LogInfo(consts.LogScopeRPC, consts.LogColorBoldGreen, "RPC service icw.core.biz initialized successfully")
	defer func() {
		_ = coreBizClient.Close()
	}()

	// 初始化路由
	router := handlers.RegisterRoutes(cfg, coreBizClient, webSocketHub)

	// 运行 icw.core.api 服务
	utils.LogInfo(consts.LogScopeHTTP, consts.LogColorBoldGreen, "icw.core.api service starts running on %s", cfg.CoreApiAddr)
	if err := router.Run(cfg.CoreApiAddr); err != nil {
		utils.LogFatal(consts.LogScopeInit, "Failed to run icw.core.api service: %v", err)
	}
}
