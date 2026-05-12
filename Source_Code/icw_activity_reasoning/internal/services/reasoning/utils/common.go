package utils

import (
	"bytes"
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
)

// ArtifactUploadPolicy 图像检测产物上传授权
type ArtifactUploadPolicy struct {
	URL       string
	KeyPrefix string
	FormData  map[string]string
}

// isSafePathPart 判断字符串是否可作为单层路径片段
func isSafePathPart(value string) bool {
	value = strings.TrimSpace(value)
	if value == "" || value == "." || value == ".." {
		return false
	}
	return !strings.ContainsAny(value, `/\`)
}

// reportPath 获取任务报告路径
func reportPath(taskDir string) string {
	return filepath.Join(taskDir, "report.json")
}

// artifactPath 获取任务产物路径
func artifactPath(taskDir, artifactName string) string {
	return filepath.Join(taskDir, artifactName)
}

// ReadCompactReportJSON 读取任务报告文件并转换为 JSON 压缩字符串
func ReadCompactReportJSON(taskDir string) (string, error) {
	bytesValue, err := os.ReadFile(reportPath(taskDir))
	if err != nil {
		return "{}", err
	}
	buffer := &bytes.Buffer{}
	if err := json.Compact(buffer, bytesValue); err != nil {
		return "{}", err
	}
	if strings.TrimSpace(buffer.String()) == "" {
		return "{}", os.ErrInvalid
	}
	return strings.TrimSpace(buffer.String()), nil
}
