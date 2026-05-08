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
	err := s.CallRPC(req, func() error {
		return s.getImageDetectionResult(ctx, req, resp)
	})
	return resp, err
}

func (s *Service) getImageDetectionResult(ctx context.Context, req *bizpb.GetImageDetectionResultRequest, resp *bizpb.GetImageDetectionResultResponse) error {
	req.ImageUuid = strings.TrimSpace(req.ImageUuid)
	if req.ImageUuid == "" {
		return rpc_error.BadRequestDefault("image uuid is required")
	}

	// 按用户 ID、项目 ID 和图像 UUID 查询图像
	image, err := s.MySQL().FindProjectImageByUuid(ctx, req.UserId, req.ProjectId, req.ImageUuid)
	if err != nil {
		return err
	}
	if image == nil {
		return rpc_error.BadRequestDefault("project image is not accessible")
	}

	// 将 MySQL 数据模型转换为 RPC 数据模型
	imageDTO, err := model.ProjectImageRecordToDTO(ctx, s.MinIO(), s.Redis(), image, s.Config().ProjectImageGetTTL)
	if err != nil {
		return err
	}

	// 获取项目图像原图下载预签名 URL
	originalURL, err := minio.PresignProjectImageOriginalURL(ctx, s.MinIO(), s.Redis(), req.ProjectId, req.ImageUuid, s.Config().ProjectImageGetTTL)
	if err != nil {
		return err
	}

	// 按用户 ID、项目 ID 和图像 UUID 查询项目图像检测主任务
	task, err := s.MySQL().FindProjectDetectionTaskByImageUuid(ctx, req.UserId, req.ProjectId, req.ImageUuid)
	if err != nil {
		return err
	}
	if task == nil {
		return nil
	}

	// 按用户 ID、项目 ID 和图像 UUID 查询项目图像检测任务状态
	status, err := s.MySQL().GetProjectDetectionTaskStatus(ctx, req.UserId, req.ProjectId, req.ImageUuid)
	if err != nil {
		return err
	}

	resp.Image = imageDTO
	resp.OriginalUrl = originalURL
	resp.Status = status
	resp.TaskCodes = make([]activitypb.DetectionTaskCode_Value, 0)

	// 将项目图像检测任务填充进 HTTP 响应
	if err := s.fillProjectDetectionTaskResults(ctx, task, resp); err != nil {
		return err
	}

	// 将项目图像检测总结报告填充进 HTTP 响应
	return s.fillProjectDetectionSummaryResult(ctx, task, resp)
}

// fillProjectDetectionTaskResults 将项目图像检测任务填充进 HTTP 响应
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
			resp.TaskCodes = append(resp.TaskCodes, enum.ParseDetectionTaskCode(taskCode))
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
			resp.TaskCodes = append(resp.TaskCodes, enum.ParseDetectionTaskCode(taskCode))
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
			resp.TaskCodes = append(resp.TaskCodes, enum.ParseDetectionTaskCode(taskCode))
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
			resp.TaskCodes = append(resp.TaskCodes, enum.ParseDetectionTaskCode(taskCode))
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
			resp.TaskCodes = append(resp.TaskCodes, enum.ParseDetectionTaskCode(taskCode))
		}
	}
	return nil
}

// projectDetectionArtifacts 将项目图像检测产物下载预签名 URL 填充进 HTTP 响应
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

// fillProjectDetectionSummaryResult 将项目图像检测总结报告填充进 HTTP 响应
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
