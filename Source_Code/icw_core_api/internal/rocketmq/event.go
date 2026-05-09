package rocketmq

import (
	"encoding/json"
	"errors"
	"fmt"

	"github.com/apache/rocketmq-client-go/v2/primitive"

	"icw_common/consts"
	"icw_common/gen/core/biz"

	"icw_core_api/internal/dto"
	"icw_core_api/internal/socket"
)

// dispatchProjectEvent 分发项目事件
func dispatchProjectEvent(hub *socket.Hub, message *primitive.MessageExt) error {
	if hub == nil {
		return errors.New("websocket hub is nil")
	}
	switch message.GetTags() {
	case consts.EventTagProjectImageStatusChanged:
		return dispatchProjectImageStatusChangedEvent(hub, message)
	case consts.EventTagProjectDetectionTaskStatusChanged:
		return dispatchProjectDetectionTaskStatusChangedEvent(hub, message)
	case consts.EventTagProjectReportStatusChanged:
		return dispatchProjectReportStatusChangedEvent(hub, message)
	default:
		return fmt.Errorf("dispatch project event failed: invalid tag %v", message.GetTags())
	}
}

// dispatchProjectImageStatusChangedEvent 分发项目图像状态变化事件
func dispatchProjectImageStatusChangedEvent(hub *socket.Hub, message *primitive.MessageExt) error {
	var event bizpb.ProjectImageStatusChangedEvent
	if err := json.Unmarshal(message.Body, &event); err != nil {
		return err
	}
	if event.EventType != consts.EventTypeProjectImageStatusChanged {
		return fmt.Errorf("unexpected event type: %s", event.EventType)
	}

	socketMessage := dto.NewProjectImageStatusChangedMessage(&event)
	messageBytes, err := json.Marshal(socketMessage)
	if err != nil {
		return err
	}

	hub.BroadcastProject(event.ProjectId, event.ProjectCode, consts.SocketScopeProjectAssets, messageBytes)
	return nil
}

// dispatchProjectDetectionTaskStatusChangedEvent 分发项目图像检测任务状态变化事件
func dispatchProjectDetectionTaskStatusChangedEvent(hub *socket.Hub, message *primitive.MessageExt) error {
	var event bizpb.ProjectDetectionTaskStatusChangedEvent
	if err := json.Unmarshal(message.Body, &event); err != nil {
		return err
	}
	if event.EventType != consts.EventTypeProjectDetectionTaskStatusChanged {
		return fmt.Errorf("unexpected event type: %s", event.EventType)
	}

	socketMessage := dto.NewProjectDetectionTaskStatusChangedMessage(&event)
	messageBytes, err := json.Marshal(socketMessage)
	if err != nil {
		return err
	}

	hub.BroadcastProject(event.ProjectId, event.ProjectCode, consts.SocketScopeProjectDetection, messageBytes)
	return nil
}

// dispatchProjectReportStatusChangedEvent 分发项目评估报告状态变化事件
func dispatchProjectReportStatusChangedEvent(hub *socket.Hub, message *primitive.MessageExt) error {
	var event bizpb.ProjectReportStatusChangedEvent
	if err := json.Unmarshal(message.Body, &event); err != nil {
		return err
	}
	if event.EventType != consts.EventTypeProjectReportStatusChanged {
		return fmt.Errorf("unexpected event type: %s", event.EventType)
	}

	socketMessage := dto.NewProjectReportStatusChangedMessage(&event)
	messageBytes, err := json.Marshal(socketMessage)
	if err != nil {
		return err
	}

	hub.BroadcastProject(event.ProjectId, event.ProjectCode, consts.SocketScopeProjectReport, messageBytes)
	return nil
}
