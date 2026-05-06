package common

import (
	"context"
	"errors"
	"reflect"

	"google.golang.org/grpc"
)

// ReadyClient BIZ 层通用 gRPC 调用能力
type ReadyClient interface {
	Ready() bool
}

// CallGRPC 执行 BIZ 层通用 gRPC 调用
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

	if isNil(client) || !client.Ready() {
		return errors.New("grpc client is nil")
	}
	if invoke == nil {
		return errors.New("grpc invoke is nil")
	}

	pbResp, err := invoke(ctx, req)
	if err != nil {
		return NormalizeError(err)
	}
	if resp == nil || pbResp == nil {
		return errors.New("grpc response is nil")
	}
	*resp = *pbResp

	return nil
}

func isNil(value interface{}) bool {
	if value == nil {
		return true
	}
	reflectValue := reflect.ValueOf(value)
	switch reflectValue.Kind() {
	case reflect.Chan, reflect.Func, reflect.Interface, reflect.Map, reflect.Ptr, reflect.Slice:
		return reflectValue.IsNil()
	default:
		return false
	}
}
