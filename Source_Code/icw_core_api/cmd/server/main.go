package main

import (
	"log"
	"net/rpc"

	"github.com/gin-gonic/gin"

	"icw_core_api/configs"
	"icw_core_api/internal/handlers"
	"icw_core_api/internal/middlewares"
)

func main() {
	// Load configs
	configs.LoadDotEnv(".env")
	cfg := configs.Load()

	// Initialize icw.core.biz service
	coreBizClient, err := rpc.Dial("tcp", cfg.CoreBizAddr)
	if err != nil {
		log.Fatalf("Failed to connect to icw.core.biz service: %v", err)
	}
	defer func(coreBizClient *rpc.Client) {
		_ = coreBizClient.Close()
	}(coreBizClient)

	// Initialize routes
	router := gin.New()
	router.Use(gin.Logger(), gin.Recovery(), middlewares.CORS())
	handlers.RegisterRoutes(router, cfg, coreBizClient)

	log.Printf("icw.core.api service starts running on %s", cfg.CoreApiAddr)

	// Run icw.core.api service
	if err := router.Run(cfg.CoreApiAddr); err != nil {
		log.Fatalf("Failed to run icw.core.api service: %v", err)
	}
}
