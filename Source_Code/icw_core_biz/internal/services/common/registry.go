package common

import (
	"runtime"
	"sync"
)

var rpcHandlerRegistry sync.Map

// RegisterRPCHandler 注册 RPC 执行函数
func RegisterRPCHandler(handler, method string) {
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
