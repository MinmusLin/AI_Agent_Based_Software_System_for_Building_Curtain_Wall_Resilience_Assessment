package rpc

import (
	"context"
	"errors"

	"google.golang.org/grpc"

	"icw_common/utils"
)

// ReadyClient gRPC 调用通用 Client 能力
type ReadyClient interface {
	Ready() bool
}

// RequestIdResolver 从调用方上下文中解析请求 ID
type RequestIdResolver func(context.Context) string

// CallOption gRPC 调用配置项类型
type CallOption func(*callOptions)

// callOptions gRPC 调用配置项
type callOptions struct {
	requestIdResolver RequestIdResolver
}

// WithRequestIdResolver 配置请求 ID 解析函数
func WithRequestIdResolver(resolver RequestIdResolver) CallOption {
	return func(options *callOptions) {
		options.requestIdResolver = resolver
	}
}

// CallGRPC 执行通用 gRPC 调用
func CallGRPC[PBReq any, PBResp any](
	ctx context.Context,
	client ReadyClient,
	req *PBReq,
	resp *PBResp,
	invoke func(context.Context, *PBReq, ...grpc.CallOption) (*PBResp, error),
	options ...CallOption,
) error {
	if ctx == nil {
		ctx = context.Background()
	}
	callOptions := callOptions{}
	for _, option := range options {
		if option != nil {
			option(&callOptions)
		}
	}
	requestId := RequestIdFromGRPCContext(ctx)
	if callOptions.requestIdResolver != nil {
		if resolved := callOptions.requestIdResolver(ctx); resolved != "" {
			requestId = resolved
		}
	}
	ctx = WithRequestIdToOutgoingContext(ctx, requestId)

	if utils.IsNil(client) || !client.Ready() {
		return errors.New("grpc client is nil")
	}
	if invoke == nil {
		return errors.New("grpc invoke is nil")
	}

	pbResp, err := invoke(ctx, req)
	if err != nil {
		return errors.New(gRPCErrorMessage(err))
	}
	if resp == nil || pbResp == nil {
		return errors.New("grpc response is nil")
	}
	*resp = *pbResp

	return nil
}
