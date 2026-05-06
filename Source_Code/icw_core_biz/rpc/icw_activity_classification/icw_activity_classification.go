package icw_activity_classification

import (
	"context"

	"icw_common/gen/activity/classification"
	"icw_core_biz/rpc/common"
)

// Start 启动分类任务
func Start(ctx context.Context, client *Client, req *classificationpb.StartRequest, resp *classificationpb.StartResponse) error {
	return common.CallGRPC[classificationpb.StartRequest, classificationpb.StartResponse](ctx, client, req, resp, client.Start)
}
