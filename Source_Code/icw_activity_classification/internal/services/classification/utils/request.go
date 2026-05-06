package utils

import (
	"errors"
	"net/url"
	"strings"

	"icw_common/gen/activity/classification"
)

// ValidateRequest 校验启动分类任务请求
func ValidateRequest(req *classificationpb.StartRequest) error {
	if req == nil {
		return errors.New("request is nil")
	}

	req.TaskUuid = strings.TrimSpace(req.TaskUuid)
	req.ImageUuid = strings.TrimSpace(req.ImageUuid)
	req.PresignGetUrl = strings.TrimSpace(req.PresignGetUrl)

	if req.TaskUuid == "" {
		return errors.New("task uuid is required")
	}
	if req.ImageUuid == "" {
		return errors.New("image uuid is required")
	}
	if req.PresignGetUrl == "" {
		return errors.New("presign get url is required")
	}
	if _, err := url.ParseRequestURI(req.PresignGetUrl); err != nil {
		return err
	}

	return nil
}
