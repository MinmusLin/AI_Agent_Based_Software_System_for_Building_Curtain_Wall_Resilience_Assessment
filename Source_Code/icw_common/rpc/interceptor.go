package rpc

import (
	"context"
	"strings"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/status"

	"icw_common/consts"
	"icw_common/utils"
)

// UnaryServerInterceptor gRPC 服务端日志拦截器
func UnaryServerInterceptor(scope string) grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		start := time.Now()
		resp, err := handler(ctx, req)
		logGRPCServer(ctx, scope, req, resp, err, time.Since(start), info)
		return resp, err
	}
}

// UnaryClientInterceptor gRPC 客户端日志拦截器（仅输出失败调用）
func UnaryClientInterceptor(scope string) grpc.UnaryClientInterceptor {
	return func(ctx context.Context, method string, req, resp interface{}, cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {
		start := time.Now()
		err := invoker(ctx, method, req, resp, cc, opts...)
		if !utils.IsEmptyError(err) {
			logGRPCClient(ctx, scope, method, req, resp, err, time.Since(start))
		}
		return err
	}
}

// logGRPCServer 输出 gRPC 服务端日志
func logGRPCServer(ctx context.Context, scope string, req, resp interface{}, err error, cost time.Duration, info *grpc.UnaryServerInfo) {
	requestId := RequestIdFromGRPCContext(ctx)
	if requestId == "" {
		requestId = "-"
	}
	fullMethod := ""
	if info != nil {
		fullMethod = info.FullMethod
	}
	if utils.IsEmptyError(err) {
		utils.LogInfo(scope, consts.LogColorBoldGreen, "[%s] %s %13v %s [%s] req=%s",
			requestId,
			consts.LogColorBoldBlackOnWhite, cost, consts.LogColorReset,
			gRPCMethod(fullMethod),
			utils.JSONF(req),
		)
		return
	}
	utils.LogError(scope, "[%s] %s %13v %s [%s] req=%s resp=%s err=%s",
		requestId,
		consts.LogColorBoldBlackOnWhite, cost, consts.LogColorReset,
		gRPCMethod(fullMethod),
		utils.JSONF(req),
		utils.JSONF(resp),
		utils.FormatErrorLog(gRPCErrorMessage(err)),
	)
}

// logGRPCClient 输出 gRPC 客户端日志
func logGRPCClient(ctx context.Context, scope, method string, req, resp interface{}, err error, cost time.Duration) {
	requestId := RequestIdFromGRPCContext(ctx)
	if requestId == "" {
		requestId = "-"
	}
	utils.LogError(scope, "[%s] %s %13v %s [%s] req=%s resp=%s err=%s",
		requestId,
		consts.LogColorBoldBlackOnWhite, cost, consts.LogColorReset,
		gRPCMethod(method),
		utils.JSONF(req),
		utils.JSONF(resp),
		utils.FormatErrorLog(gRPCErrorMessage(err)),
	)
}

// gRPCMethod 获取 gRPC 方法名称
func gRPCMethod(fullMethod string) string {
	fullMethod = strings.TrimSpace(strings.TrimPrefix(fullMethod, "/"))
	if fullMethod == "" {
		return "unknown"
	}
	return fullMethod
}

// gRPCErrorMessage 获取 gRPC 错误文本
func gRPCErrorMessage(err error) string {
	if err == nil {
		return ""
	}
	if grpcStatus, ok := status.FromError(err); ok {
		return grpcStatus.Message()
	}
	return err.Error()
}
