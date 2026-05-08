package utils

import (
	"context"
	"encoding/json"
	"errors"
	"sort"
	"strings"

	"icw_common/enum"
	"icw_common/gen/activity"
	"icw_common/gen/activity/classification"

	"icw_activity_classification/internal/agent"
)

// ExecuteClassification 执行图像分类并返回任务代码列表和模型原始输出
func ExecuteClassification(ctx context.Context, req *classificationpb.StartRequest, agentClient *agent.Client, imageSize int) ([]string, string, error) {
	imageBytes, contentType, err := DownloadAndResizeImage(ctx, req.PresignGetUrl, imageSize)
	if err != nil {
		return nil, "", err
	}
	output, err := agentClient.Chat(ctx, agent.Message{
		Text:        "图像已上传，请执行分类",
		Image:       imageBytes,
		ContentType: contentType,
	})
	if err != nil {
		return nil, output, err
	}
	taskCodes, err := ParseTaskCodes(output)
	return taskCodes, output, err
}

type classificationOutput struct {
	TaskCodes []string `json:"task_codes"`
}

// ParseTaskCodes 从模型 JSON 输出中解析任务代码列表
func ParseTaskCodes(output string) ([]string, error) {
	output = strings.TrimSpace(output)
	if output == "" {
		return nil, errors.New("classification output is empty")
	}
	result := classificationOutput{}
	if err := json.Unmarshal([]byte(output), &result); err != nil {
		return nil, err
	}
	return normalizeTaskCodes(result.TaskCodes)
}

// normalizeTaskCodes 标准化并去重任务代码
func normalizeTaskCodes(taskCodes []string) ([]string, error) {
	seen := map[activitypb.DetectionTaskCode_Value]struct{}{}
	codes := make([]string, 0, len(taskCodes))
	for _, taskCode := range taskCodes {
		code := enum.ParseDetectionTaskCode(taskCode)
		if code == activitypb.DetectionTaskCode_Unknown {
			return nil, errors.New("classification task code is invalid")
		}
		if _, ok := seen[code]; ok {
			continue
		}
		seen[code] = struct{}{}
		codes = append(codes, enum.DetectionTaskCodeString(code))
	}
	if len(codes) == 0 {
		return []string{}, nil
	}
	sort.Strings(codes)
	return codes, nil
}
