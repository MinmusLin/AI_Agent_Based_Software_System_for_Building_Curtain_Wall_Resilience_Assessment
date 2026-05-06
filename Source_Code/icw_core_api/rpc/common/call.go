package common

import (
	"context"
	"errors"

	"google.golang.org/grpc"

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
	ctx = utils.AppendRequestIdToOutgoingContext(ctx, apiUtils.GetXRequestId(ctx))

	if utils.IsNil(client) || !client.Ready() {
		return errors.New("grpc client is nil")
	}
	if invoke == nil {
		return errors.New("grpc invoke is nil")
	}

	pbResp, err := invoke(ctx, req)
	if err != nil {
		return errors.New(utils.GRPCErrorMessage(err))
	}
	if resp == nil || pbResp == nil {
		return errors.New("grpc response is nil")
	}
	*resp = *pbResp

	return nil
}
