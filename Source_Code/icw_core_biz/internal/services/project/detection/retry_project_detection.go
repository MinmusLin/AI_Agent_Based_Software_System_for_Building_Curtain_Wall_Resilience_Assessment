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
	imageUuids, err := s.MySQL().ListProjectDetectionImageUuidsByStatus(ctx, req.UserId, req.ProjectId, bizpb.ProjectDetectionTaskStatus_Failed)
	if err != nil {
		return err
	}
	if len(imageUuids) == 0 {
		return nil
	}

	for _, imageUuid := range imageUuids {
		prefix, err := minio.GenProjectDetectionArtifactPrefix(req.ProjectId, imageUuid)
		if err != nil {
			return err
		}
		if s.Redis() != nil {
			keys, err := s.MinIO().ListObjectKeysByPrefix(ctx, prefix)
			if err != nil {
				return err
			}
			for _, key := range keys {
				_ = s.Redis().ClearPresignURL(ctx, key)
			}
		}
		if err := s.MinIO().RemoveObjectsByPrefix(ctx, prefix); err != nil {
			return err
		}
	}

	if err := s.MySQL().DeleteProjectDetectionTasksByStatus(ctx, req.UserId, req.ProjectId, bizpb.ProjectDetectionTaskStatus_Failed); err != nil {
		return err
	}

	tasks, err := s.MySQL().CreateProjectDetectionTasksByImageUuids(ctx, req.UserId, req.ProjectId, imageUuids)
	if err != nil {
		return err
	}
	resp.TaskCount = s.enqueueProjectDetectionTasks(ctx, tasks)
	return nil
}
