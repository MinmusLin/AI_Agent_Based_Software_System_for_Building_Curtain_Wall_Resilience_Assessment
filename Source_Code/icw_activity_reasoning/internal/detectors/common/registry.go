package common

import (
	"context"
	"fmt"
	"path/filepath"
	"sort"

	"icw_common/utils"
)

// Detector 原子检测能力
type Detector interface {
	Code() string
	Description() string
	Detect(ctx context.Context, imageUuid string) error
}

// DetectorMeta 原子检测能力元数据
type DetectorMeta struct {
	Code        string
	Description string
	Path        string
}

func NewDetectorMeta(code, description string) *DetectorMeta {
	return &DetectorMeta{
		Code:        code,
		Description: description,
		Path:        filepath.Join("python/detectors", code),
	}
}

// Registry 原子检测能力注册表
type Registry struct {
	detectors map[string]Detector
}

func NewRegistry(runtimeRoot string, metas []*DetectorMeta) *Registry {
	items := make([]Detector, 0, len(metas))
	for _, item := range metas {
		if item == nil {
			continue
		}
		items = append(items, NewPythonDetector(item.Code, item.Description, item.Path, runtimeRoot))
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
	return registry
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
	for _, taskCode := range taskCodes {
		detector := registry.detectors[taskCode]
		if detector == nil {
			continue
		}
		recordsTaskCodes = append(recordsTaskCodes, detector.Code())
		recordsDescriptions = append(recordsDescriptions, detector.Description())
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
	})
}
