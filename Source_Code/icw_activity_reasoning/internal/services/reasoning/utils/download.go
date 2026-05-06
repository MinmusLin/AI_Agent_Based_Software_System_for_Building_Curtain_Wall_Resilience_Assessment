package utils

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"icw_common/gen/activity/reasoning"
)

// DownloadSourceImage 下载待检测原图到任务临时目录
func DownloadSourceImage(ctx context.Context, req *reasoningpb.StartRequest, taskDir string, timeout time.Duration) error {
	downloadCtx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	httpReq, err := http.NewRequestWithContext(downloadCtx, http.MethodGet, req.PresignGetUrl, nil)
	if err != nil {
		return err
	}
	resp, err := http.DefaultClient.Do(httpReq)
	if err != nil {
		return err
	}
	defer func() {
		_ = resp.Body.Close()
	}()
	if resp.StatusCode < http.StatusOK || resp.StatusCode >= http.StatusMultipleChoices {
		return fmt.Errorf("download original image failed: status=%d", resp.StatusCode)
	}

	imagePath := filepath.Join(taskDir, "original.png")
	file, err := os.Create(imagePath)
	if err != nil {
		return err
	}
	defer func() {
		_ = file.Close()
	}()
	if _, err := io.Copy(file, resp.Body); err != nil {
		return err
	}
	return nil
}
