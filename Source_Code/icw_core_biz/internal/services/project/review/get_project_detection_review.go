package review

import (
	"context"
	"strings"

	"icw_common/gen/core/biz"
	"icw_common/rpc/error"
)

// GetProjectDetectionReview 获取图像检测人工复核信息
func (s *Service) GetProjectDetectionReview(ctx context.Context, req *bizpb.GetProjectDetectionReviewRequest) (*bizpb.GetProjectDetectionReviewResponse, error) {
	resp := &bizpb.GetProjectDetectionReviewResponse{}
	err := s.CallRPC(req, func() error {
		return s.getProjectDetectionReview(ctx, req, resp)
	})
	return resp, err
}

func (s *Service) getProjectDetectionReview(ctx context.Context, req *bizpb.GetProjectDetectionReviewRequest, resp *bizpb.GetProjectDetectionReviewResponse) error {
	req.TaskUuid = strings.TrimSpace(req.TaskUuid)
	if req.TaskUuid == "" {
		return rpc_error.BadRequestDefault("task uuid is required")
	}

	review, err := s.MySQL().GetProjectDetectionReview(ctx, req.UserId, req.ProjectId, req.TaskUuid)
	if err != nil {
		return err
	}

	resp.Review = review

	return nil
}
