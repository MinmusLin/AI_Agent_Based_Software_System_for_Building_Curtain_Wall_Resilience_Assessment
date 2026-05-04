package main

import (
	"github.com/gin-gonic/gin"

	"icw_core_api/configs"
	"icw_core_api/internal/handlers"
	"icw_core_api/internal/handlers/common"
	"icw_core_api/internal/middlewares"
	"icw_core_api/internal/rocketmq"
	"icw_core_api/internal/socket"
	"icw_core_biz/consts"
	"icw_core_biz/utils"
)

// main icw.core.api 服务入口
func main() {
	// 加载服务配置
	configs.LoadDotEnv(".env")
	cfg, err := configs.Load()
	if err != nil {
		utils.LogFatal(consts.LogScopeInit, "Failed to load config: %v", err)
	}

	// 初始化 icw.core.biz 服务
	coreBizClient, err := common.NewRPCClient(common.CoreBizPSM, cfg.CoreBizAddr)
	if err != nil {
		utils.LogFatal(consts.LogScopeInit, "Failed to connect to icw.core.biz service: %v", err)
	}
	defer func() {
		_ = coreBizClient.Close()
	}()

	// 初始化 WebSocket Hub 和 RocketMQ 事件消费者
	webSocketHub := socket.NewHub()
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

	// 初始化路由
	gin.SetMode(gin.ReleaseMode)
	router := gin.New()
	router.Use(middlewares.RequestId(), middlewares.Logger(), gin.Recovery(), middlewares.CORS())
	handlers.RegisterRoutes(router, cfg, coreBizClient, webSocketHub)

	// 运行 icw.core.api 服务
	utils.LogInfo(consts.LogScopeInit, "", "icw.core.api service starts running on %s", cfg.CoreApiAddr)
	if err := router.Run(cfg.CoreApiAddr); err != nil {
		utils.LogFatal(consts.LogScopeInit, "Failed to run icw.core.api service: %v", err)
	}
}
