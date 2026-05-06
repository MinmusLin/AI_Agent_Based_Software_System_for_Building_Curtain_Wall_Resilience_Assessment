package icw_activity_summary

import (
	"context"

	"icw_common/gen/activity/summary"
	"icw_common/rpc"
)

// Ping 总结能力服务探活
func Ping(ctx context.Context, client *Client, req *summarypb.PingRequest, resp *summarypb.PingResponse) error {
	return rpc.CallGRPC[summarypb.PingRequest, summarypb.PingResponse](ctx, client, req, resp, client.Ping)
}
