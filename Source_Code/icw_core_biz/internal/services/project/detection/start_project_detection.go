package detection

import (
	"context"

	"icw_common/gen/core/biz"

	"icw_core_biz/repositories/mysql/model"
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
	// 按用户 ID 和项目 ID 查询项目图像检测任务状态
	currentTasks, err := s.MySQL().GetProjectDetectionTasks(ctx, req.UserId, req.ProjectId)
	if err != nil {
		return err
	}
	if len(currentTasks) > 0 {
		return nil
	}

	// 按用户 ID 和项目 ID 创建项目图像检测主任务
	tasks, err := s.MySQL().CreateProjectDetectionTasks(ctx, req.UserId, req.ProjectId)
	if err != nil {
		return err
	}

	// 投递项目图像检测任务
	resp.TaskCount = s.enqueueProjectDetectionTasks(ctx, tasks)

	return nil
}

// enqueueProjectDetectionTasks 投递项目图像检测任务
func (s *Service) enqueueProjectDetectionTasks(ctx context.Context, tasks []*model.ProjectDetectionTaskRecord) uint32 {
	var taskCount uint32
	for _, task := range tasks {
		if task == nil {
			continue
		}
		s.DetectionWorker().Enqueue(ctx, task.Id)
		taskCount++
	}
	return taskCount
}
