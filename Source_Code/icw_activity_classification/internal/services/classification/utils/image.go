package utils

import (
	"bytes"
	"context"
	"fmt"
	"image"
	"image/png"
	"io"
	"net/http"

	"golang.org/x/image/draw"

	"icw_activity_classification/utils"
)

const (
	// ClassificationImageContentType 项目图像 MIME 类型
	ClassificationImageContentType = "image/png"
)

// init 注册 png 图像解码格式
func init() {
	image.RegisterFormat("png", "\x89PNG\r\n\x1a\n", png.Decode, png.DecodeConfig)
}

// DownloadAndResizeImage 下载并缩放项目图像
func DownloadAndResizeImage(ctx context.Context, imageURL string, size int) ([]byte, string, error) {
	request, err := http.NewRequestWithContext(ctx, http.MethodGet, imageURL, nil)
	if err != nil {
		return nil, "", err
	}

	response, err := http.DefaultClient.Do(request)
	if err != nil {
		return nil, "", err
	}

	defer func() {
		_ = response.Body.Close()
	}()

	if response.StatusCode < http.StatusOK || response.StatusCode >= http.StatusMultipleChoices {
		return nil, "", fmt.Errorf("unexpected http status code %d: %s", response.StatusCode, http.StatusText(response.StatusCode))
	}

	data, err := io.ReadAll(io.LimitReader(response.Body, 32<<20))
	if err != nil {
		return nil, "", err
	}

	config, _, err := image.DecodeConfig(bytes.NewReader(data))
	if err != nil {
		return nil, "", err
	}

	if config.Width <= size || config.Height <= size || size <= 0 {
		return data, utils.FirstNotEmpty(response.Header.Get("Content-Type"), http.DetectContentType(data)), nil
	}

	source, _, err := image.Decode(bytes.NewReader(data))
	if err != nil {
		return nil, "", err
	}

	target := resizeImage(source, size)
	var buffer bytes.Buffer

	if err := png.Encode(&buffer, target); err != nil {
		return nil, "", err
	}

	return buffer.Bytes(), ClassificationImageContentType, nil
}

// resizeImage 缩放项目图像
func resizeImage(source image.Image, size int) image.Image {
	bounds := source.Bounds()
	sourceWidth := bounds.Dx()
	sourceHeight := bounds.Dy()

	if sourceWidth <= 0 || sourceHeight <= 0 || size <= 0 {
		return source
	}
	if sourceWidth <= size || sourceHeight <= size {
		return source
	}

	scale := float64(size) / float64(min(sourceWidth, sourceHeight))
	targetWidth := max(1, int(float64(sourceWidth)*scale))
	targetHeight := max(1, int(float64(sourceHeight)*scale))
	destination := image.NewRGBA(image.Rect(0, 0, targetWidth, targetHeight))

	draw.CatmullRom.Scale(destination, destination.Bounds(), source, bounds, draw.Over, nil)

	return destination
}
