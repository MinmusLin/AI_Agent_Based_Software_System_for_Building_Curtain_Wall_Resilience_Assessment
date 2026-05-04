package events

import (
	"context"
	"time"

	"github.com/google/uuid"

	"icw_core_biz/pkg/dto/project"
	"icw_core_biz/repositories/rocketmq"
	"icw_core_biz/utils"
)

// PublishProjectImageStatusChangedEvent 发布项目图像状态变化事件
func PublishProjectImageStatusChangedEvent(ctx context.Context, rocketMQ *rocketmq.Repository, userId, projectId uint64, image *project.ProjectImage) {
	if image == nil {
		return
	}
	event := &project.ProjectImageStatusChangedEvent{
		EventId:     uuid.NewString(),
		EventType:   project.EventTypeProjectImageStatusChanged,
		ProjectId:   projectId,
		ProjectCode: utils.Encode(projectId),
		UserId:      userId,
		Image:       image,
		OccurredAt:  time.Now().Format("2006-01-02 15:04:05"),
	}

	// 发布项目图像状态变化事件
	_ = rocketMQ.PublishProjectImageStatusChangedEvent(ctx, event)
}
