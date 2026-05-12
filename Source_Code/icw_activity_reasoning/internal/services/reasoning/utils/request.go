package utils

import (
	"errors"
	"net/url"
	"strings"

	"icw_common/enum"
	"icw_common/gen/activity/reasoning"

	"icw_activity_reasoning/internal/detectors/common"
)

// ValidateRequest 校验启动原子检测任务请求
func ValidateRequest(req *reasoningpb.StartRequest, registry *common.Registry) error {
	if req == nil {
		return errors.New("request is nil")
	}

	req.TaskUuid = strings.TrimSpace(req.TaskUuid)
	if req.TaskUuid == "" {
		return errors.New("task uuid is required")
	}

	taskCode := enum.DetectionTaskCodeString(req.TaskCode)
	if !isSafePathPart(taskCode) {
		return errors.New("task code is invalid")
	}
	if _, err := registry.Get(taskCode); err != nil {
		return err
	}

	req.ImageUuid = strings.TrimSpace(req.ImageUuid)
	if !isSafePathPart(req.ImageUuid) {
		return errors.New("image uuid is invalid")
	}

	if req.UserId == 0 {
		return errors.New("user id is required")
	}
	if req.ProjectId == 0 {
		return errors.New("project id is required")
	}

	if req.ArtifactPolicy == nil {
		return errors.New("artifact policy is required")
	}
	req.ArtifactPolicy.Url = strings.TrimSpace(req.ArtifactPolicy.Url)
	req.ArtifactPolicy.KeyPrefix = strings.TrimSpace(req.ArtifactPolicy.KeyPrefix)
	if req.ArtifactPolicy.Url == "" {
		return errors.New("artifact policy url is required")
	}
	if req.ArtifactPolicy.KeyPrefix == "" {
		return errors.New("artifact policy key prefix is required")
	}
	if _, err := url.ParseRequestURI(req.ArtifactPolicy.Url); err != nil {
		return err
	}

	return nil
}
