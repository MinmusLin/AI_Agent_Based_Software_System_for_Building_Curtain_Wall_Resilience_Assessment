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

	req.TaskUuid = strings.TrimSpace(req.TaskUuid)
	req.TaskCode = strings.TrimSpace(req.TaskCode)
	req.ImageUuid = strings.TrimSpace(req.ImageUuid)
	req.PresignGetUrl = strings.TrimSpace(req.PresignGetUrl)

	if req.TaskUuid == "" {
		return errors.New("task uuid is required")
	}
	if !isSafePathPart(req.TaskCode) {
		return errors.New("task code is invalid")
	}
	if _, err := registry.Get(req.TaskCode); err != nil {
		return err
	}
	if !isSafePathPart(req.ImageUuid) {
		return errors.New("image uuid is invalid")
	}
	if req.PresignGetUrl == "" {
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

		artifact.Name = strings.TrimSpace(artifact.Name)
		artifact.PresignUploadUrl = strings.TrimSpace(artifact.PresignUploadUrl)
		artifact.ContentType = strings.TrimSpace(artifact.ContentType)

		if !isSafePathPart(artifact.Name) {
			return errors.New("artifact name is invalid")
		}
		if _, ok := artifactNames[artifact.Name]; ok {
			return errors.New("artifact name is duplicated")
		}
		artifactNames[artifact.Name] = struct{}{}
		if artifact.PresignUploadUrl == "" {
			return errors.New("artifact presign upload url is required")
		}
		if _, err := url.ParseRequestURI(artifact.PresignUploadUrl); err != nil {
			return err
		}
		if artifact.ContentType == "" {
			return errors.New("artifact content type is required")
		}
	}

	return nil
}
