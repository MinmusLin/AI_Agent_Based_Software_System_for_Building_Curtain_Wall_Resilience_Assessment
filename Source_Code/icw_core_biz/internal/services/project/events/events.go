package events

import (
	"context"
	"time"

	"github.com/google/uuid"

	"icw_common/consts"
	"icw_common/enum"
	"icw_common/gen/core/biz"
	"icw_common/utils"
	"icw_core_biz/repositories/rocketmq"
)

const (
	// DetectionNodeCodeClassification 检测任务节点代码：分类阶段
	DetectionNodeCodeClassification = "classification"
	// DetectionNodeCodeReasoning 检测任务节点代码：推理阶段
	DetectionNodeCodeReasoning = "reasoning"
	// DetectionNodeCodeSummary 检测任务节点代码：总结阶段
	DetectionNodeCodeSummary = "summary"
)

var (
	// DetectionNodeStatusPending 检测任务节点状态：等待中/执行中
	DetectionNodeStatusPending = enum.ProjectDetectionSubTaskStatusString(bizpb.ProjectDetectionSubTaskStatus_Pending)
	// DetectionNodeStatusSucceeded 检测任务节点状态：成功
	DetectionNodeStatusSucceeded = enum.ProjectDetectionSubTaskStatusString(bizpb.ProjectDetectionSubTaskStatus_Succeeded)
	// DetectionNodeStatusFailed 检测任务节点状态：失败
	DetectionNodeStatusFailed = enum.ProjectDetectionSubTaskStatusString(bizpb.ProjectDetectionSubTaskStatus_Failed)
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
	_ = rocketMQ.PublishProjectImageStatusChangedEvent(ctx, event)
}

// PublishProjectDetectionNodeStatusChangedEvent 发布项目图像检测任务状态变化事件
func PublishProjectDetectionNodeStatusChangedEvent(ctx context.Context, rocketMQ *rocketmq.Repository, userId, projectId uint64, imageUuid, nodeCode, mainTaskId, mainStatus, subTaskId, subStatus string) {
	event := &bizpb.ProjectDetectionTaskStatusChangedEvent{
		EventId:     uuid.NewString(),
		EventType:   consts.EventTypeProjectDetectionTaskStatusChanged,
		ProjectId:   projectId,
		ProjectCode: utils.Encode(projectId),
		UserId:      userId,
		ImageUuid:   imageUuid,
		NodeCode:    nodeCode,
		MainTaskId:  mainTaskId,
		MainStatus:  mainStatus,
		SubTaskId:   subTaskId,
		SubStatus:   subStatus,
		OccurredAt:  time.Now().Format("2006-01-02 15:04:05"),
	}
	_ = rocketMQ.PublishProjectDetectionTaskStatusChangedEvent(ctx, event)
}

// ReasoningNodeCode 生成推理阶段检测任务节点代码
func ReasoningNodeCode(taskCode string) string {
	if taskCode == "" {
		return DetectionNodeCodeReasoning
	}
	return DetectionNodeCodeReasoning + ":" + taskCode
}
