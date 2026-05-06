package utils

import (
	"bytes"
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
)

// IsSafePathPart 判断字符串是否可作为单层路径片段
func IsSafePathPart(value string) bool {
	value = strings.TrimSpace(value)
	if value == "" || value == "." || value == ".." {
		return false
	}
	return !strings.ContainsAny(value, `/\`)
}

// ReportPath 获取任务报告路径
func ReportPath(taskDir string) string {
	return filepath.Join(taskDir, "report.json")
}

// ArtifactPath 获取任务产物路径
func ArtifactPath(taskDir, artifactName string) string {
	return filepath.Join(taskDir, strings.TrimSpace(artifactName))
}

// ReadCompactReportJSON 读取任务报告文件并转换为 JSON 压缩字符串
func ReadCompactReportJSON(taskDir string) (string, error) {
	bytesValue, err := os.ReadFile(ReportPath(taskDir))
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
