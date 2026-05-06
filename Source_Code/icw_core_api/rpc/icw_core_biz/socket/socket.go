package socket

import (
	"context"

	"icw_common/gen/core/biz"
	"icw_common/rpc"
	"icw_core_api/rpc/icw_core_biz"
	"icw_core_api/utils"
)

// CreateSocketTicket 创建 WebSocket 连接票据
func CreateSocketTicket(ctx context.Context, client *icw_core_biz.Client, req *bizpb.CreateSocketTicketRequest, resp *bizpb.CreateSocketTicketResponse) error {
	return rpc.CallGRPC[bizpb.CreateSocketTicketRequest, bizpb.CreateSocketTicketResponse](ctx, client, req, resp, client.Socket().CreateSocketTicket, rpc.WithRequestIdResolver(utils.GetXRequestId))
}

// ValidateSocketTicket 校验 WebSocket 连接票据
func ValidateSocketTicket(ctx context.Context, client *icw_core_biz.Client, req *bizpb.ValidateSocketTicketRequest, resp *bizpb.ValidateSocketTicketResponse) error {
	return rpc.CallGRPC[bizpb.ValidateSocketTicketRequest, bizpb.ValidateSocketTicketResponse](ctx, client, req, resp, client.Socket().ValidateSocketTicket, rpc.WithRequestIdResolver(utils.GetXRequestId))
}
