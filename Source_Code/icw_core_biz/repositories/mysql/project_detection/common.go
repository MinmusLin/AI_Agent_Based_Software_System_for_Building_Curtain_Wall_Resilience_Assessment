package project_detection

import (
	"strings"

	"icw_common/enum"
	"icw_common/gen/core/biz"
	"icw_common/gen/core/common"

	"icw_core_biz/repositories/mysql/model"
)

// NormalizeDetectionTaskCodes 标准化图像检测子任务代码
func NormalizeDetectionTaskCodes(taskCodes []string) ([]string, error) {
	normalized := make([]string, 0, len(taskCodes))
	seen := map[commonpb.DetectionTaskCode_Value]bool{}
	for _, taskCode := range taskCodes {
		if strings.TrimSpace(taskCode) == "" {
			continue
		}
		code := enum.ParseDetectionTaskCode(taskCode)
		if code == commonpb.DetectionTaskCode_Unknown {
			return nil, model.ErrUnsupportedDetectionTaskCode
		}
		if seen[code] {
			continue
		}
		normalized = append(normalized, enum.DetectionTaskCodeString(code))
		seen[code] = true
	}
	return normalized, nil
}

// ClassificationNodeStatus 根据图像检测主任务状态推导图像检测分类任务状态
func ClassificationNodeStatus(task *model.ProjectDetectionTaskRecord) bizpb.ProjectDetectionSubTaskStatus_Value {
	switch task.Status {
	case bizpb.ProjectDetectionTaskStatus_Classifying:
		return bizpb.ProjectDetectionSubTaskStatus_Pending
	case bizpb.ProjectDetectionTaskStatus_Detecting,
		bizpb.ProjectDetectionTaskStatus_Summarizing,
		bizpb.ProjectDetectionTaskStatus_Succeeded:
		return bizpb.ProjectDetectionSubTaskStatus_Succeeded
	case bizpb.ProjectDetectionTaskStatus_Failed:
		if !task.CorrosionShouldExecute &&
			!task.CrackShouldExecute &&
			!task.StainShouldExecute &&
			!task.FlatnessShouldExecute &&
			!task.SpallingShouldExecute &&
			!task.SummaryShouldExecute {
			return bizpb.ProjectDetectionSubTaskStatus_Failed
		}
		return bizpb.ProjectDetectionSubTaskStatus_Succeeded
	default:
		return bizpb.ProjectDetectionSubTaskStatus_Unknown
	}
}
