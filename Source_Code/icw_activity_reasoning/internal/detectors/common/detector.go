package common

import (
	"context"
	"errors"
	"fmt"
	"os"
	"path/filepath"
)

// PythonDetector Python 原子检测能力执行器
type PythonDetector struct {
	code        string
	description string
	runtimeRoot string
	worker      *pythonWorker
}

// NewPythonDetector 创建 Python 原子检测能力执行器
func NewPythonDetector(code, description, runtimeRoot string) *PythonDetector {
	return &PythonDetector{
		code:        code,
		description: description,
		runtimeRoot: runtimeRoot,
		worker:      newPythonWorker(code, runtimeRoot),
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
	taskDir := filepath.Join(d.runtimeRoot, d.code, imageUuid)
	imagePath := filepath.Join(taskDir, "original.png")
	if info, err := os.Stat(imagePath); err != nil || info.IsDir() {
		return fmt.Errorf("original image not found: %s", imagePath)
	}
	return d.worker.Run(ctx, imageUuid)
}
