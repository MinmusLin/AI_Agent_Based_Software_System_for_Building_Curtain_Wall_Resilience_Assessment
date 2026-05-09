package utils

import (
	"errors"
	"net/url"
	"strings"

	"icw_common/gen/activity/summary"
)

// ValidateStartDetectionSummaryRequest 校验启动图像检测总结任务请求
func ValidateStartDetectionSummaryRequest(req *summarypb.StartDetectionSummaryRequest) error {
	if req == nil {
		return errors.New("request is nil")
	}

	req.TaskUuid = strings.TrimSpace(req.TaskUuid)
	req.ImageUuid = strings.TrimSpace(req.ImageUuid)
	req.SourceJson = strings.TrimSpace(req.SourceJson)

	if req.TaskUuid == "" {
		return errors.New("task uuid is required")
	}
	if req.ImageUuid == "" {
		return errors.New("image uuid is required")
	}
	if req.SourceJson == "" {
		return errors.New("source json is required")
	}

	return nil
}

// ValidateStartProjectSummaryRequest 校验启动项目总结任务请求
func ValidateStartProjectSummaryRequest(req *summarypb.StartProjectSummaryRequest) error {
	if req == nil {
		return errors.New("request is nil")
	}

	req.ProjectId = strings.TrimSpace(req.ProjectId)
	req.SourceUrl = strings.TrimSpace(req.SourceUrl)

	if req.ProjectId == "" {
		return errors.New("project id is required")
	}
	if req.SourceUrl == "" {
		return errors.New("source url is required")
	}

	parsedURL, err := url.ParseRequestURI(req.SourceUrl)
	if err != nil || parsedURL.Scheme == "" || parsedURL.Host == "" {
		return errors.New("source url is invalid")
	}

	return nil
}
