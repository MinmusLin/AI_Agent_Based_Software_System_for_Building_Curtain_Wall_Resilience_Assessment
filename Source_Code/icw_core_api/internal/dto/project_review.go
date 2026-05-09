package dto

import (
	"icw_common/gen/core/api"
	"icw_common/gen/core/biz"
)

func NewGetProjectDetectionReviewResponse(resp *bizpb.GetProjectDetectionReviewResponse) *apipb.GetProjectDetectionReviewResponse {
	if resp == nil {
		return nil
	}
	return &apipb.GetProjectDetectionReviewResponse{
		Review: resp.Review,
	}
}

func NewUpdateProjectDetectionReviewResponse(resp *bizpb.UpdateProjectDetectionReviewResponse) *apipb.UpdateProjectDetectionReviewResponse {
	if resp == nil {
		return nil
	}
	return &apipb.UpdateProjectDetectionReviewResponse{
		Review: resp.Review,
	}
}
