package common

import (
	"context"
	"errors"

	"google.golang.org/grpc"
	"google.golang.org/grpc/status"

	"icw_common/rpc_err"
	"icw_common/utils"
	apiUtils "icw_core_api/utils"
)

// ReadyClient 定义可被 API 层调用的 gRPC Client 基础能力
type ReadyClient interface {
	Ready() bool
}

// CallGRPC 执行 API 层通用 gRPC 调用
func CallGRPC[PBReq any, PBResp any](
	ctx context.Context,
	client ReadyClient,
	req *PBReq,
	resp *PBResp,
	invoke func(context.Context, *PBReq, ...grpc.CallOption) (*PBResp, error),
) error {
	if ctx == nil {
		ctx = context.Background()
	}

	requestId := apiUtils.GetRequestId(ctx)
	ctx = utils.AppendRequestIdToOutgoingContext(ctx, requestId)

	if client == nil || !client.Ready() {
		return rpc_err.InternalErrorDefault("grpc client is nil")
	}
	if invoke == nil {
		return rpc_err.InternalErrorDefault("grpc invoke is nil")
	}

	pbResp, err := invoke(ctx, req)
	if err != nil {
		return normalizeGRPCError(err)
	}
	if resp == nil || pbResp == nil {
		return rpc_err.InternalErrorDefault("grpc response is nil")
	}
	*resp = *pbResp
	return nil
}

func normalizeGRPCError(err error) error {
	if err == nil {
		return nil
	}
	if grpcStatus, ok := status.FromError(err); ok {
		return errors.New(grpcStatus.Message())
	}
	return err
}
