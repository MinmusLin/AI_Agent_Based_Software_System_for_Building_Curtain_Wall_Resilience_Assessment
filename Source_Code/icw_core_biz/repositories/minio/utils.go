package minio

import (
	"errors"
	"fmt"
	"net/url"
	"strings"

	"icw_common/utils"
)

// normalizeEndpoint 将 Endpoint 配置标准化为 MinIO SDK 需要的 <host>:<port> 格式
func normalizeEndpoint(endpoint string) (string, bool) {
	endpoint = strings.TrimSpace(endpoint)
	if endpoint == "" {
		return "", false
	}
	parsedURL, err := url.Parse(endpoint)
	if err != nil || parsedURL.Host == "" {
		return "", false
	}
	return parsedURL.Host, parsedURL.Scheme == "https"
}

// GenDefaultAvatarKey 生成用户默认头像对象 Key
func GenDefaultAvatarKey(emailHash string) string {
	return fmt.Sprintf("users/avatars/%s/default.svg", strings.TrimSpace(emailHash))
}

// GenCustomAvatarKey 生成用户自定义头像对象 Key
func GenCustomAvatarKey(emailHash string) string {
	return fmt.Sprintf("users/avatars/%s/custom.png", strings.TrimSpace(emailHash))
}

// GenProjectThumbnailKey 生成项目缩略图对象 Key
func GenProjectThumbnailKey(projectId uint64) (string, error) {
	projectCode := utils.Encode(projectId)
	if projectCode == "" {
		return "", errors.New("project id is invalid")
	}
	return fmt.Sprintf("projects/%s/thumbnail.png", projectCode), nil
}

// GenProjectImageOriginalKey 生成项目图像原图对象 Key
func GenProjectImageOriginalKey(projectId uint64, imageUuid string) (string, error) {
	projectCode := utils.Encode(projectId)
	if projectCode == "" {
		return "", errors.New("project id is invalid")
	}
	imageUuid = strings.TrimSpace(imageUuid)
	if imageUuid == "" {
		return "", errors.New("image uuid is invalid")
	}
	return fmt.Sprintf("projects/%s/assets/%s/original.png", projectCode, imageUuid), nil
}

// GenProjectImageThumbnailKey 生成项目图像缩略图对象 Key
func GenProjectImageThumbnailKey(projectId uint64, imageUuid string) (string, error) {
	projectCode := utils.Encode(projectId)
	if projectCode == "" {
		return "", errors.New("project id is invalid")
	}
	imageUuid = strings.TrimSpace(imageUuid)
	if imageUuid == "" {
		return "", errors.New("image uuid is invalid")
	}
	return fmt.Sprintf("projects/%s/assets/%s/thumbnail.png", projectCode, imageUuid), nil
}

// GenProjectDetectionArtifactKey 生成项目检测产物对象 Key
func GenProjectDetectionArtifactKey(projectId uint64, imageUuid, taskCode, artifactName string) (string, error) {
	projectCode := utils.Encode(projectId)
	if projectCode == "" {
		return "", errors.New("project id is invalid")
	}
	imageUuid = strings.TrimSpace(imageUuid)
	if imageUuid == "" {
		return "", errors.New("image uuid is invalid")
	}
	taskCode = strings.TrimSpace(taskCode)
	if taskCode == "" {
		return "", errors.New("task code is invalid")
	}
	artifactName = strings.TrimSpace(artifactName)
	if artifactName == "" {
		return "", errors.New("artifact name is invalid")
	}
	return fmt.Sprintf("projects/%s/detections/%s/%s/%s", projectCode, imageUuid, taskCode, artifactName), nil
}
