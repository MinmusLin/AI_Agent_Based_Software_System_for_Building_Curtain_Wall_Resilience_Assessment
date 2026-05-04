package dto

import (
	"icw_core_api/internal/dto/project"
	bizDto "icw_core_biz/pkg/dto"
	bizProjectDto "icw_core_biz/pkg/dto/project"
)

const (
	// SocketScopeProjectAssets 图像资产 WebSocket 连接范围
	SocketScopeProjectAssets = "ws_project_assets"
)

type CreateSocketTicketRequest struct {
	ProjectId string `json:"project_id"`
}

type CreateSocketTicketResponse struct {
	Ticket    string `json:"ticket"`
	ExpiresIn int64  `json:"expires_in"`
}

func NewCreateSocketTicketResponse(resp *bizDto.CreateSocketTicketResponse) *CreateSocketTicketResponse {
	if resp == nil {
		return nil
	}
	return &CreateSocketTicketResponse{
		Ticket:    resp.Ticket,
		ExpiresIn: resp.ExpiresIn,
	}
}

type ProjectImageStatusChangedMessage struct {
	Type      string                `json:"type"`
	ProjectId string                `json:"project_id"`
	Image     *project.ProjectImage `json:"image"`
}

func NewProjectImageStatusChangedMessage(event *bizProjectDto.ProjectImageStatusChangedEvent) *ProjectImageStatusChangedMessage {
	if event == nil {
		return nil
	}
	return &ProjectImageStatusChangedMessage{
		Type:      event.EventType,
		ProjectId: event.ProjectCode,
		Image:     project.NewProjectImage(event.Image),
	}
}
