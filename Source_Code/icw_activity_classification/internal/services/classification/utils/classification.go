package utils

import (
	"context"
	"encoding/json"
	"errors"
	"sort"
	"strings"

	"icw_common/enum"
	"icw_common/gen/activity/classification"
	"icw_common/gen/core/common"
	"icw_common/utils"

	"icw_activity_classification/internal/agent"
)

// ExecuteClassification 执行图像分类任务并返回检测任务代码列表和模型原始输出
func ExecuteClassification(ctx context.Context, req *classificationpb.StartRequest, client *agent.Client, imageSize int) ([]string, string, error) {
	imageBytes, contentType, err := DownloadAndResizeImage(ctx, req.PresignGetUrl, imageSize)
	if err != nil {
		return nil, "", err
	}

	output, err := client.Chat(ctx, agent.Message{
		Text:        "图像已上传，请根据指令进行符合要求的输出。",
		Image:       imageBytes,
		ContentType: contentType,
	})
	if err != nil {
		return nil, output, err
	}

	normalizedOutput, err := utils.CompactAgentJSONObjectString(output)
	if err != nil {
		return nil, output, err
	}
	taskCodes, err := parseTaskCodes(normalizedOutput)

	return taskCodes, normalizedOutput, err
}

// parseTaskCodes 从模型 JSON 输出中解析检测任务代码列表
func parseTaskCodes(output string) ([]string, error) {
	output = strings.TrimSpace(output)
	if output == "" {
		return nil, errors.New("classification output is required")
	}
	result := struct {
		TaskCodes []string `json:"task_codes"`
	}{}
	if err := json.Unmarshal([]byte(output), &result); err != nil {
		return nil, err
	}
	return normalizeTaskCodes(result.TaskCodes)
}

// normalizeTaskCodes 标准化并去重检测任务代码
func normalizeTaskCodes(taskCodes []string) ([]string, error) {
	seen := map[commonpb.DetectionTaskCode_Value]struct{}{}
	codes := make([]string, 0, len(taskCodes))
	for _, taskCode := range taskCodes {
		code := enum.ParseDetectionTaskCode(taskCode)
		if code == commonpb.DetectionTaskCode_Unknown {
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
