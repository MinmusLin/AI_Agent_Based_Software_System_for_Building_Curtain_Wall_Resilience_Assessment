package minio

import (
	"errors"
	"fmt"
	"net/url"
	"strings"

	"icw_core_biz/utils"
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
	if projectId == 0 {
		return "", errors.New("project id is invalid")
	}
	projectCode := utils.Encode(projectId)
	if projectCode == "" {
		return "", errors.New("project id is invalid")
	}
	return fmt.Sprintf("projects/%s/thumbnail.png", projectCode), nil
}
