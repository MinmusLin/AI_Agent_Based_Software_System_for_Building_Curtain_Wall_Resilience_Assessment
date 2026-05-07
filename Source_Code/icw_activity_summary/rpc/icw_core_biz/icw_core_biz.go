package icw_core_biz

import (
	"context"

	"icw_common/gen/core/biz"
	"icw_common/rpc"
)

// ReportDetectionSummaryResult 上报图像检测总结结果
func ReportDetectionSummaryResult(ctx context.Context, client *Client, req *bizpb.ReportDetectionSummaryResultRequest, resp *bizpb.ReportDetectionSummaryResultResponse) error {
	return rpc.CallGRPC[bizpb.ReportDetectionSummaryResultRequest, bizpb.ReportDetectionSummaryResultResponse](ctx, client, req, resp, client.ReportDetectionSummaryResult)
}

// ReportProjectSummaryResult 上报项目总结结果
func ReportProjectSummaryResult(ctx context.Context, client *Client, req *bizpb.ReportProjectSummaryResultRequest, resp *bizpb.ReportProjectSummaryResultResponse) error {
	return rpc.CallGRPC[bizpb.ReportProjectSummaryResultRequest, bizpb.ReportProjectSummaryResultResponse](ctx, client, req, resp, client.ReportProjectSummaryResult)
}
