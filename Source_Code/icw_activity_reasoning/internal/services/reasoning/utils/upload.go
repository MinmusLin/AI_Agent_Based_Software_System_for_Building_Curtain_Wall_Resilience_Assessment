package utils

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"io"
	"net/http"
	"os"
	"sync"
	"time"

	"icw_common/gen/activity/reasoning"
	"icw_common/gen/core/biz"
)

// UploadArtifacts 并发上传任务产物
func UploadArtifacts(ctx context.Context, plans []*reasoningpb.ReasoningArtifactUploadPlan, taskDir string, timeout time.Duration) []*bizpb.ReasoningArtifactUploadResult {
	uploadPlans := make([]*reasoningpb.ReasoningArtifactUploadPlan, 0, len(plans))
	for _, plan := range plans {
		if plan != nil {
			uploadPlans = append(uploadPlans, plan)
		}
	}

	results := make([]*bizpb.ReasoningArtifactUploadResult, len(uploadPlans))

	wg := sync.WaitGroup{}
	for index, plan := range uploadPlans {
		wg.Add(1)
		go func(resultIndex int, item *reasoningpb.ReasoningArtifactUploadPlan) {
			defer wg.Done()
			results[resultIndex] = uploadArtifact(ctx, item, artifactPath(taskDir, item.Name), timeout)
		}(index, plan)
	}
	wg.Wait()

	return results
}

// uploadArtifact 上传任务产物
func uploadArtifact(ctx context.Context, plan *reasoningpb.ReasoningArtifactUploadPlan, path string, timeout time.Duration) *bizpb.ReasoningArtifactUploadResult {
	result := &bizpb.ReasoningArtifactUploadResult{
		Name: plan.Name,
	}

	file, err := os.Open(path)
	if err != nil {
		return result
	}
	defer func() {
		_ = file.Close()
	}()

	info, err := file.Stat()
	if err != nil {
		return result
	}

	uploadCtx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	hash := sha256.New()
	body := io.TeeReader(file, hash)
	httpReq, err := http.NewRequestWithContext(uploadCtx, http.MethodPut, plan.PresignUploadUrl, body)
	if err != nil {
		return result
	}

	httpReq.ContentLength = info.Size()
	if plan.ContentType != "" {
		httpReq.Header.Set("Content-Type", plan.ContentType)
	}

	resp, err := http.DefaultClient.Do(httpReq)
	if err != nil {
		return result
	}
	defer func() {
		_ = resp.Body.Close()
	}()

	if resp.StatusCode < http.StatusOK || resp.StatusCode >= http.StatusMultipleChoices {
		return result
	}

	result.Uploaded = true
	result.Sha256 = hex.EncodeToString(hash.Sum(nil))

	return result
}
