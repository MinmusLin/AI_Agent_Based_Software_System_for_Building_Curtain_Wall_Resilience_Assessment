package model

// ProjectAssetsReadyStats 项目图像状态校验
type ProjectAssetsReadyStats struct {
	PendingImageCount  uint64
	UploadedImageCount uint64
	FailedImageCount   uint64
	EmptyGroupCount    uint64
}
