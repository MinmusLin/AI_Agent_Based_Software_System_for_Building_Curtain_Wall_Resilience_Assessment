package detection

import (
	"context"
	"strings"

	"icw_common/enum"
	"icw_common/gen/core/biz"
	"icw_common/rpc/error"

	"icw_core_biz/repositories/minio"
)

// GetReasoningArtifactPolicy 获取图像检测推理产物上传授权
func (s *Service) GetReasoningArtifactPolicy(ctx context.Context, req *bizpb.GetReasoningArtifactPolicyRequest) (*bizpb.GetReasoningArtifactPolicyResponse, error) {
	resp := &bizpb.GetReasoningArtifactPolicyResponse{}
	err := s.CallRPC(req, func() error {
		return s.getReasoningArtifactPolicy(ctx, req, resp)
	})
	return resp, err
}

func (s *Service) getReasoningArtifactPolicy(ctx context.Context, req *bizpb.GetReasoningArtifactPolicyRequest, resp *bizpb.GetReasoningArtifactPolicyResponse) error {
	if req.UserId == 0 {
		return rpc_error.BadRequestDefault("user id is required")
	}
	if req.ProjectId == 0 {
		return rpc_error.BadRequestDefault("project id is required")
	}

	req.ImageUuid = strings.TrimSpace(req.ImageUuid)
	if req.ImageUuid == "" {
		return rpc_error.BadRequestDefault("image uuid is required")
	}
	req.TaskUuid = strings.TrimSpace(req.TaskUuid)
	if req.TaskUuid == "" {
		return rpc_error.BadRequestDefault("task uuid is required")
	}

	taskCode := enum.DetectionTaskCodeString(req.TaskCode)
	if taskCode == "" {
		return rpc_error.BadRequestDefault("detection task code is invalid")
	}

	task, subTask, err := s.MySQL().FindProjectDetectionReasoningTask(ctx, taskCode, req.TaskUuid)
	if err != nil {
		return err
	}
	if task == nil || subTask == nil {
		return rpc_error.BadRequestDefault("reasoning task is not accessible")
	}
	if task.UserId != req.UserId || task.ProjectId != req.ProjectId || task.ImageUuid != req.ImageUuid {
		return rpc_error.BadRequestDefault("reasoning task is not accessible")
	}
	if subTask.UserId != req.UserId || subTask.ProjectId != req.ProjectId || subTask.ImageId != task.ImageId || subTask.MainTaskId != task.Id {
		return rpc_error.BadRequestDefault("reasoning task is not accessible")
	}
	if subTask.Status != bizpb.ProjectDetectionSubTaskStatus_Pending {
		return rpc_error.BadRequestDefault("reasoning task status is invalid")
	}

	keyPrefix, err := minio.GenProjectDetectionArtifactPrefixByTask(req.ProjectId, req.ImageUuid, taskCode)
	if err != nil {
		return err
	}
	policyURL, formData, err := s.MinIO().PresignPostPolicy(ctx, keyPrefix, s.Config().ProjectImageUploadTTL)
	if err != nil {
		return err
	}

	resp.Url = policyURL
	resp.KeyPrefix = keyPrefix
	resp.FormData = formData

	return nil
}
