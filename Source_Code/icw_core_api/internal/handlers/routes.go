package handlers

import (
	"github.com/gin-gonic/gin"

	"icw_core_api/configs"
	"icw_core_api/internal/handlers/auth"
	"icw_core_api/internal/handlers/common"
	"icw_core_api/internal/handlers/project/assets"
	"icw_core_api/internal/handlers/project/core"
	"icw_core_api/internal/handlers/project/detection"
	"icw_core_api/internal/handlers/project/profile"
	"icw_core_api/internal/handlers/project/report"
	"icw_core_api/internal/handlers/project/review"
	"icw_core_api/internal/handlers/user"
	"icw_core_api/internal/middlewares"
)

// RegisterRoutes 注册路由
func RegisterRoutes(router *gin.Engine, cfg configs.Config, coreBizClient *common.RPCClient) {
	// 创建 API Handler 的公共依赖集合
	handlerDeps := common.NewDeps(cfg, coreBizClient)

	// 登录鉴权 Handler
	authHandler := auth.NewHandler(handlerDeps)
	authRouter := router.Group("/auth")
	{
		protectedAuth := authRouter.Group("")
		// 登录
		authRouter.POST("/login", authHandler.Login)
		// 登出
		authRouter.POST("/logout", authHandler.Logout)
		protectedAuth.Use(middlewares.AuthRequired(coreBizClient.Raw()))
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
	userHandler := user.NewHandler(handlerDeps)
	userRouter := router.Group("/user")
	userRouter.Use(middlewares.AuthRequired(coreBizClient.Raw()))
	{
		// 获取用户头像
		userRouter.GET("/avatar", userHandler.GetAvatar)
		// 上传用户自定义头像
		userRouter.POST("/avatar", userHandler.UploadAvatar)
		// 删除用户自定义头像
		userRouter.DELETE("/avatar", userHandler.DeleteAvatar)
	}

	// 项目流程 Router
	projectRouter := router.Group("/project")
	projectRouter.Use(middlewares.AuthRequired(coreBizClient.Raw()))
	{
		// 项目核心 Handler
		projectCoreHandler := core.NewHandler(handlerDeps)
		projectCoreRouter := projectRouter.Group("/core")
		{
			// 项目进度流转
			projectCoreRouter.POST("/advance", projectCoreHandler.AdvanceProject)
			// 创建项目
			projectCoreRouter.POST("/create", projectCoreHandler.CreateProject)
			// 删除项目
			projectCoreRouter.POST("/delete", projectCoreHandler.DeleteProject)
			// 获取项目列表
			projectCoreRouter.GET("/list", projectCoreHandler.ListProjects)
		}

		// 基础信息 Handler
		projectProfileHandler := profile.NewHandler(handlerDeps)
		projectProfileRouter := projectRouter.Group("/profile")
		{
			// 获取项目基础信息
			projectProfileRouter.GET("/detail", projectProfileHandler.GetProjectProfile)
			// 获取项目缩略图
			projectProfileRouter.GET("/thumbnail", projectProfileHandler.GetProjectThumbnail)
			// 上传项目缩略图
			projectProfileRouter.POST("/thumbnail", projectProfileHandler.UploadProjectThumbnail)
			// 删除项目缩略图
			projectProfileRouter.DELETE("/thumbnail", projectProfileHandler.DeleteProjectThumbnail)
			// 更新项目基础信息
			projectProfileRouter.POST("/update", projectProfileHandler.UpdateProjectProfile)
		}

		// 图像资产 Handler
		projectAssetsHandler := assets.NewHandler(handlerDeps)
		projectAssetsRouter := projectRouter.Group("/assets")
		{
			// 获取项目图像列表
			projectAssetsRouter.GET("/list", projectAssetsHandler.GetProjectAssets)

			// 图像组 Router
			projectAssetsGroupRouter := projectAssetsRouter.Group("/group")
			{
				// 创建图像组
				projectAssetsGroupRouter.POST("/create", projectAssetsHandler.CreateProjectGroup)
				// 删除图像组
				projectAssetsGroupRouter.POST("/delete", projectAssetsHandler.DeleteProjectGroup)
				// 移动图像组
				projectAssetsGroupRouter.POST("/move", projectAssetsHandler.MoveProjectGroup)
				// 更新图像组
				projectAssetsGroupRouter.POST("/update", projectAssetsHandler.UpdateProjectGroup)
			}

			// 图像 Router
			projectAssetsImageRouter := projectAssetsRouter.Group("/image")
			{
				// 删除图像
				projectAssetsImageRouter.POST("/delete", projectAssetsHandler.DeleteProjectImage)
				// 获取原图
				projectAssetsImageRouter.GET("/original", projectAssetsHandler.GetProjectImageOriginal)
				// 移动图像
				projectAssetsImageRouter.POST("/move", projectAssetsHandler.MoveProjectImage)
				// 上报图像
				projectAssetsImageRouter.POST("/report", projectAssetsHandler.ReportProjectImage)
				// 上传图像
				projectAssetsImageRouter.POST("/upload", projectAssetsHandler.UploadProjectImage)
			}
		}

		// 智能检测 Handler
		projectDetectionHandler := detection.NewHandler(handlerDeps)
		projectDetectionRouter := projectRouter.Group("/detection")
		{
			if projectDetectionHandler == nil || projectDetectionRouter == nil {
				// TODO: Prevent errors on unused variables
			}
		}

		// 人工复核 Handler
		projectReviewHandler := review.NewHandler(handlerDeps)
		projectReviewRouter := projectRouter.Group("/review")
		{
			if projectReviewHandler == nil || projectReviewRouter == nil {
				// TODO: Prevent errors on unused variables
			}
		}

		// 评估报告 Handler
		projectReportHandler := report.NewHandler(handlerDeps)
		projectReportRouter := projectRouter.Group("/report")
		{
			if projectReportHandler == nil || projectReportRouter == nil {
				// TODO: Prevent errors on unused variables
			}
		}
	}
}
