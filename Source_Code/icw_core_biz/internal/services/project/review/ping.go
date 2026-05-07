package review

import (
	"context"

	"icw_common/gen/core/biz"
)

// Ping .
func (s *Service) Ping(ctx context.Context, req *bizpb.PingReviewRequest) (*bizpb.PingReviewResponse, error) {
	resp := &bizpb.PingReviewResponse{}
	err := s.CallRPC(ctx, req, func() error {
		return s.ping(req, resp)
	})
	return resp, err
}

func (s *Service) ping(_ *bizpb.PingReviewRequest, _ *bizpb.PingReviewResponse) error {
	return nil
}
