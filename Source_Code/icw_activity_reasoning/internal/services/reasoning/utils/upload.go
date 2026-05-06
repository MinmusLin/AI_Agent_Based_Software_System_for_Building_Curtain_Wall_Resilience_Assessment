package utils

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"io"
	"net/http"
	"os"
	"strings"
	"time"

	"icw_common/gen/activity/reasoning"
	"icw_common/gen/core/biz"
)

// UploadArtifacts 按 BIZ 下发的预签名 URL 上传检测产物
func UploadArtifacts(ctx context.Context, plans []*reasoningpb.ReasoningArtifactUploadPlan, taskDir string, timeout time.Duration) []*bizpb.ReasoningArtifactUploadResult {
	results := make([]*bizpb.ReasoningArtifactUploadResult, 0, len(plans))
	for _, plan := range plans {
		if plan == nil {
			continue
		}
		results = append(results, uploadArtifact(ctx, plan, ArtifactPath(taskDir, plan.Name), timeout))
	}
	return results
}

// uploadArtifact 上传单个检测产物
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
	if strings.TrimSpace(plan.ContentType) != "" {
		httpReq.Header.Set("Content-Type", strings.TrimSpace(plan.ContentType))
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
