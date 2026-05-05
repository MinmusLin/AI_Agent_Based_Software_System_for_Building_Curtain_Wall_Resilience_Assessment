package services

import (
	"context"
	"fmt"
	"net/rpc"
	"reflect"
	"runtime"
	"sort"

	"icw_core_biz/internal/services/auth"
	"icw_core_biz/internal/services/common"
	"icw_core_biz/internal/services/project/assets"
	"icw_core_biz/internal/services/project/core"
	"icw_core_biz/internal/services/project/detection"
	"icw_core_biz/internal/services/project/profile"
	"icw_core_biz/internal/services/project/report"
	"icw_core_biz/internal/services/project/review"
	"icw_core_biz/internal/services/socket"
	"icw_core_biz/internal/services/user"
	"icw_core_biz/utils"
)

var (
	// errorType 校验 RPC 方法的返回值类型必须是 error
	errorType = reflect.TypeOf((*error)(nil)).Elem()
)

// rpcServiceMeta RPC 服务元数据
type rpcServiceMeta struct {
	name        string
	description string
	service     interface{}
	methods     []rpcMethodMeta
}

// rpcMethodMeta RPC 方法元数据
type rpcMethodMeta struct {
	name        string
	description string
}

// registeredRPCMethod 已注册 RPC 方法元数据
type registeredRPCMethod struct {
	serviceName        string
	serviceDescription string
	methodName         string
	methodDescription  string
	handler            string
}

// RegisterRPCServices 注册 RPC 服务
func RegisterRPCServices(ctx context.Context, serviceDeps *common.Deps) {
	registeredMethods := make([]registeredRPCMethod, 0)
	for _, meta := range rpcRegistry(ctx, serviceDeps) {
		methods, err := registerRPCService(meta)
		if err != nil {
			common.RpcFatal("Failed to register RPC service %s: %v", meta.name, err)
		}
		registeredMethods = append(registeredMethods, methods...)
	}
	common.RpcInfo("RPC methods registered, waiting for requests:\n%s", formatRPCRegistryTable(registeredMethods))
}

// rpcRegistry RPC 服务注册表
func rpcRegistry(ctx context.Context, serviceDeps *common.Deps) []rpcServiceMeta {
	return []rpcServiceMeta{
		{
			name:        "SocketService",
			description: "WebSocket 连接票据服务",
			service:     socket.NewService(ctx, serviceDeps),
			methods: []rpcMethodMeta{
				{name: "CreateSocketTicket", description: "创建 WebSocket 连接票据"},
				{name: "ValidateSocketTicket", description: "校验 WebSocket 连接票据"},
			},
		},
		{
			name:        "UserService",
			description: "用户业务服务",
			service:     user.NewService(ctx, serviceDeps),
			methods: []rpcMethodMeta{
				{name: "GetAvatar", description: "获取用户头像"},
				{name: "UploadAvatar", description: "上传用户自定义头像"},
				{name: "DeleteAvatar", description: "删除用户自定义头像"},
			},
		},
		{
			name:        "AuthService",
			description: "登录鉴权服务",
			service:     auth.NewService(ctx, serviceDeps),
			methods: []rpcMethodMeta{
				{name: "Login", description: "登录"},
				{name: "Logout", description: "登出"},
				{name: "Me", description: "获取用户信息"},
				{name: "Refresh", description: "刷新 Token"},
				{name: "Register", description: "注册"},
				{name: "ResetPassword", description: "重置密码"},
				{name: "SendEmailCode", description: "发送邮箱验证码"},
			},
		},
		{
			name:        "ProjectCoreService",
			description: "项目核心服务",
			service:     core.NewService(ctx, serviceDeps),
			methods: []rpcMethodMeta{
				{name: "AdvanceProject", description: "项目进度流转"},
				{name: "CreateProject", description: "创建项目"},
				{name: "DeleteProject", description: "删除项目"},
				{name: "ListProjects", description: "获取项目列表"},
				{name: "CheckProjectAccess", description: "校验项目访问权限"},
			},
		},
		{
			name:        "ProjectProfileService",
			description: "基础信息服务",
			service:     profile.NewService(ctx, serviceDeps),
			methods: []rpcMethodMeta{
				{name: "GetProjectProfile", description: "获取项目基础信息"},
				{name: "GetProjectThumbnail", description: "获取项目缩略图"},
				{name: "UploadProjectThumbnail", description: "上传项目缩略图"},
				{name: "DeleteProjectThumbnail", description: "删除项目缩略图"},
				{name: "UpdateProjectProfile", description: "更新项目基础信息"},
			},
		},
		{
			name:        "ProjectAssetsService",
			description: "图像资产服务",
			service:     assets.NewService(ctx, serviceDeps),
			methods: []rpcMethodMeta{
				{name: "GetProjectAssets", description: "获取项目图像列表"},
				{name: "CreateProjectGroup", description: "创建图像组"},
				{name: "DeleteProjectGroup", description: "删除图像组"},
				{name: "MoveProjectGroup", description: "移动图像组"},
				{name: "UpdateProjectGroup", description: "更新图像组"},
				{name: "DeleteProjectImage", description: "删除图像"},
				{name: "GetProjectImageOriginal", description: "获取原图"},
				{name: "MoveProjectImage", description: "移动图像"},
				{name: "ReportProjectImage", description: "上报图像"},
				{name: "UploadProjectImage", description: "上传图像"},
			},
		},
		{
			name:        "ProjectDetectionService",
			description: "智能检测服务",
			service:     detection.NewService(ctx, serviceDeps),
			methods: []rpcMethodMeta{
				{name: "Ping", description: "智能检测服务探活"},
				{name: "ReportClassificationResult", description: "上报图像检测分类结果"},
				{name: "ReportReasoningResult", description: "上报图像检测推理结果"},
				{name: "ReportSummaryResult", description: "上报图像检测总结结果"},
			},
		},
		{
			name:        "ProjectReviewService",
			description: "人工复核服务",
			service:     review.NewService(ctx, serviceDeps),
			methods: []rpcMethodMeta{
				{name: "Ping", description: "人工复核服务探活"},
			},
		},
		{
			name:        "ProjectReportService",
			description: "评估报告服务",
			service:     report.NewService(ctx, serviceDeps),
			methods: []rpcMethodMeta{
				{name: "Ping", description: "评估报告服务探活"},
			},
		},
	}
}

// registerRPCService 注册单个 RPC 服务并返回注册表行数据
func registerRPCService(meta rpcServiceMeta) ([]registeredRPCMethod, error) {
	methods, err := resolveRPCMethods(meta)
	if err != nil {
		return nil, err
	}
	if err := rpc.RegisterName(meta.name, meta.service); err != nil {
		return nil, err
	}
	return methods, nil
}

// resolveRPCMethods 校验并解析 RPC 方法
func resolveRPCMethods(meta rpcServiceMeta) ([]registeredRPCMethod, error) {
	serviceType := reflect.TypeOf(meta.service)
	if serviceType == nil {
		return nil, fmt.Errorf("service is nil")
	}

	availableMethods := rpcMethodMap(serviceType)
	methodDescriptions := make(map[string]string, len(meta.methods))
	for _, method := range meta.methods {
		if _, exists := methodDescriptions[method.name]; exists {
			return nil, fmt.Errorf("method %s is duplicated", method.name)
		}
		if _, ok := availableMethods[method.name]; !ok {
			return nil, fmt.Errorf("method %s is not a suitable RPC method", method.name)
		}
		methodDescriptions[method.name] = method.description
	}

	availableNames := make([]string, 0, len(availableMethods))
	for methodName := range availableMethods {
		if _, ok := methodDescriptions[methodName]; ok {
			continue
		}
		availableNames = append(availableNames, methodName)
	}
	sort.Strings(availableNames)
	if len(availableNames) > 0 {
		return nil, fmt.Errorf("method description is missing: %v", availableNames)
	}

	rows := make([]registeredRPCMethod, 0, len(meta.methods))
	for _, method := range meta.methods {
		reflectedMethod := availableMethods[method.name]
		handler := handlerName(reflectedMethod)
		common.RegisterRPCHandler(handler, meta.name+"."+method.name)
		rows = append(rows, registeredRPCMethod{
			serviceName:        meta.name,
			serviceDescription: meta.description,
			methodName:         method.name,
			methodDescription:  method.description,
			handler:            handler,
		})
	}
	return rows, nil
}

// rpcMethodMap 获取 service 上符合 net/rpc 规范的方法
func rpcMethodMap(serviceType reflect.Type) map[string]reflect.Method {
	methods := make(map[string]reflect.Method)
	for index := 0; index < serviceType.NumMethod(); index++ {
		method := serviceType.Method(index)
		if !isSuitableRPCMethod(method) {
			continue
		}
		methods[method.Name] = method
	}
	return methods
}

// isSuitableRPCMethod 判断方法是否符合 net/rpc 方法签名
func isSuitableRPCMethod(method reflect.Method) bool {
	methodType := method.Type
	return method.PkgPath == "" &&
		methodType.NumIn() == 3 &&
		methodType.In(1).Kind() == reflect.Ptr &&
		methodType.In(2).Kind() == reflect.Ptr &&
		methodType.NumOut() == 1 &&
		methodType.Out(0) == errorType
}

// handlerName 获取 RPC 方法真实执行函数名
func handlerName(method reflect.Method) string {
	fn := runtime.FuncForPC(method.Func.Pointer())
	if fn == nil {
		return method.Name
	}
	return fn.Name()
}

// formatRPCRegistryTable 格式化 RPC 注册表
func formatRPCRegistryTable(methods []registeredRPCMethod) string {
	serviceNames := make([]string, 0, len(methods))
	serviceDescriptions := make([]string, 0, len(methods))
	methodNames := make([]string, 0, len(methods))
	methodDescriptions := make([]string, 0, len(methods))
	handlers := make([]string, 0, len(methods))
	for _, method := range methods {
		serviceNames = append(serviceNames, method.serviceName)
		serviceDescriptions = append(serviceDescriptions, method.serviceDescription)
		methodNames = append(methodNames, method.methodName)
		methodDescriptions = append(methodDescriptions, method.methodDescription)
		handlers = append(handlers, method.handler)
	}
	return utils.FormatTable([]utils.TableColumn{
		{
			Header: "serviceName",
			Values: serviceNames,
		},
		{
			Header: "serviceDescription",
			Values: serviceDescriptions,
		},
		{
			Header: "methodName",
			Values: methodNames,
		},
		{
			Header: "methodDescription",
			Values: methodDescriptions,
		},
		{
			Header: "handler",
			Values: handlers,
		},
	})
}
