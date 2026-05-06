package utils

import (
	"errors"
	"net/url"
	"strings"

	"icw_activity_reasoning/internal/detectors/common"
	"icw_common/gen/activity/reasoning"
)

// ValidateRequest 校验启动原子检测任务请求
func ValidateRequest(req *reasoningpb.StartRequest, registry *common.Registry) error {
	if req == nil {
		return errors.New("request is nil")
	}
	if strings.TrimSpace(req.TaskUuid) == "" {
		return errors.New("task uuid is required")
	}
	if strings.TrimSpace(req.TaskCode) == "" {
		return errors.New("task code is required")
	}
	if _, err := registry.Get(req.TaskCode); err != nil {
		return err
	}
	if strings.TrimSpace(req.ImageUuid) == "" {
		return errors.New("image uuid is required")
	}
	if !IsSafePathPart(req.ImageUuid) {
		return errors.New("image uuid is invalid")
	}
	if strings.TrimSpace(req.PresignGetUrl) == "" {
		return errors.New("presign get url is required")
	}
	if _, err := url.ParseRequestURI(req.PresignGetUrl); err != nil {
		return err
	}
	artifactNames := map[string]struct{}{}
	for _, artifact := range req.Artifacts {
		if artifact == nil {
			return errors.New("artifact is nil")
		}
		artifactName := strings.TrimSpace(artifact.Name)
		if artifactName == "" {
			return errors.New("artifact name is required")
		}
		if !IsSafePathPart(artifactName) {
			return errors.New("artifact name is invalid")
		}
		if _, ok := artifactNames[artifactName]; ok {
			return errors.New("artifact name is duplicated")
		}
		artifactNames[artifactName] = struct{}{}
		if strings.TrimSpace(artifact.PresignUploadUrl) == "" {
			return errors.New("artifact presign upload url is required")
		}
		if _, err := url.ParseRequestURI(artifact.PresignUploadUrl); err != nil {
			return err
		}
		if strings.TrimSpace(artifact.ContentType) == "" {
			return errors.New("artifact content type is required")
		}
	}
	return nil
}
