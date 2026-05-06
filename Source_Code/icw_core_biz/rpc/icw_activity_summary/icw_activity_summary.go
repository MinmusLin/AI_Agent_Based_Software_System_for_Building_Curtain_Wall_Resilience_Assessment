package icw_activity_summary

import (
	"context"

	"icw_common/gen/activity/summary"
	"icw_core_biz/rpc/common"
)

// Ping 智能总结服务探活
func Ping(ctx context.Context, client *Client, req *summarypb.PingRequest, resp *summarypb.PingResponse) error {
	return common.CallGRPC[summarypb.PingRequest, summarypb.PingResponse](ctx, client, req, resp, client.Ping)
}
