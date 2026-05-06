package classification

import (
	"context"

	"icw_common/enum"
	"icw_common/gen/activity"
	"icw_common/gen/activity/classification"
)

// Start 启动分类任务
func (s *Service) Start(ctx context.Context, req *classificationpb.StartRequest) (*classificationpb.StartResponse, error) {
	resp := &classificationpb.StartResponse{}
	err := s.CallRPC(ctx, req, func() error {
		resp.TaskCodes = []string{
			enum.DetectionTaskCodeString(activitypb.DetectionTaskCode_Corrosion),
			enum.DetectionTaskCodeString(activitypb.DetectionTaskCode_Crack),
			enum.DetectionTaskCodeString(activitypb.DetectionTaskCode_Stain),
			enum.DetectionTaskCodeString(activitypb.DetectionTaskCode_Flatness),
			enum.DetectionTaskCodeString(activitypb.DetectionTaskCode_Spalling),
		}
		return nil
	})
	return resp, err
}
