package handlers

import (
	"net/rpc"

	"github.com/gin-gonic/gin"

	"icw_core_api/configs"
	"icw_core_api/internal/handlers/auth"
	"icw_core_api/internal/handlers/user"
	"icw_core_api/internal/middlewares"
)

// RegisterRoutes 注册路由
func RegisterRoutes(router *gin.Engine, cfg configs.Config, coreBizClient *rpc.Client) {
	// 登录鉴权 Handler
	authHandler := auth.NewHandler(cfg, coreBizClient)
	authRouter := router.Group("/auth")
	{
		protectedAuth := authRouter.Group("")
		// 登录
		authRouter.POST("/login", authHandler.Login)
		// 登出
		authRouter.POST("/logout", authHandler.Logout)
		protectedAuth.Use(middlewares.AuthRequired(coreBizClient))
		{
			// 获取用户信息（登录态接口，不允许匿名访问）
			protectedAuth.GET("/me", authHandler.Me)
		}
		// 刷新 Token
		authRouter.POST("/refresh", authHandler.Refresh)
		// 注册
		authRouter.POST("/register", authHandler.Register)
		// 重置密码
		authRouter.POST("/reset-password", authHandler.ResetPassword)
		// 发送邮箱验证码
		authRouter.POST("/send-email-code", authHandler.SendEmailCode)
	}

	// 用户业务 Handler
	userHandler := user.NewHandler(cfg, coreBizClient)
	userRouter := router.Group("/user")
	userRouter.Use(middlewares.AuthRequired(coreBizClient))
	{
		// 获取用户头像
		userRouter.GET("/avatar", userHandler.GetAvatar)
		// 上传用户自定义头像
		userRouter.POST("/avatar", userHandler.UploadAvatar)
		// 删除用户自定义头像
		userRouter.DELETE("/avatar", userHandler.DeleteAvatar)
	}
}
