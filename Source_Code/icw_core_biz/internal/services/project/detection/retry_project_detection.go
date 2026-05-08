package detection

import (
	"context"

	"icw_common/gen/core/biz"

	"icw_core_biz/repositories/minio"
)

// RetryProjectDetection 重试项目智能检测
func (s *Service) RetryProjectDetection(ctx context.Context, req *bizpb.RetryProjectDetectionRequest) (*bizpb.RetryProjectDetectionResponse, error) {
	resp := &bizpb.RetryProjectDetectionResponse{}
	err := s.CallRPC(ctx, req, func() error {
		return s.retryProjectDetection(ctx, req, resp)
	})
	return resp, err
}

func (s *Service) retryProjectDetection(ctx context.Context, req *bizpb.RetryProjectDetectionRequest, resp *bizpb.RetryProjectDetectionResponse) error {
	imageUuids, tasks, err := s.MySQL().RetryProjectDetectionTasks(ctx, req.UserId, req.ProjectId)
	if err != nil {
		return err
	}
	if len(imageUuids) == 0 {
		return nil
	}

	for _, imageUuid := range imageUuids {
		// 生成项目检测产物对象 Key 前缀
		prefix, err := minio.GenProjectDetectionArtifactPrefix(req.ProjectId, imageUuid)
		if err != nil {
			return err
		}

		if s.Redis() != nil {
			// 按对象 Key 前缀查询对象 Key 列表
			keys, err := s.MinIO().ListObjectKeysByPrefix(ctx, prefix)
			if err != nil {
				return err
			}
			for _, key := range keys {
				_ = s.Redis().ClearPresignURL(ctx, key)
			}
		}

		// 按对象 Key 前缀删除对象
		if err := s.MinIO().RemoveObjectsByPrefix(ctx, prefix); err != nil {
			return err
		}
	}

	if len(tasks) == 0 {
		return nil
	}

	// 投递项目图像检测任务
	resp.TaskCount = s.enqueueProjectDetectionTasks(ctx, tasks)

	return nil
}
