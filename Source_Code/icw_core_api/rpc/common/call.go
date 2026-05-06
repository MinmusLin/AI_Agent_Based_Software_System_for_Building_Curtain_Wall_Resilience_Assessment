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

// ReadyClient API 层通用 gRPC 调用能力
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

	requestId := apiUtils.GetXRequestId(ctx)
	ctx = utils.AppendRequestIdToOutgoingContext(ctx, requestId)

	if !client.Ready() {
		return rpc_err.InternalErrorDefault("grpc client is nil")
	}
	if invoke == nil {
		return rpc_err.InternalErrorDefault("grpc invoke is nil")
	}

	pbResp, err := invoke(ctx, req)
	if err != nil {
		if grpcStatus, ok := status.FromError(err); ok {
			return errors.New(grpcStatus.Message())
		}
		return err
	}
	if resp == nil || pbResp == nil {
		return rpc_err.InternalErrorDefault("grpc response is nil")
	}
	*resp = *pbResp

	return nil
}
