package icw_activity_summary

import (
	"context"

	"icw_common/gen/activity/summary"
	"icw_common/rpc"
)

// StartDetectionSummary 启动图像检测总结任务
func StartDetectionSummary(ctx context.Context, client *Client, req *summarypb.StartDetectionSummaryRequest, resp *summarypb.StartDetectionSummaryResponse) error {
	return rpc.CallGRPC[summarypb.StartDetectionSummaryRequest, summarypb.StartDetectionSummaryResponse](ctx, client, req, resp, client.StartDetectionSummary)
}
