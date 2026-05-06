package utils

import (
	"context"
	"strings"
	"time"

	"google.golang.org/grpc"

	"icw_common/consts"
)

// GRPCUnaryServerInterceptor gRPC 服务端日志拦截器
func GRPCUnaryServerInterceptor(scope string) grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		start := time.Now()
		resp, err := handler(ctx, req)
		logGRPCServer(ctx, scope, req, resp, err, time.Since(start), info)
		return resp, err
	}
}

// GRPCUnaryClientInterceptor gRPC 客户端日志拦截器
func GRPCUnaryClientInterceptor(scope string) grpc.UnaryClientInterceptor {
	return func(ctx context.Context, method string, req, resp interface{}, cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {
		start := time.Now()
		err := invoker(ctx, method, req, resp, cc, opts...)
		if !IsEmptyError(err) {
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
	if IsEmptyError(err) {
		LogInfo(scope, consts.LogColorBoldGreen, "[%s] %s %13v %s [%s] req=%s resp=%s",
			requestId,
			consts.LogColorBoldBlackOnWhite, cost, consts.LogColorReset,
			gRPCMethod(fullMethod),
			JSONF(req),
			JSONF(resp),
		)
		return
	}
	LogError(scope, "[%s] %s %13v %s [%s] req=%s resp=%s err=%s",
		requestId,
		consts.LogColorBoldBlackOnWhite, cost, consts.LogColorReset,
		gRPCMethod(fullMethod),
		JSONF(req),
		JSONF(resp),
		FormatErrorLog(GRPCErrorMessage(err)),
	)
}

// logGRPCClient 输出 gRPC 客户端日志
func logGRPCClient(ctx context.Context, scope, method string, req, resp interface{}, err error, cost time.Duration) {
	requestId := RequestIdFromGRPCContext(ctx)
	if requestId == "" {
		requestId = "-"
	}
	LogError(scope, "[%s] %s %13v %s [%s] req=%s resp=%s err=%s",
		requestId,
		consts.LogColorBoldBlackOnWhite, cost, consts.LogColorReset,
		gRPCMethod(method),
		JSONF(req),
		JSONF(resp),
		FormatErrorLog(GRPCErrorMessage(err)),
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
