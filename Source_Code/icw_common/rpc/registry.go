package rpc

import (
	"context"
	"errors"
	"fmt"
	"reflect"
	"runtime"
	"sort"
	"strings"

	"google.golang.org/grpc"

	"icw_common/consts"
	"icw_common/utils"
)

var (
	// contextType 校验 RPC 方法第一个参数类型必须是 context.Context
	contextType = reflect.TypeOf((*context.Context)(nil)).Elem()
	// errorType 校验 RPC 方法的返回值类型必须是 error
	errorType = reflect.TypeOf((*error)(nil)).Elem()
)

// serviceMeta RPC 服务元数据
type serviceMeta struct {
	Name        string
	Description string
	Service     interface{}
	Register    func(grpc.ServiceRegistrar, interface{})
	Methods     []methodMeta
}

// methodMeta RPC 方法元数据
type methodMeta struct {
	Name        string
	Description string
}

// registeredMethodMeta 已注册 RPC 方法元数据
type registeredMethodMeta struct {
	serviceName       string
	methodName        string
	methodDescription string
	handler           string
}

// RegisterServices 注册 gRPC 服务并输出 RPC 注册表
func RegisterServices(grpcServer *grpc.Server, metas []serviceMeta) {
	registeredMethods := make([]registeredMethodMeta, 0)
	for _, meta := range metas {
		methods, err := resolveMethods(meta)
		if err != nil {
			utils.LogFatal(consts.LogScopeRPC, "Failed to register RPC service %s: %v", meta.Name, err)
		}
		if meta.Register == nil {
			utils.LogFatal(consts.LogScopeRPC, "Failed to register RPC service %s: register function is nil", meta.Name)
		}
		meta.Register(grpcServer, meta.Service)
		registeredMethods = append(registeredMethods, methods...)
	}
	utils.LogInfo(consts.LogScopeRPC, consts.LogColorBoldGreen, "RPC methods registered, waiting for requests:\n%s", formatRegistryTable(registeredMethods))
}

// resolveMethods 校验并解析 RPC 方法
func resolveMethods(meta serviceMeta) ([]registeredMethodMeta, error) {
	serviceType := reflect.TypeOf(meta.Service)
	if serviceType == nil {
		return nil, errors.New("service is nil")
	}

	availableMethods := rpcMethodMap(serviceType)
	methodDescriptions := make(map[string]string, len(meta.Methods))
	for _, method := range meta.Methods {
		if _, exists := methodDescriptions[method.Name]; exists {
			return nil, fmt.Errorf("method %s is duplicated", method.Name)
		}
		if _, ok := availableMethods[method.Name]; !ok {
			return nil, fmt.Errorf("method %s is not a suitable rpc method", method.Name)
		}
		methodDescriptions[method.Name] = method.Description
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

	methods := make([]registeredMethodMeta, 0, len(meta.Methods))
	for _, method := range meta.Methods {
		reflectedMethod := availableMethods[method.Name]
		handler := handlerName(reflectedMethod)
		methods = append(methods, registeredMethodMeta{
			serviceName:       meta.Name,
			methodName:        method.Name,
			methodDescription: method.Description,
			handler:           handler,
		})
	}

	return methods, nil
}

// rpcMethodMap 获取符合 gRPC 入口规范的方法
func rpcMethodMap(serviceType reflect.Type) map[string]reflect.Method {
	methods := make(map[string]reflect.Method)
	for index := 0; index < serviceType.NumMethod(); index++ {
		method := serviceType.Method(index)
		methodType := method.Type
		if !(method.PkgPath == "" && methodType.NumIn() == 3 &&
			methodType.In(1) == contextType && methodType.In(2).Kind() == reflect.Ptr &&
			methodType.NumOut() == 2 && methodType.Out(0).Kind() == reflect.Ptr && methodType.Out(1) == errorType) {
			continue
		}
		if strings.HasPrefix(handlerName(method), "icw_common/") {
			continue
		}
		methods[method.Name] = method
	}
	return methods
}

// handlerName 获取 RPC 方法执行函数名
func handlerName(method reflect.Method) string {
	fn := runtime.FuncForPC(method.Func.Pointer())
	if fn == nil {
		return method.Name
	}
	name := fn.Name()
	const marker = "/internal/services/"
	if index := strings.Index(name, marker); index >= 0 {
		return name[index+len(marker):]
	}
	return name
}

// formatRegistryTable 格式化 RPC 注册表
func formatRegistryTable(methods []registeredMethodMeta) string {
	serviceMethods := make([]string, 0, len(methods))
	methodDescriptions := make([]string, 0, len(methods))
	handlers := make([]string, 0, len(methods))
	for _, method := range methods {
		serviceMethods = append(serviceMethods, method.serviceName+"."+method.methodName)
		methodDescriptions = append(methodDescriptions, method.methodDescription)
		handlers = append(handlers, method.handler)
	}
	return utils.FormatTable([]*utils.TableColumn{
		{
			Header: "service.method",
			Values: serviceMethods,
		},
		{
			Header: "description",
			Values: methodDescriptions,
		},
		{
			Header: "handler",
			Values: handlers,
		},
	})
}
