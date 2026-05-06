package events

import (
	"context"
	"time"

	"github.com/google/uuid"

	"icw_common/consts"
	"icw_common/gen/core/biz"
	"icw_common/utils"
	"icw_core_biz/repositories/rocketmq"
)

// PublishProjectImageStatusChangedEvent 发布项目图像状态变化事件
func PublishProjectImageStatusChangedEvent(ctx context.Context, rocketMQ *rocketmq.Repository, userId, projectId uint64, image *bizpb.ProjectImage) {
	if image == nil {
		return
	}
	event := &bizpb.ProjectImageStatusChangedEvent{
		EventId:     uuid.NewString(),
		EventType:   consts.EventTypeProjectImageStatusChanged,
		ProjectId:   projectId,
		ProjectCode: utils.Encode(projectId),
		UserId:      userId,
		Image:       image,
		OccurredAt:  time.Now().Format("2006-01-02 15:04:05"),
	}

	// 发布项目图像状态变化事件
	_ = rocketMQ.PublishProjectImageStatusChangedEvent(ctx, event)
}
