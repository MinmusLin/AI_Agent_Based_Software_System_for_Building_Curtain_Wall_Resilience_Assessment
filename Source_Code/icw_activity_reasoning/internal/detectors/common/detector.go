package common

import (
	"context"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"icw_activity_reasoning/internal/detectors/utils"
)

// PythonDetector Python 原子检测能力执行器
type PythonDetector struct {
	code        string
	description string
	pythonBin   string
	runner      string
	path        string
	runtimeRoot string
}

// NewPythonDetector 创建 Python 原子检测能力执行器
func NewPythonDetector(code, description, pythonBin, detectorPath, runtimeRoot string) *PythonDetector {
	return &PythonDetector{
		code:        strings.TrimSpace(code),
		description: strings.TrimSpace(description),
		pythonBin:   strings.TrimSpace(pythonBin),
		runner:      utils.AbsPath("python/runner.py"),
		path:        utils.AbsPath(detectorPath),
		runtimeRoot: utils.AbsPath(runtimeRoot),
	}
}

// Code 获取原子检测能力代码
func (d *PythonDetector) Code() string {
	if d == nil {
		return ""
	}
	return d.code
}

// Description 获取原子检测能力描述
func (d *PythonDetector) Description() string {
	if d == nil {
		return ""
	}
	return d.description
}

// Detect 执行 Python 原子检测脚本并读取产物
func (d *PythonDetector) Detect(ctx context.Context, imageUuid string) error {
	if d == nil {
		return errors.New("python detector is nil")
	}
	imageUuid = strings.TrimSpace(imageUuid)
	if imageUuid == "" {
		return errors.New("image uuid is required")
	}
	taskDir := filepath.Join(d.runtimeRoot, d.code, imageUuid)
	imagePath := filepath.Join(taskDir, "original.png")
	if info, err := os.Stat(imagePath); err != nil || info.IsDir() {
		return fmt.Errorf("original image not found: %s", imagePath)
	}
	args := []string{
		d.runner,
		"--task-code", d.code,
		"--detector-path", d.path,
		"--image-uuid", imageUuid,
		"--runtime-root", d.runtimeRoot,
	}
	cmd := exec.CommandContext(
		ctx,
		d.pythonBin,
		args...,
	)
	cmd.Dir = filepath.Dir(d.runner)
	_, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("run python detector %s failed: %v", d.code, err)
	}
	return nil
}
