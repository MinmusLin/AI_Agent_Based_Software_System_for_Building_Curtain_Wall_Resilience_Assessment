package handlers

import (
	"fmt"
	"net/http"
	"strings"

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
	"icw_core_api/internal/handlers/socket"
	"icw_core_api/internal/handlers/user"
	"icw_core_api/internal/middlewares"
	ws "icw_core_api/internal/socket"
	"icw_core_biz/consts"
	"icw_core_biz/utils"
)

// routeGroup 封装 Gin 路由注册器
type routeGroup struct {
	router       *gin.RouterGroup
	descriptions map[string]string
}

// newRouteGroup 创建路由注册器
func newRouteGroup(router *gin.RouterGroup, descriptions map[string]string) routeGroup {
	return routeGroup{
		router:       router,
		descriptions: descriptions,
	}
}

// handle 注册 Gin 路由
func (r routeGroup) handle(method, path, description string, handlers ...gin.HandlerFunc) {
	if r.router == nil {
		return
	}
	r.router.Handle(method, path, handlers...)
	if r.descriptions != nil {
		r.descriptions[routeKey(method, joinPath(r.router.BasePath(), path))] = description
	}
}

// GET 注册 GET 路由
func (r routeGroup) GET(path, description string, handlers ...gin.HandlerFunc) {
	r.handle(http.MethodGet, path, description, handlers...)
}

// POST 注册 POST 路由
func (r routeGroup) POST(path, description string, handlers ...gin.HandlerFunc) {
	r.handle(http.MethodPost, path, description, handlers...)
}

// DELETE 注册 DELETE 路由
func (r routeGroup) DELETE(path, description string, handlers ...gin.HandlerFunc) {
	r.handle(http.MethodDelete, path, description, handlers...)
}

// RegisterRoutes 注册路由
func RegisterRoutes(cfg configs.Config, coreBizClient *common.RPCClient, socketHub *ws.Hub) *gin.Engine {
	gin.SetMode(gin.ReleaseMode)
	router := gin.New()
	router.Use(middlewares.RequestId(), middlewares.Logger(), gin.Recovery(), middlewares.CORS())
	routeDescriptions := make(map[string]string)

	// 创建 API Handler 的公共依赖集合
	handlerDeps := common.NewDeps(cfg, coreBizClient)

	// 登录鉴权 Handler
	authHandler := auth.NewHandler(handlerDeps)
	authRouter := router.Group("/auth")
	authRoutes := newRouteGroup(authRouter, routeDescriptions)
	{
		protectedAuth := authRouter.Group("")
		authRoutes.POST("/login", "登录", authHandler.Login)
		authRoutes.POST("/logout", "登出", authHandler.Logout)
		protectedAuth.Use(middlewares.AuthRequired(coreBizClient))
		{
			newRouteGroup(protectedAuth, routeDescriptions).GET("/me", "获取用户信息", authHandler.Me)
		}
		authRoutes.POST("/refresh", "刷新 Token", authHandler.Refresh)
		authRoutes.POST("/register", "注册", authHandler.Register)
		authRoutes.POST("/reset-password", "重置密码", authHandler.ResetPassword)
		authRoutes.POST("/send-email-code", "发送邮箱验证码", authHandler.SendEmailCode)
	}

	// 用户业务 Handler
	userHandler := user.NewHandler(handlerDeps)
	userRouter := router.Group("/user")
	userRoutes := newRouteGroup(userRouter, routeDescriptions)
	userRouter.Use(middlewares.AuthRequired(coreBizClient))
	{
		userRoutes.GET("/avatar", "获取用户头像", userHandler.GetAvatar)
		userRoutes.POST("/avatar", "上传用户自定义头像", userHandler.UploadAvatar)
		userRoutes.DELETE("/avatar", "删除用户自定义头像", userHandler.DeleteAvatar)
	}

	// WebSocket Handler
	socketHandler := socket.NewHandler(handlerDeps, socketHub)
	socketRouter := router.Group("/socket")
	{
		// Socket 建连 Router
		socketSetupRouter := socketRouter.Group("/setup")
		{
			newRouteGroup(socketSetupRouter, routeDescriptions).GET("/assets", "建立图像资产 WebSocket 连接", socketHandler.SetupAssetsWebSocket)
		}

		// Socket 票据 Router
		socketTicketRouter := socketRouter.Group("")
		socketTicketRouter.Use(middlewares.AuthRequired(coreBizClient))
		{
			newRouteGroup(socketTicketRouter, routeDescriptions).POST("/ticket", "创建 WebSocket 连接票据", middlewares.ProjectAccessible(coreBizClient), socketHandler.CreateSocketTicket)
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
		projectCoreRoutes := newRouteGroup(projectCoreRouter, routeDescriptions)
		{
			projectCoreRoutes.POST("/advance", "项目进度流转", projectAccessible, projectCoreHandler.AdvanceProject)
			projectCoreRoutes.POST("/create", "创建项目", projectCoreHandler.CreateProject)
			projectCoreRoutes.POST("/delete", "删除项目", projectAccessible, projectCoreHandler.DeleteProject)
			projectCoreRoutes.GET("/list", "获取项目列表", projectCoreHandler.ListProjects)
		}

		// 基础信息 Handler
		projectProfileHandler := profile.NewHandler(handlerDeps)
		projectProfileRouter := projectRouter.Group("/profile")
		projectProfileRoutes := newRouteGroup(projectProfileRouter, routeDescriptions)
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
		projectAssetsRoutes := newRouteGroup(projectAssetsRouter, routeDescriptions)
		{
			projectAssetsRoutes.GET("/list", "获取项目图像列表", projectAccessible, projectAssetsHandler.GetProjectAssets)

			// 图像组 Router
			projectAssetsGroupRouter := projectAssetsRouter.Group("/group")
			projectAssetsGroupRoutes := newRouteGroup(projectAssetsGroupRouter, routeDescriptions)
			{
				projectAssetsGroupRoutes.POST("/create", "创建图像组", projectAssetsEditable, projectAssetsHandler.CreateProjectGroup)
				projectAssetsGroupRoutes.POST("/delete", "删除图像组", projectAssetsEditable, projectAssetsHandler.DeleteProjectGroup)
				projectAssetsGroupRoutes.POST("/move", "移动图像组", projectAssetsEditable, projectAssetsHandler.MoveProjectGroup)
				projectAssetsGroupRoutes.POST("/update", "更新图像组", projectAssetsEditable, projectAssetsHandler.UpdateProjectGroup)
			}

			// 图像 Router
			projectAssetsImageRouter := projectAssetsRouter.Group("/image")
			projectAssetsImageRoutes := newRouteGroup(projectAssetsImageRouter, routeDescriptions)
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

	utils.LogInfo(consts.LogScopeHTTP, consts.LogColorBoldGreen, "HTTP routes registered, waiting for requests:\n%s", formatRoutesTable(router.Routes(), routeDescriptions))
	return router
}

// routeKey 生成路由描述 Key
func routeKey(method, path string) string {
	return fmt.Sprintf("%s:%s", method, path)
}

// joinPath 拼接路由路径
func joinPath(basePath, path string) string {
	if basePath == "/" {
		basePath = ""
	}
	fullPath := strings.TrimRight(basePath, "/") + "/" + strings.TrimLeft(path, "/")
	if fullPath == "" {
		return "/"
	}
	return fullPath
}

// formatRoutesTable 将 Gin 路由表格式化为表格
func formatRoutesTable(routes gin.RoutesInfo, descriptions map[string]string) string {
	methodValues := make([]string, 0, len(routes))
	pathValues := make([]string, 0, len(routes))
	descriptionValues := make([]string, 0, len(routes))
	handlerValues := make([]string, 0, len(routes))
	for _, route := range routes {
		methodValues = append(methodValues, route.Method)
		pathValues = append(pathValues, route.Path)
		descriptionValues = append(descriptionValues, descriptions[routeKey(route.Method, route.Path)])
		handlerValues = append(handlerValues, strings.TrimSuffix(strings.TrimPrefix(route.Handler, "icw_core_api/internal/handlers/"), "-fm"))
	}
	return utils.FormatTable([]utils.TableColumn{
		{
			Header: "method",
			Values: methodValues,
		},
		{
			Header: "path",
			Values: pathValues,
		},
		{
			Header: "description",
			Values: descriptionValues,
		},
		{
			Header: "handler",
			Values: handlerValues,
		},
	})
}
