package icw_activity_reasoning

import (
	"context"

	"icw_common/gen/activity/reasoning"
	"icw_core_biz/rpc/common"
)

// Start 启动推理任务
func Start(ctx context.Context, client *Client, req *reasoningpb.StartRequest, resp *reasoningpb.StartResponse) error {
	return common.CallGRPC[reasoningpb.StartRequest, reasoningpb.StartResponse](ctx, client, req, resp, client.Start)
}
