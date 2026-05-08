package detection

import (
	"context"

	"icw_common/gen/core/biz"
)

// GetProjectDetectionTasks 获取项目检测任务列表
func (s *Service) GetProjectDetectionTasks(ctx context.Context, req *bizpb.GetProjectDetectionTasksRequest) (*bizpb.GetProjectDetectionTasksResponse, error) {
	resp := &bizpb.GetProjectDetectionTasksResponse{}
	err := s.CallRPC(req, func() error {
		return s.getProjectDetectionTasks(ctx, req, resp)
	})
	return resp, err
}

func (s *Service) getProjectDetectionTasks(ctx context.Context, req *bizpb.GetProjectDetectionTasksRequest, resp *bizpb.GetProjectDetectionTasksResponse) error {
	tasks, err := s.MySQL().GetProjectDetectionTasksStatus(ctx, req.UserId, req.ProjectId)
	if err != nil {
		return err
	}
	resp.Tasks = tasks
	return nil
}
