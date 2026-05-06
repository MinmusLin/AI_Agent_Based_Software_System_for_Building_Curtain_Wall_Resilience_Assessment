package rpc

import (
	"context"
	"strings"

	"google.golang.org/grpc/metadata"
)

const (
	// MetadataRequestId gRPC metadata 中透传请求 ID 的 Key
	MetadataRequestId = "x-request-id"
)

// Context 生成 gRPC 出站调用上下文
func Context(ctx context.Context) context.Context {
	if ctx == nil {
		ctx = context.Background()
	}
	return AppendRequestIdToOutgoingContext(ctx, RequestIdFromGRPCContext(ctx))
}

// AppendRequestIdToOutgoingContext 将请求 ID 写入 gRPC 出站元数据
func AppendRequestIdToOutgoingContext(ctx context.Context, requestId string) context.Context {
	if ctx == nil {
		ctx = context.Background()
	}
	requestId = strings.TrimSpace(requestId)
	if requestId == "" {
		return ctx
	}
	if RequestIdFromOutgoingContext(ctx) == requestId {
		return ctx
	}
	return metadata.AppendToOutgoingContext(ctx, MetadataRequestId, requestId)
}

// RequestIdFromGRPCContext 从 gRPC 入站元数据 / 出站元数据中读取请求 ID
func RequestIdFromGRPCContext(ctx context.Context) string {
	requestId := RequestIdFromIncomingContext(ctx)
	if requestId != "" {
		return requestId
	}
	return RequestIdFromOutgoingContext(ctx)
}

// RequestIdFromIncomingContext 从 gRPC 入站元数据中读取请求 ID
func RequestIdFromIncomingContext(ctx context.Context) string {
	if ctx == nil {
		return ""
	}
	values := metadata.ValueFromIncomingContext(ctx, MetadataRequestId)
	if len(values) == 0 {
		return ""
	}
	return strings.TrimSpace(values[0])
}

// RequestIdFromOutgoingContext 从 gRPC 出站元数据中读取请求 ID
func RequestIdFromOutgoingContext(ctx context.Context) string {
	if ctx == nil {
		return ""
	}
	md, ok := metadata.FromOutgoingContext(ctx)
	if !ok {
		return ""
	}
	values := md.Get(MetadataRequestId)
	if len(values) == 0 {
		return ""
	}
	return strings.TrimSpace(values[len(values)-1])
}
