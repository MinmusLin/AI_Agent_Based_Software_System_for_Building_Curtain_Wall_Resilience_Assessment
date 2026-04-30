package minio

import (
	"net/url"
	"strings"
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
