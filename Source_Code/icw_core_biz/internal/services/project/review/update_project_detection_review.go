package review

import (
	"context"
	"strings"

	"icw_common/gen/core/biz"
	"icw_common/rpc/error"
)

// UpdateProjectDetectionReview 更新图像检测人工复核信息
func (s *Service) UpdateProjectDetectionReview(ctx context.Context, req *bizpb.UpdateProjectDetectionReviewRequest) (*bizpb.UpdateProjectDetectionReviewResponse, error) {
	resp := &bizpb.UpdateProjectDetectionReviewResponse{}
	err := s.CallRPC(req, func() error {
		return s.updateProjectDetectionReview(ctx, req, resp)
	})
	return resp, err
}

func (s *Service) updateProjectDetectionReview(ctx context.Context, req *bizpb.UpdateProjectDetectionReviewRequest, resp *bizpb.UpdateProjectDetectionReviewResponse) error {
	req.TaskUuid = strings.TrimSpace(req.TaskUuid)
	if req.TaskUuid == "" {
		return rpc_error.BadRequestDefault("task uuid is required")
	}

	review, err := s.MySQL().UpdateProjectDetectionReview(ctx, req.UserId, req.ProjectId, req.TaskUuid, req.Verdict, req.Comment)
	if err != nil {
		return err
	}

	resp.Review = review

	return nil
}
