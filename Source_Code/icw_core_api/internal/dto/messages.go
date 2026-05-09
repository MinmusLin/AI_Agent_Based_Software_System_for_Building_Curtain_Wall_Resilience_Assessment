package dto

import (
	"icw_common/gen/core/api"
	"icw_common/gen/core/biz"
)

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

func NewProjectDetectionTaskStatusChangedMessage(event *bizpb.ProjectDetectionTaskStatusChangedEvent) *apipb.ProjectDetectionTaskStatusChangedMessage {
	if event == nil {
		return nil
	}
	return &apipb.ProjectDetectionTaskStatusChangedMessage{
		Type:         event.EventType,
		ProjectId:    event.ProjectCode,
		ImageUuid:    event.ImageUuid,
		NodeCode:     event.NodeCode,
		MainTaskUuid: event.MainTaskUuid,
		MainStatus:   event.MainStatus,
		SubTaskUuid:  event.SubTaskUuid,
		SubStatus:    event.SubStatus,
		OccurredAt:   event.OccurredAt,
	}
}

func NewProjectReportStatusChangedMessage(event *bizpb.ProjectReportStatusChangedEvent) *apipb.ProjectReportStatusChangedMessage {
	if event == nil {
		return nil
	}
	return &apipb.ProjectReportStatusChangedMessage{
		Type:       event.EventType,
		ProjectId:  event.ProjectCode,
		ReportUuid: event.ReportUuid,
		Status:     event.Status,
		OccurredAt: event.OccurredAt,
	}
}
