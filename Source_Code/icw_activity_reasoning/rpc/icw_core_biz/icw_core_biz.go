package icw_core_biz

import (
	"context"

	"icw_common/gen/core/biz"
	"icw_common/rpc"
)

// GetProjectImageOriginal 获取项目图像原图下载地址
func GetProjectImageOriginal(ctx context.Context, client *Client, req *bizpb.GetProjectImageOriginalRequest, resp *bizpb.GetProjectImageOriginalResponse) error {
	return rpc.CallGRPC[bizpb.GetProjectImageOriginalRequest, bizpb.GetProjectImageOriginalResponse](ctx, client, req, resp, client.GetProjectImageOriginal)
}

// ReportReasoningResult 上报图像检测推理结果
func ReportReasoningResult(ctx context.Context, client *Client, req *bizpb.ReportReasoningResultRequest, resp *bizpb.ReportReasoningResultResponse) error {
	return rpc.CallGRPC[bizpb.ReportReasoningResultRequest, bizpb.ReportReasoningResultResponse](ctx, client, req, resp, client.ReportReasoningResult)
}
