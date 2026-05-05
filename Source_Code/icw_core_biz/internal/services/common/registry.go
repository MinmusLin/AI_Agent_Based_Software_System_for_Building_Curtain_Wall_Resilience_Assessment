package common

import (
	"errors"
	"fmt"
	"net/rpc"
	"reflect"
	"runtime"
	"sort"
	"sync"

	"icw_core_biz/utils"
)

var (
	// rpcHandlerRegistry RPC Handler 全局映射表
	rpcHandlerRegistry sync.Map
	// errorType 校验 RPC 方法的返回值类型必须是 error
	errorType = reflect.TypeOf((*error)(nil)).Elem()
)

// RPCServiceMeta RPC 服务元数据
type RPCServiceMeta struct {
	Name        string
	Description string
	Service     interface{}
	Methods     []RPCMethodMeta
}

// RPCMethodMeta RPC 方法元数据
type RPCMethodMeta struct {
	Name        string
	Description string
}

// RegisteredRPCMethodMeta 已注册 RPC 方法元数据
type RegisteredRPCMethodMeta struct {
	serviceName        string
	serviceDescription string
	methodName         string
	methodDescription  string
	handler            string
}

// RegisterRPCService 注册单个 RPC 服务并返回 RPC 方法元数据
func RegisterRPCService(meta RPCServiceMeta) ([]RegisteredRPCMethodMeta, error) {
	methods, err := resolveRPCMethods(meta)
	if err != nil {
		return nil, err
	}
	if err := rpc.RegisterName(meta.Name, meta.Service); err != nil {
		return nil, err
	}
	return methods, nil
}

// resolveRPCMethods 校验并解析 RPC 方法
func resolveRPCMethods(meta RPCServiceMeta) ([]RegisteredRPCMethodMeta, error) {
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
			return nil, fmt.Errorf("method %s is not a suitable RPC method", method.Name)
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

	rows := make([]RegisteredRPCMethodMeta, 0, len(meta.Methods))
	for _, method := range meta.Methods {
		reflectedMethod := availableMethods[method.Name]
		handler := handlerName(reflectedMethod)
		registerRPCHandler(handler, meta.Name+"."+method.Name)
		rows = append(rows, RegisteredRPCMethodMeta{
			serviceName:        meta.Name,
			serviceDescription: meta.Description,
			methodName:         method.Name,
			methodDescription:  method.Description,
			handler:            handler,
		})
	}
	return rows, nil
}

// rpcMethodMap 获取符合 net/rpc 规范的方法
func rpcMethodMap(serviceType reflect.Type) map[string]reflect.Method {
	methods := make(map[string]reflect.Method)
	for index := 0; index < serviceType.NumMethod(); index++ {
		method := serviceType.Method(index)
		methodType := method.Type
		if !(method.PkgPath == "" && methodType.NumIn() == 3 &&
			methodType.In(1).Kind() == reflect.Ptr && methodType.In(2).Kind() == reflect.Ptr &&
			methodType.NumOut() == 1 && methodType.Out(0) == errorType) {
			continue
		}
		methods[method.Name] = method
	}
	return methods
}

// registerRPCHandler 注册 RPC 执行函数
func registerRPCHandler(handler, method string) {
	if handler == "" || method == "" {
		return
	}
	rpcHandlerRegistry.Store(handler, method)
}

// rpcMethodFromPC 根据调用方 PC 获取 RPC 方法名
func rpcMethodFromPC(pc uintptr) string {
	fn := runtime.FuncForPC(pc)
	if fn == nil {
		return "unknown"
	}
	handler := fn.Name()
	if method, ok := rpcHandlerRegistry.Load(handler); ok {
		if methodName, ok := method.(string); ok && methodName != "" {
			return methodName
		}
	}
	return handler
}

// handlerName 获取 RPC 方法执行函数名
func handlerName(method reflect.Method) string {
	fn := runtime.FuncForPC(method.Func.Pointer())
	if fn == nil {
		return method.Name
	}
	return fn.Name()
}

// FormatRegistryTable 格式化 RPC 注册表
func FormatRegistryTable(methods []RegisteredRPCMethodMeta) string {
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
			Header: "service",
			Values: serviceNames,
		},
		{
			Header: "description",
			Values: serviceDescriptions,
		},
		{
			Header: "method",
			Values: methodNames,
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
