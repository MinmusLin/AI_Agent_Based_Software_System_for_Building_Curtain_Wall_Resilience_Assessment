package detection

import (
	"context"

	"icw_common/gen/core/biz"
)

// StartProjectDetection 启动项目智能检测
func (s *Service) StartProjectDetection(ctx context.Context, req *bizpb.StartProjectDetectionRequest) (*bizpb.StartProjectDetectionResponse, error) {
	resp := &bizpb.StartProjectDetectionResponse{}
	err := s.CallRPC(ctx, req, func() error {
		return s.startProjectDetection(ctx, req, resp)
	})
	return resp, err
}

func (s *Service) startProjectDetection(ctx context.Context, req *bizpb.StartProjectDetectionRequest, resp *bizpb.StartProjectDetectionResponse) error {
	tasks, err := s.MySQL().CreateProjectDetectionTasks(ctx, req.UserId, req.ProjectId)
	if err != nil {
		return err
	}
	var taskCount uint32
	for _, task := range tasks {
		if task == nil {
			continue
		}
		s.DetectionWorker().Enqueue(ctx, task.Id)
		taskCount++
	}
	resp.TaskCount = taskCount
	return nil
}
