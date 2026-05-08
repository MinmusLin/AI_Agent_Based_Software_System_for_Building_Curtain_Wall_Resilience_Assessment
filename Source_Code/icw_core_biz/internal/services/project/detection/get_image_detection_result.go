package detection

import (
	"context"
	"path"
	"strings"

	"icw_common/enum"
	"icw_common/gen/activity"
	"icw_common/gen/core/biz"
	"icw_common/rpc/error"

	"icw_core_biz/repositories/minio"
	"icw_core_biz/repositories/mysql/model"
)

// GetImageDetectionResult 获取图像检测结果
func (s *Service) GetImageDetectionResult(ctx context.Context, req *bizpb.GetImageDetectionResultRequest) (*bizpb.GetImageDetectionResultResponse, error) {
	resp := &bizpb.GetImageDetectionResultResponse{}
	err := s.CallRPC(ctx, req, func() error {
		return s.getImageDetectionResult(ctx, req, resp)
	})
	return resp, err
}

func (s *Service) getImageDetectionResult(ctx context.Context, req *bizpb.GetImageDetectionResultRequest, resp *bizpb.GetImageDetectionResultResponse) error {
	req.ImageUuid = strings.TrimSpace(req.ImageUuid)
	if req.ImageUuid == "" {
		return rpc_error.BadRequestDefault("image uuid is required")
	}

	image, err := s.MySQL().FindProjectImageByUuid(ctx, req.UserId, req.ProjectId, req.ImageUuid)
	if err != nil {
		return err
	}
	if image == nil {
		return rpc_error.BadRequestDefault("project image is not accessible")
	}

	imageDTO, err := model.ProjectImageRecordToDTO(ctx, s.MinIO(), s.Redis(), image, s.Config().ProjectImageGetTTL)
	if err != nil {
		return err
	}
	originalURL, err := minio.PresignProjectImageOriginalURL(ctx, s.MinIO(), s.Redis(), req.ProjectId, req.ImageUuid, s.Config().ProjectImageGetTTL)
	if err != nil {
		return err
	}

	resp.Image = imageDTO
	resp.OriginalUrl = originalURL
	resp.TaskCodes = make([]string, 0)

	task, err := s.MySQL().FindProjectDetectionTaskByImageUuid(ctx, req.UserId, req.ProjectId, req.ImageUuid)
	if err != nil {
		return err
	}
	if task == nil {
		return nil
	}

	statusItems, err := s.MySQL().GetProjectDetectionTasks(ctx, req.UserId, req.ProjectId)
	if err != nil {
		return err
	}
	for _, item := range statusItems {
		if item != nil && item.ImageUuid == req.ImageUuid {
			resp.Status = item
			break
		}
	}

	if err := s.fillProjectDetectionTaskResults(ctx, task, resp); err != nil {
		return err
	}
	return s.fillProjectDetectionSummaryResult(ctx, task, resp)
}

func (s *Service) projectDetectionArtifacts(ctx context.Context, task *model.ProjectDetectionTaskRecord, taskCode string) (map[string]string, error) {
	prefix, err := minio.GenProjectDetectionArtifactPrefixByTask(task.ProjectId, task.ImageUuid, taskCode)
	if err != nil {
		return nil, err
	}
	keys, err := s.MinIO().ListObjectKeysByPrefix(ctx, prefix)
	if err != nil {
		return nil, err
	}
	artifacts := make(map[string]string, len(keys))
	for _, key := range keys {
		name := strings.TrimPrefix(key, prefix)
		if name == "" || path.Ext(name) == "" {
			continue
		}
		url, err := minio.PresignProjectDetectionArtifactURL(ctx, s.MinIO(), s.Redis(), key, s.Config().ProjectImageGetTTL)
		if err != nil {
			return nil, err
		}
		artifacts[name] = url
	}
	return artifacts, nil
}

func (s *Service) fillProjectDetectionTaskResults(ctx context.Context, task *model.ProjectDetectionTaskRecord, resp *bizpb.GetImageDetectionResultResponse) error {
	if task.CorrosionShouldExecute && task.CorrosionTaskId.Valid {
		taskCode := enum.DetectionTaskCodeString(activitypb.DetectionTaskCode_Corrosion)
		result, err := s.MySQL().GetProjectDetectionCorrosionResult(ctx, uint64(task.CorrosionTaskId.Int64))
		if err != nil {
			return err
		}
		if result != nil {
			result.Artifacts, err = s.projectDetectionArtifacts(ctx, task, taskCode)
			if err != nil {
				return err
			}
			resp.CorrosionResult = result
			resp.TaskCodes = append(resp.TaskCodes, taskCode)
		}
	}
	if task.CrackShouldExecute && task.CrackTaskId.Valid {
		taskCode := enum.DetectionTaskCodeString(activitypb.DetectionTaskCode_Crack)
		result, err := s.MySQL().GetProjectDetectionCrackResult(ctx, uint64(task.CrackTaskId.Int64))
		if err != nil {
			return err
		}
		if result != nil {
			result.Artifacts, err = s.projectDetectionArtifacts(ctx, task, taskCode)
			if err != nil {
				return err
			}
			resp.CrackResult = result
			resp.TaskCodes = append(resp.TaskCodes, taskCode)
		}
	}
	if task.StainShouldExecute && task.StainTaskId.Valid {
		taskCode := enum.DetectionTaskCodeString(activitypb.DetectionTaskCode_Stain)
		result, err := s.MySQL().GetProjectDetectionStainResult(ctx, uint64(task.StainTaskId.Int64))
		if err != nil {
			return err
		}
		if result != nil {
			result.Artifacts, err = s.projectDetectionArtifacts(ctx, task, taskCode)
			if err != nil {
				return err
			}
			resp.StainResult = result
			resp.TaskCodes = append(resp.TaskCodes, taskCode)
		}
	}
	if task.FlatnessShouldExecute && task.FlatnessTaskId.Valid {
		taskCode := enum.DetectionTaskCodeString(activitypb.DetectionTaskCode_Flatness)
		result, err := s.MySQL().GetProjectDetectionFlatnessResult(ctx, uint64(task.FlatnessTaskId.Int64))
		if err != nil {
			return err
		}
		if result != nil {
			result.Artifacts, err = s.projectDetectionArtifacts(ctx, task, taskCode)
			if err != nil {
				return err
			}
			resp.FlatnessResult = result
			resp.TaskCodes = append(resp.TaskCodes, taskCode)
		}
	}
	if task.SpallingShouldExecute && task.SpallingTaskId.Valid {
		taskCode := enum.DetectionTaskCodeString(activitypb.DetectionTaskCode_Spalling)
		result, err := s.MySQL().GetProjectDetectionSpallingResult(ctx, uint64(task.SpallingTaskId.Int64))
		if err != nil {
			return err
		}
		if result != nil {
			result.Artifacts, err = s.projectDetectionArtifacts(ctx, task, taskCode)
			if err != nil {
				return err
			}
			resp.SpallingResult = result
			resp.TaskCodes = append(resp.TaskCodes, taskCode)
		}
	}
	return nil
}

func (s *Service) fillProjectDetectionSummaryResult(ctx context.Context, task *model.ProjectDetectionTaskRecord, resp *bizpb.GetImageDetectionResultResponse) error {
	if !task.SummaryShouldExecute || !task.SummaryTaskId.Valid {
		return nil
	}
	summaryResult, err := s.MySQL().GetProjectDetectionSummaryTypedResult(ctx, uint64(task.SummaryTaskId.Int64))
	if err != nil {
		return err
	}
	resp.SummaryResult = summaryResult
	return nil
}
