package icw_core_biz

import (
	"context"

	"icw_common/gen/core/biz"
	"icw_common/rpc"
)

// ReportClassificationResult 上报图像检测分类结果
func ReportClassificationResult(ctx context.Context, client *Client, req *bizpb.ReportClassificationResultRequest, resp *bizpb.ReportClassificationResultResponse) error {
	return rpc.CallGRPC[bizpb.ReportClassificationResultRequest, bizpb.ReportClassificationResultResponse](ctx, client, req, resp, client.ReportClassificationResult)
}
