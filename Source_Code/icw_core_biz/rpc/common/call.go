package common

import (
	"context"
	"errors"

	"google.golang.org/grpc"
	"google.golang.org/grpc/status"

	"icw_common/utils"
)

// ReadyClient RPC 层通用 gRPC 调用能力
type ReadyClient interface {
	Ready() bool
}

// CallGRPC 执行 RPC 层通用 gRPC 调用能力
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
	ctx = Context(ctx)

	if utils.IsNil(client) || !client.Ready() {
		return errors.New("grpc client is nil")
	}
	if invoke == nil {
		return errors.New("grpc invoke is nil")
	}

	pbResp, err := invoke(ctx, req)
	if err != nil {
		if grpcStatus, ok := status.FromError(err); ok {
			return errors.New(grpcStatus.Message())
		}
		return err
	}
	if resp == nil || pbResp == nil {
		return errors.New("grpc response is nil")
	}
	*resp = *pbResp

	return nil
}
