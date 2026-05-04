package project

// 事件类型
const (
	// EventTypeProjectImageStatusChanged 项目图像状态变化事件
	EventTypeProjectImageStatusChanged = "project_image_status_changed"
)

// 事件 Tag
const (
	// EventTagProjectImageStatusChanged 项目图像状态变化事件
	EventTagProjectImageStatusChanged = "PROJECT_IMAGE_STATUS_CHANGED"
)

// ProjectImageStatusChangedEvent 项目图像状态变化事件
type ProjectImageStatusChangedEvent struct {
	EventId     string
	EventType   string
	ProjectId   uint64
	ProjectCode string
	UserId      uint64
	Image       *ProjectImage
	OccurredAt  string
}
