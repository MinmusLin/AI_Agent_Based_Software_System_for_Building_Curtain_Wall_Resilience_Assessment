package assets

import (
	"log"
	"time"

	"github.com/google/uuid"

	"icw_core_biz/internal/services/common"
	"icw_core_biz/pkg/dto/project"
	"icw_core_biz/utils"
)

// publishProjectImageStatusChangedEvent 发布项目图像状态变化事件
func (s *Service) publishProjectImageStatusChangedEvent(userId, projectId uint64, image *project.ProjectImage) {
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
	if err := s.RocketMQ().PublishProjectImageStatusChangedEvent(s.Ctx(), event); err != nil {
		log.Printf("%s Publish project image status event failed, project_id: %d, image_uuid: %s, err: %v", common.RpcWarnPrefix(), projectId, image.Uuid, err)
	}
}
