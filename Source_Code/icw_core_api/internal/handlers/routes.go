package handlers

import (
	"net/rpc"

	"github.com/gin-gonic/gin"

	"icw_core_api/configs"
	authHandlers "icw_core_api/internal/handlers/auth"
)

func RegisterRoutes(router *gin.Engine, cfg configs.Config, rpcClient *rpc.Client) {
	authHandler := authHandlers.NewHandler(cfg, rpcClient)
	auth := router.Group("/auth")
	{
		auth.POST("/login", authHandler.Login)
		auth.POST("/logout", authHandler.Logout)
		auth.GET("/me", authHandler.Me)
		auth.POST("/refresh", authHandler.Refresh)
		auth.POST("/register", authHandler.Register)
		auth.POST("/reset-password", authHandler.ResetPassword)
		auth.POST("/send-email-code", authHandler.SendEmailCode)
	}
}
