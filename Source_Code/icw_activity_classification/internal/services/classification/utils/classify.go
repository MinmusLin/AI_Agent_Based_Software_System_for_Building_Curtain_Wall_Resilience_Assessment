package utils

import (
	"context"
	"errors"
	"math/rand"
	"sort"
	"strings"

	"icw_common/enum"
	"icw_common/gen/activity"
	"icw_common/gen/activity/classification"

	"icw_activity_classification/internal/agent"
)

func randomShowWithProbability(items []string) []string {
	if len(items) == 0 {
		return []string{}
	}
	r := rand.Intn(6)
	count := r
	if count > len(items) {
		count = len(items)
	}
	rand.Shuffle(len(items), func(i, j int) {
		items[i], items[j] = items[j], items[i]
	})
	return items[:count]
}

func ExecuteClassification(ctx context.Context, req *classificationpb.StartRequest, agentClient *agent.Client, imageSize int) ([]string, string, error) {
	//time.Sleep(time.Second)
	//return randomShowWithProbability([]string{"corrosion", "crack", "stain", "flatness", "spalling"}), "", nil
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

func ParseTaskCodes(output string) ([]string, error) {
	output = strings.TrimSpace(output)
	if output == "" {
		return nil, errors.New("classification output is empty")
	}
	return parseTaskCodesFromText(output)
}

func parseTaskCodesFromText(output string) ([]string, error) {
	normalizedOutput := strings.ToLower(output)
	candidates := []string{
		enum.DetectionTaskCodeString(activitypb.DetectionTaskCode_Corrosion),
		enum.DetectionTaskCodeString(activitypb.DetectionTaskCode_Crack),
		enum.DetectionTaskCodeString(activitypb.DetectionTaskCode_Stain),
		enum.DetectionTaskCodeString(activitypb.DetectionTaskCode_Flatness),
		enum.DetectionTaskCodeString(activitypb.DetectionTaskCode_Spalling),
	}
	codes := make([]string, 0, len(candidates))
	for _, candidate := range candidates {
		if strings.Contains(normalizedOutput, candidate) {
			codes = append(codes, candidate)
		}
	}
	return normalizeTaskCodes(codes)
}

func normalizeTaskCodes(codes []string) ([]string, error) {
	seen := map[string]struct{}{}
	normalizedCodes := make([]string, 0, len(codes))
	for _, code := range codes {
		code = strings.ToLower(strings.TrimSpace(code))
		if enum.ParseDetectionTaskCode(code) == activitypb.DetectionTaskCode_Unknown {
			continue
		}
		if _, ok := seen[code]; ok {
			continue
		}
		seen[code] = struct{}{}
		normalizedCodes = append(normalizedCodes, code)
	}
	if len(normalizedCodes) == 0 {
		return []string{}, nil
	}
	sort.Strings(normalizedCodes)
	return normalizedCodes, nil
}
