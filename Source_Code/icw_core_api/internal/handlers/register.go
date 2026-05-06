package handlers

import (
	"github.com/gin-gonic/gin"

	"icw_common/consts"
	"icw_common/utils"

	"icw_core_api/configs"
	"icw_core_api/internal/handlers/auth"
	"icw_core_api/internal/handlers/common"
	"icw_core_api/internal/handlers/project/assets"
	"icw_core_api/internal/handlers/project/core"
	"icw_core_api/internal/handlers/project/detection"
	"icw_core_api/internal/handlers/project/profile"
	"icw_core_api/internal/handlers/project/report"
	"icw_core_api/internal/handlers/project/review"
	"icw_core_api/internal/handlers/socket"
	"icw_core_api/internal/handlers/user"
	"icw_core_api/internal/middlewares"
	ws "icw_core_api/internal/socket"
	"icw_core_api/rpc/icw_core_biz"
)

// RegisterRoutes 注册 HTTP 路由
func RegisterRoutes(cfg configs.Config, coreBizClient *icw_core_biz.Client, socketHub *ws.Hub) *gin.Engine {
	gin.SetMode(gin.ReleaseMode)
	router := gin.New()
	router.Use(middlewares.RequestId(), middlewares.Logger(), gin.Recovery(), middlewares.CORS())
	routeDescriptions := make(map[string]string)

	// 创建 API Handler 的公共依赖集合
	handlerDeps := common.NewDeps(cfg, coreBizClient)
	registry(router, handlerDeps, coreBizClient, socketHub, routeDescriptions)

	utils.LogInfo(consts.LogScopeHTTP, consts.LogColorBoldGreen, "HTTP methods registered, waiting for requests:\n%s", common.FormatRoutesTable(router.Routes(), routeDescriptions))
	return router
}

// registry HTTP 路由注册表
func registry(router *gin.Engine, handlerDeps *common.Deps, coreBizClient *icw_core_biz.Client, socketHub *ws.Hub, routeDescriptions map[string]string) {
	// 登录鉴权 Handler
	authHandler := auth.NewHandler(handlerDeps)
	authRouter := router.Group("/auth")
	authRoutes := common.NewRouteGroup(authRouter, routeDescriptions)
	{
		authRoutes.POST("/login", "登录", authHandler.Login)
		authRoutes.POST("/logout", "登出", authHandler.Logout)
		protectedAuth := authRouter.Group("")
		protectedAuth.Use(middlewares.AuthRequired(coreBizClient))
		{
			common.NewRouteGroup(protectedAuth, routeDescriptions).GET("/me", "获取用户信息", authHandler.Me)
		}
		authRoutes.POST("/refresh", "刷新 Token", authHandler.Refresh)
		authRoutes.POST("/register", "注册", authHandler.Register)
		authRoutes.POST("/reset-password", "重置密码", authHandler.ResetPassword)
		authRoutes.POST("/send-email-code", "发送邮箱验证码", authHandler.SendEmailCode)
	}

	// 用户业务 Handler
	userHandler := user.NewHandler(handlerDeps)
	userRouter := router.Group("/user")
	userRoutes := common.NewRouteGroup(userRouter, routeDescriptions)
	userRouter.Use(middlewares.AuthRequired(coreBizClient))
	{
		userRoutes.DELETE("/avatar", "删除用户自定义头像", userHandler.DeleteAvatar)
		userRoutes.GET("/avatar", "获取用户头像", userHandler.GetAvatar)
		userRoutes.POST("/avatar", "上传用户自定义头像", userHandler.UploadAvatar)
	}

	// WebSocket Handler
	socketHandler := socket.NewHandler(handlerDeps, socketHub)
	socketRouter := router.Group("/socket")
	{
		// Socket 建连 Router
		socketSetupRouter := socketRouter.Group("/setup")
		{
			common.NewRouteGroup(socketSetupRouter, routeDescriptions).GET("/assets", "建立图像资产 WebSocket 连接", socketHandler.SetupAssetsWebSocket)
		}

		// Socket 票据 Router
		socketTicketRouter := socketRouter.Group("")
		socketTicketRouter.Use(middlewares.AuthRequired(coreBizClient))
		{
			common.NewRouteGroup(socketTicketRouter, routeDescriptions).POST("/ticket", "创建 WebSocket 连接票据", middlewares.ProjectAccessible(coreBizClient), socketHandler.CreateSocketTicket)
		}
	}

	// 项目流程 Router
	projectRouter := router.Group("/project")
	projectRouter.Use(middlewares.AuthRequired(coreBizClient))
	{
		// 校验项目访问权限
		projectAccessible := middlewares.ProjectAccessible(coreBizClient)
		// 校验项目基础信息阶段编辑权限
		projectProfileEditable := middlewares.ProjectProfileEditable(coreBizClient)
		// 校验图像资产构建阶段编辑权限
		projectAssetsEditable := middlewares.ProjectAssetsEditable(coreBizClient)
		// 校验 Agent 智能检测阶段编辑权限
		projectDetectionEditable := middlewares.ProjectDetectionEditable(coreBizClient)
		// 校验人工复核确认阶段编辑权限
		projectReviewEditable := middlewares.ProjectReviewEditable(coreBizClient)
		// 校验评估报告生成阶段编辑权限
		projectReportEditable := middlewares.ProjectReportEditable(coreBizClient)

		// 项目核心 Handler
		projectCoreHandler := core.NewHandler(handlerDeps)
		projectCoreRouter := projectRouter.Group("/core")
		projectCoreRoutes := common.NewRouteGroup(projectCoreRouter, routeDescriptions)
		{
			projectCoreRoutes.POST("/advance", "项目进度流转", projectAccessible, projectCoreHandler.AdvanceProject)
			projectCoreRoutes.POST("/create", "创建项目", projectCoreHandler.CreateProject)
			projectCoreRoutes.POST("/delete", "删除项目", projectAccessible, projectCoreHandler.DeleteProject)
			projectCoreRoutes.GET("/list", "获取项目列表", projectCoreHandler.ListProjects)
		}

		// 基础信息 Handler
		projectProfileHandler := profile.NewHandler(handlerDeps)
		projectProfileRouter := projectRouter.Group("/profile")
		projectProfileRoutes := common.NewRouteGroup(projectProfileRouter, routeDescriptions)
		{
			projectProfileRoutes.GET("/detail", "获取项目基础信息", projectAccessible, projectProfileHandler.GetProjectProfile)
			projectProfileRoutes.GET("/thumbnail", "获取项目缩略图", projectAccessible, projectProfileHandler.GetProjectThumbnail)
			projectProfileRoutes.POST("/thumbnail", "上传项目缩略图", projectProfileEditable, projectProfileHandler.UploadProjectThumbnail)
			projectProfileRoutes.DELETE("/thumbnail", "删除项目缩略图", projectProfileEditable, projectProfileHandler.DeleteProjectThumbnail)
			projectProfileRoutes.POST("/update", "更新项目基础信息", projectProfileEditable, projectProfileHandler.UpdateProjectProfile)
		}

		// 图像资产 Handler
		projectAssetsHandler := assets.NewHandler(handlerDeps)
		projectAssetsRouter := projectRouter.Group("/assets")
		projectAssetsRoutes := common.NewRouteGroup(projectAssetsRouter, routeDescriptions)
		{
			projectAssetsRoutes.GET("/list", "获取项目图像列表", projectAccessible, projectAssetsHandler.GetProjectAssets)

			// 图像组 Router
			projectAssetsGroupRouter := projectAssetsRouter.Group("/group")
			projectAssetsGroupRoutes := common.NewRouteGroup(projectAssetsGroupRouter, routeDescriptions)
			{
				projectAssetsGroupRoutes.POST("/create", "创建图像组", projectAssetsEditable, projectAssetsHandler.CreateProjectGroup)
				projectAssetsGroupRoutes.POST("/delete", "删除图像组", projectAssetsEditable, projectAssetsHandler.DeleteProjectGroup)
				projectAssetsGroupRoutes.POST("/move", "移动图像组", projectAssetsEditable, projectAssetsHandler.MoveProjectGroup)
				projectAssetsGroupRoutes.POST("/update", "更新图像组", projectAssetsEditable, projectAssetsHandler.UpdateProjectGroup)
			}

			// 图像 Router
			projectAssetsImageRouter := projectAssetsRouter.Group("/image")
			projectAssetsImageRoutes := common.NewRouteGroup(projectAssetsImageRouter, routeDescriptions)
			{
				projectAssetsImageRoutes.POST("/delete", "删除图像", projectAssetsEditable, projectAssetsHandler.DeleteProjectImage)
				projectAssetsImageRoutes.GET("/original", "获取原图", projectAccessible, projectAssetsHandler.GetProjectImageOriginal)
				projectAssetsImageRoutes.POST("/move", "移动图像", projectAssetsEditable, projectAssetsHandler.MoveProjectImage)
				projectAssetsImageRoutes.POST("/report", "上报图像", projectAssetsEditable, projectAssetsHandler.ReportProjectImage)
				projectAssetsImageRoutes.POST("/upload", "上传图像", projectAssetsEditable, projectAssetsHandler.UploadProjectImage)
			}
		}

		// 智能检测 Handler
		projectDetectionHandler := detection.NewHandler(handlerDeps)
		projectDetectionRouter := projectRouter.Group("/detection")
		{
			if projectDetectionHandler == nil || projectDetectionRouter == nil || projectDetectionEditable == nil {
				// TODO: Prevent errors on unused variables
			}
		}

		// 人工复核 Handler
		projectReviewHandler := review.NewHandler(handlerDeps)
		projectReviewRouter := projectRouter.Group("/review")
		{
			if projectReviewHandler == nil || projectReviewRouter == nil || projectReviewEditable == nil {
				// TODO: Prevent errors on unused variables
			}
		}

		// 评估报告 Handler
		projectReportHandler := report.NewHandler(handlerDeps)
		projectReportRouter := projectRouter.Group("/report")
		{
			if projectReportHandler == nil || projectReportRouter == nil || projectReportEditable == nil {
				// TODO: Prevent errors on unused variables
			}
		}
	}
}
