package events

import (
	"context"
	"time"

	"github.com/google/uuid"

	"icw_common/enum"
	"icw_common/gen/core/biz"
	"icw_common/gen/core/common"
	"icw_common/utils"

	"icw_core_biz/repositories/rocketmq"
)

const (
	// DetectionNodeStatusPending 检测任务节点状态：等待中/执行中
	DetectionNodeStatusPending = bizpb.ProjectDetectionSubTaskStatus_Pending
	// DetectionNodeStatusSucceeded 检测任务节点状态：成功
	DetectionNodeStatusSucceeded = bizpb.ProjectDetectionSubTaskStatus_Succeeded
	// DetectionNodeStatusFailed 检测任务节点状态：失败
	DetectionNodeStatusFailed = bizpb.ProjectDetectionSubTaskStatus_Failed
)

// ReasoningNodeCode 生成推理阶段检测任务节点代码
func ReasoningNodeCode(taskCode string) string {
	reasoningCode := enum.DetectionNodeCodeString(commonpb.DetectionNodeCode_Reasoning)
	if taskCode == "" {
		return reasoningCode
	}
	return reasoningCode + ":" + taskCode
}

// PublishProjectImageStatusChangedEvent 发布项目图像状态变化事件
func PublishProjectImageStatusChangedEvent(ctx context.Context, rocketMQ *rocketmq.Producer, userId, projectId uint64, image *bizpb.ProjectImage) {
	if image == nil {
		return
	}
	event := &bizpb.ProjectImageStatusChangedEvent{
		EventId:     uuid.NewString(),
		EventType:   commonpb.ProjectEventType_ImageStatusChanged,
		ProjectId:   projectId,
		ProjectCode: utils.Encode(projectId),
		UserId:      userId,
		Image:       image,
		OccurredAt:  time.Now().Format("2006-01-02 15:04:05"),
	}
	_ = rocketMQ.PublishProjectImageStatusChangedEvent(ctx, event)
}

// PublishProjectDetectionNodeStatusChangedEvent 发布项目图像检测任务状态变化事件
func PublishProjectDetectionNodeStatusChangedEvent(ctx context.Context, rocketMQ *rocketmq.Producer, userId, projectId uint64, imageUuid, nodeCode, mainTaskUuid string, mainStatus bizpb.ProjectDetectionTaskStatus_Value, subTaskUuid string, subStatus bizpb.ProjectDetectionSubTaskStatus_Value) {
	event := &bizpb.ProjectDetectionTaskStatusChangedEvent{
		EventId:      uuid.NewString(),
		EventType:    commonpb.ProjectEventType_DetectionTaskStatusChanged,
		ProjectId:    projectId,
		ProjectCode:  utils.Encode(projectId),
		UserId:       userId,
		ImageUuid:    imageUuid,
		NodeCode:     nodeCode,
		MainTaskUuid: mainTaskUuid,
		MainStatus:   mainStatus,
		SubTaskUuid:  subTaskUuid,
		SubStatus:    subStatus,
		OccurredAt:   time.Now().Format("2006-01-02 15:04:05"),
	}
	_ = rocketMQ.PublishProjectDetectionTaskStatusChangedEvent(ctx, event)
}

// PublishProjectReportStatusChangedEvent 发布项目评估报告状态变化事件
func PublishProjectReportStatusChangedEvent(ctx context.Context, rocketMQ *rocketmq.Producer, userId, projectId uint64, reportUuid string, status bizpb.ProjectReportStatus_Value) {
	event := &bizpb.ProjectReportStatusChangedEvent{
		EventId:     uuid.NewString(),
		EventType:   commonpb.ProjectEventType_ReportStatusChanged,
		ProjectId:   projectId,
		ProjectCode: utils.Encode(projectId),
		UserId:      userId,
		ReportUuid:  reportUuid,
		Status:      status,
		OccurredAt:  time.Now().Format("2006-01-02 15:04:05"),
	}
	_ = rocketMQ.PublishProjectReportStatusChangedEvent(ctx, event)
}
