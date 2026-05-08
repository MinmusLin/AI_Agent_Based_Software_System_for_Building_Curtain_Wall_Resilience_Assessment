package utils

import (
	"bytes"
	"context"
	"image"
	"image/color"
	"image/jpeg"
	"image/png"
	"io"
	"net/http"
	"strconv"

	"golang.org/x/image/draw"
)

const (
	classificationImageContentType = "image/png"
)

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
		return nil, "", errUnexpectedStatus(response.StatusCode)
	}
	source, _, err := image.Decode(io.LimitReader(response.Body, 32<<20))
	if err != nil {
		return nil, "", err
	}
	thumbnail := squareThumbnail(source, size)
	var buffer bytes.Buffer
	if err := png.Encode(&buffer, thumbnail); err != nil {
		return nil, "", err
	}
	return buffer.Bytes(), classificationImageContentType, nil
}

func squareThumbnail(source image.Image, size int) image.Image {
	bounds := source.Bounds()
	sourceWidth := bounds.Dx()
	sourceHeight := bounds.Dy()
	if sourceWidth <= 0 || sourceHeight <= 0 || size <= 0 {
		return image.NewRGBA(image.Rect(0, 0, size, size))
	}

	destination := image.NewRGBA(image.Rect(0, 0, size, size))
	draw.Draw(destination, destination.Bounds(), &image.Uniform{C: color.White}, image.Point{}, draw.Src)

	scale := float64(size) / float64(max(sourceWidth, sourceHeight))
	targetWidth := max(1, int(float64(sourceWidth)*scale))
	targetHeight := max(1, int(float64(sourceHeight)*scale))
	targetRect := image.Rect(
		(size-targetWidth)/2,
		(size-targetHeight)/2,
		(size+targetWidth)/2,
		(size+targetHeight)/2,
	)
	draw.CatmullRom.Scale(destination, targetRect, source, bounds, draw.Over, nil)
	return destination
}

func errUnexpectedStatus(statusCode int) error {
	return &unexpectedStatusError{statusCode: statusCode}
}

type unexpectedStatusError struct {
	statusCode int
}

func (e *unexpectedStatusError) Error() string {
	statusText := http.StatusText(e.statusCode)
	if statusText == "" {
		statusText = strconv.Itoa(e.statusCode)
	}
	return "unexpected http status: " + statusText
}

func init() {
	image.RegisterFormat("jpeg", "\xff\xd8", jpeg.Decode, jpeg.DecodeConfig)
	image.RegisterFormat("png", "\x89PNG\r\n\x1a\n", png.Decode, png.DecodeConfig)
}
