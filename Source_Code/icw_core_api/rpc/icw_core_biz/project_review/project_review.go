package project_review

import (
	"context"

	"icw_common/gen/core/biz"
	"icw_common/rpc"

	"icw_core_api/rpc/icw_core_biz"
	"icw_core_api/utils"
)

// GetProjectDetectionReview 获取图像检测人工复核信息
func GetProjectDetectionReview(ctx context.Context, client *icw_core_biz.Client, req *bizpb.GetProjectDetectionReviewRequest, resp *bizpb.GetProjectDetectionReviewResponse) error {
	return rpc.CallGRPC[bizpb.GetProjectDetectionReviewRequest, bizpb.GetProjectDetectionReviewResponse](ctx, client, req, resp, client.ProjectReview().GetProjectDetectionReview, rpc.WithRequestIdResolver(utils.GetXRequestId))
}

// UpdateProjectDetectionReview 更新图像检测人工复核信息
func UpdateProjectDetectionReview(ctx context.Context, client *icw_core_biz.Client, req *bizpb.UpdateProjectDetectionReviewRequest, resp *bizpb.UpdateProjectDetectionReviewResponse) error {
	return rpc.CallGRPC[bizpb.UpdateProjectDetectionReviewRequest, bizpb.UpdateProjectDetectionReviewResponse](ctx, client, req, resp, client.ProjectReview().UpdateProjectDetectionReview, rpc.WithRequestIdResolver(utils.GetXRequestId))
}
