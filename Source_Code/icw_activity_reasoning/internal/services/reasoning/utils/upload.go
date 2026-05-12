package utils

import (
	"bytes"
	"context"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"net/textproto"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"sync"
	"time"

	"icw_common/gen/core/biz"
)

// UploadArtifacts 并发上传任务产物
func UploadArtifacts(ctx context.Context, policy *ArtifactUploadPolicy, taskDir string, timeout time.Duration) ([]*bizpb.ReasoningArtifactUploadResult, error) {
	if policy == nil {
		return nil, errors.New("reasoning artifact upload policy is nil")
	}

	artifactNames := listArtifactNames(taskDir)
	if len(artifactNames) == 0 {
		return make([]*bizpb.ReasoningArtifactUploadResult, 0), nil
	}

	results := make([]*bizpb.ReasoningArtifactUploadResult, len(artifactNames))

	wg := sync.WaitGroup{}
	for index, artifactName := range artifactNames {
		wg.Add(1)
		go func(resultIndex int, name string) {
			defer wg.Done()
			results[resultIndex] = uploadArtifact(ctx, policy, name, artifactPath(taskDir, name), timeout)
		}(index, artifactName)
	}
	wg.Wait()

	failedNames := make([]string, 0)
	for index, result := range results {
		name := artifactNames[index]
		if result != nil && strings.TrimSpace(result.Name) != "" {
			name = strings.TrimSpace(result.Name)
		}
		if result == nil || !result.Uploaded || strings.TrimSpace(result.Sha256) == "" {
			failedNames = append(failedNames, name)
		}
	}
	if len(failedNames) > 0 {
		return results, fmt.Errorf(
			"reasoning artifact upload failed: uploaded=%d total=%d failed=%s",
			len(results)-len(failedNames),
			len(results),
			strings.Join(failedNames, ","),
		)
	}

	return results, nil
}

// uploadArtifact 上传任务产物
func uploadArtifact(ctx context.Context, policy *ArtifactUploadPolicy, name, path string, timeout time.Duration) *bizpb.ReasoningArtifactUploadResult {
	result := &bizpb.ReasoningArtifactUploadResult{
		Name: name,
	}
	if policy == nil {
		return result
	}

	file, err := os.Open(path)
	if err != nil {
		return result
	}
	defer func() {
		_ = file.Close()
	}()

	uploadCtx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	hash := sha256.New()
	bodyReader := io.TeeReader(file, hash)
	var requestBody bytes.Buffer
	writer := multipart.NewWriter(&requestBody)

	for key, value := range policy.FormData {
		if key == "key" || key == "Content-Type" {
			continue
		}
		if err := writer.WriteField(key, value); err != nil {
			return result
		}
	}

	if err := writer.WriteField("key", policy.KeyPrefix+name); err != nil {
		return result
	}
	if err := writer.WriteField("Content-Type", "image/png"); err != nil {
		return result
	}

	header := make(textproto.MIMEHeader)
	header.Set("Content-Disposition", `form-data; name="file"; filename="`+name+`"`)
	header.Set("Content-Type", "image/png")

	part, err := writer.CreatePart(header)
	if err != nil {
		return result
	}
	if _, err := io.Copy(part, bodyReader); err != nil {
		return result
	}
	if err := writer.Close(); err != nil {
		return result
	}

	httpReq, err := http.NewRequestWithContext(uploadCtx, http.MethodPost, policy.URL, bytes.NewReader(requestBody.Bytes()))
	if err != nil {
		return result
	}
	httpReq.Header.Set("Content-Type", writer.FormDataContentType())

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

// listArtifactNames 查询需要上传的图像检测产物名称
func listArtifactNames(taskDir string) []string {
	entries, err := os.ReadDir(taskDir)
	if err != nil {
		return make([]string, 0)
	}
	names := make([]string, 0, len(entries))
	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}
		name := entry.Name()
		if strings.EqualFold(name, "report.json") {
			continue
		}
		if strings.ToLower(filepath.Ext(name)) != ".png" {
			continue
		}
		names = append(names, name)
	}
	sort.Strings(names)
	return names
}
