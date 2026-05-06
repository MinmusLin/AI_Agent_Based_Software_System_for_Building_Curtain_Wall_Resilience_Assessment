package dto

import (
	"icw_common/gen/core/api"
	"icw_common/gen/core/biz"
)

// NewCreateSocketTicketResponse 将 BIZ WebSocket 票据响应转换为 API WebSocket 票据响应
func NewCreateSocketTicketResponse(resp *bizpb.CreateSocketTicketResponse) *apipb.CreateSocketTicketResponse {
	if resp == nil {
		return nil
	}
	return &apipb.CreateSocketTicketResponse{
		Ticket:    resp.Ticket,
		ExpiresIn: resp.ExpiresIn,
	}
}

// NewProjectImageStatusChangedMessage 将 BIZ 图像状态变化事件转换为 API WebSocket 消息
func NewProjectImageStatusChangedMessage(event *bizpb.ProjectImageStatusChangedEvent) *apipb.ProjectImageStatusChangedMessage {
	if event == nil {
		return nil
	}
	return &apipb.ProjectImageStatusChangedMessage{
		Type:      event.EventType,
		ProjectId: event.ProjectCode,
		Image:     event.Image,
	}
}
