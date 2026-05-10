package common

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"icw_common/utils"
)

// Detector 原子检测能力
type Detector interface {
	Code() string
	Description() string
	ModelName() string
	Detect(ctx context.Context, imageUuid string) error
}

// DetectorMeta 原子检测能力元数据
type DetectorMeta struct {
	Code        string
	Description string
	ModelName   string
}

// NewDetectorMeta 创建原子检测能力元数据
func NewDetectorMeta(code, description, modelName string) *DetectorMeta {
	return &DetectorMeta{
		Code:        code,
		Description: description,
		ModelName:   modelName,
	}
}

// Registry 原子检测能力注册表
type Registry struct {
	detectors map[string]Detector
}

// NewRegistry 创建原子检测能力注册表
func NewRegistry(runtimeRoot, modelRoot string, metas []*DetectorMeta) (*Registry, error) {
	items := make([]Detector, 0, len(metas))
	for _, item := range metas {
		if item == nil {
			continue
		}
		modelPath, err := resolveModelPath(modelRoot, item.Code, item.ModelName)
		if err != nil {
			return nil, err
		}
		if err := validateModelFile(modelPath); err != nil {
			return nil, err
		}
		items = append(items, NewPythonDetector(item.Code, item.Description, item.ModelName, modelPath, runtimeRoot))
	}
	registry := &Registry{
		detectors: make(map[string]Detector),
	}
	for _, item := range items {
		if item == nil {
			continue
		}
		registry.detectors[item.Code()] = item
	}
	return registry, nil
}

// Get 获取原子检测能力
func (r *Registry) Get(code string) (Detector, error) {
	if r == nil {
		return nil, fmt.Errorf("unsupported module code: %s", code)
	}
	detector, ok := r.detectors[code]
	if !ok {
		return nil, fmt.Errorf("unsupported module code: %s", code)
	}
	return detector, nil
}

// FormatRegistryTable 格式化原子检测能力注册表
func FormatRegistryTable(registry *Registry) string {
	if registry == nil {
		return ""
	}
	taskCodes := make([]string, 0, len(registry.detectors))
	for taskCode := range registry.detectors {
		taskCodes = append(taskCodes, taskCode)
	}
	sort.Strings(taskCodes)
	recordsTaskCodes := make([]string, 0, len(taskCodes))
	recordsDescriptions := make([]string, 0, len(taskCodes))
	recordsModels := make([]string, 0, len(taskCodes))
	for _, taskCode := range taskCodes {
		detector := registry.detectors[taskCode]
		if detector == nil {
			continue
		}
		recordsTaskCodes = append(recordsTaskCodes, detector.Code())
		recordsDescriptions = append(recordsDescriptions, detector.Description())
		recordsModels = append(recordsModels, detector.ModelName())
	}
	return utils.FormatTable([]*utils.TableColumn{
		{
			Header: "detector",
			Values: recordsTaskCodes,
		},
		{
			Header: "description",
			Values: recordsDescriptions,
		},
		{
			Header: "model",
			Values: recordsModels,
		},
	})
}

// resolveModelPath 解析模型路径
func resolveModelPath(modelRoot, taskCode, modelName string) (string, error) {
	modelRoot = strings.TrimSpace(modelRoot)
	taskCode = strings.TrimSpace(taskCode)
	modelName = strings.TrimSpace(modelName)
	if modelRoot == "" {
		return "", fmt.Errorf("model root is required")
	}
	if taskCode == "" {
		return "", fmt.Errorf("task code is required")
	}
	if modelName == "" {
		return "", fmt.Errorf("model name is required")
	}
	return filepath.Join(absPath(modelRoot), taskCode, modelName), nil
}

// validateModelFile 校验模型文件存在
func validateModelFile(modelPath string) error {
	info, err := os.Stat(modelPath)
	if err != nil {
		return fmt.Errorf("model file not found: %s", modelPath)
	}
	if info.IsDir() {
		return fmt.Errorf("model path is a directory: %s", modelPath)
	}
	file, err := os.Open(modelPath)
	if err != nil {
		return fmt.Errorf("open model file failed: %s: %v", modelPath, err)
	}
	defer func() {
		_ = file.Close()
	}()

	buffer := make([]byte, 128)
	n, err := file.Read(buffer)
	if err != nil {
		return fmt.Errorf("read model file failed: %s: %v", modelPath, err)
	}
	if strings.HasPrefix(string(buffer[:n]), "version https://git-lfs.github.com/spec/v1") {
		return fmt.Errorf("model file is git lfs pointer: %s", modelPath)
	}

	return nil
}
