package utils

import (
	"path/filepath"
	"strings"
)

// AbsPath 将相对路径转换为绝对路径
func AbsPath(path string) string {
	path = strings.TrimSpace(path)
	if path == "" {
		return ""
	}
	abs, err := filepath.Abs(path)
	if err != nil {
		return path
	}
	return abs
}
