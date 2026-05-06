package utils

import (
	"context"
	"strings"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/status"

	"icw_common/consts"
)

// GRPCUnaryServerInterceptor 生成 gRPC 服务端入口日志拦截器
func GRPCUnaryServerInterceptor(scope string) grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		start := time.Now()
		resp, err := handler(ctx, req)
		logGRPCServer(scope, ctx, req, resp, info, time.Since(start), err)
		return resp, err
	}
}

// GRPCUnaryClientInterceptor 生成只记录失败调用的 gRPC 客户端日志拦截器
func GRPCUnaryClientInterceptor(scope, psm string) grpc.UnaryClientInterceptor {
	return func(ctx context.Context, method string, req, resp interface{}, cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {
		start := time.Now()
		err := invoker(ctx, method, req, resp, cc, opts...)
		if !IsEmptyError(err) {
			logGRPCClient(scope, ctx, method, req, resp, time.Since(start), err)
		}
		return err
	}
}

// gRPCMethod 将 /package.Service/Method 格式转换为 Service.Method
func gRPCMethod(fullMethod string) string {
	fullMethod = strings.TrimSpace(strings.TrimPrefix(fullMethod, "/"))
	if fullMethod == "" {
		return "unknown"
	}
	return fullMethod
}

func logGRPCServer(scope string, ctx context.Context, req, resp interface{}, info *grpc.UnaryServerInfo, cost time.Duration, err error) {
	requestId := RequestIdFromGRPCContext(ctx)
	if requestId == "" {
		requestId = "-"
	}
	fullMethod := ""
	if info != nil {
		fullMethod = info.FullMethod
	}
	if IsEmptyError(err) {
		LogInfo(scope, consts.LogColorBoldGreen, "[%s] %s %13v %s [%s] req=%s resp=%s err=%s",
			requestId,
			consts.LogColorBoldBlackOnWhite, cost, consts.LogColorReset,
			gRPCMethod(fullMethod),
			JSONF(req),
			JSONF(resp),
			FormatErrorLog(grpcErrorMessage(err)),
		)
		return
	}
	LogError(scope, "[%s] %s %13v %s [%s] req=%s resp=%s err=%s",
		requestId,
		consts.LogColorBoldBlackOnWhite, cost, consts.LogColorReset,
		gRPCMethod(fullMethod),
		JSONF(req),
		JSONF(resp),
		FormatErrorLog(grpcErrorMessage(err)),
	)
}

func logGRPCClient(scope string, ctx context.Context, method string, req, resp interface{}, cost time.Duration, err error) {
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
		FormatErrorLog(grpcErrorMessage(err)),
	)
}

func grpcErrorMessage(err error) string {
	if err == nil {
		return ""
	}
	if grpcStatus, ok := status.FromError(err); ok {
		return grpcStatus.Message()
	}
	return err.Error()
}
