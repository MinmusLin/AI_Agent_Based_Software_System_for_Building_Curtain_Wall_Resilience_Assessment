package icw_activity_reasoning

import (
	"context"

	"icw_common/gen/activity/reasoning"
	"icw_common/rpc"
)

// Start 启动推理任务
func Start(ctx context.Context, client *Client, req *reasoningpb.StartRequest, resp *reasoningpb.StartResponse) error {
	return rpc.CallGRPC[reasoningpb.StartRequest, reasoningpb.StartResponse](ctx, client, req, resp, client.Start)
}
