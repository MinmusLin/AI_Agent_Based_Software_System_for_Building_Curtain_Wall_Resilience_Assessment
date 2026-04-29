package handlers

import (
	"net/rpc"

	"github.com/gin-gonic/gin"

	"icw_core_api/configs"
	authHandlers "icw_core_api/internal/handlers/auth"
	"icw_core_api/internal/middlewares"
)

// RegisterRoutes 注册路由
func RegisterRoutes(router *gin.Engine, cfg configs.Config, coreBizClient *rpc.Client) {
	// 登录鉴权 Handler
	authHandler := authHandlers.NewHandler(cfg, coreBizClient)
	auth := router.Group("/auth")
	{
		protectedAuth := auth.Group("")
		// 登录
		auth.POST("/login", authHandler.Login)
		// 登出
		auth.POST("/logout", authHandler.Logout)
		protectedAuth.Use(middlewares.AuthRequired(coreBizClient))
		{
			// 获取用户信息（登录态接口，不允许匿名访问）
			protectedAuth.GET("/me", authHandler.Me)
		}
		// 刷新 Token
		auth.POST("/refresh", authHandler.Refresh)
		// 注册
		auth.POST("/register", authHandler.Register)
		// 重置密码
		auth.POST("/reset-password", authHandler.ResetPassword)
		// 发送邮箱验证码
		auth.POST("/send-email-code", authHandler.SendEmailCode)
	}
}
