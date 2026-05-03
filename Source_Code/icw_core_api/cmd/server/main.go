package main

import (
	"log"

	"github.com/gin-gonic/gin"

	"icw_core_api/configs"
	"icw_core_api/internal/handlers"
	"icw_core_api/internal/handlers/common"
	"icw_core_api/internal/middlewares"
)

// main icw.core.api 服务入口
func main() {
	// 加载服务配置
	configs.LoadDotEnv(".env")
	cfg, err := configs.Load()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	// 初始化 icw.core.biz 服务
	coreBizClient, err := common.NewRPCClient(common.CoreBizPSM, cfg.CoreBizAddr)
	if err != nil {
		log.Fatalf("Failed to connect to icw.core.biz service: %v", err)
	}
	defer func(coreBizRPCClient *common.RPCClient) {
		_ = coreBizRPCClient.Close()
	}(coreBizClient)

	// 初始化路由
	gin.SetMode(gin.ReleaseMode)
	router := gin.New()
	router.Use(gin.Logger(), gin.Recovery(), middlewares.CORS())
	handlers.RegisterRoutes(router, cfg, coreBizClient)

	// 运行 icw.core.api 服务
	log.Printf("icw.core.api service starts running on %s", cfg.CoreApiAddr)
	if err := router.Run(cfg.CoreApiAddr); err != nil {
		log.Fatalf("Failed to run icw.core.api service: %v", err)
	}
}
