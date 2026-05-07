package dto

import (
	"icw_common/gen/core/api"
	"icw_common/gen/core/biz"
)

func NewGetProjectDetectionTasksResponse(resp *bizpb.GetProjectDetectionTasksResponse) *apipb.GetProjectDetectionTasksResponse {
	if resp == nil {
		return nil
	}
	return &apipb.GetProjectDetectionTasksResponse{
		Tasks: resp.Tasks,
	}
}

func NewRetryProjectDetectionResponse(resp *bizpb.RetryProjectDetectionResponse) *apipb.RetryProjectDetectionResponse {
	if resp == nil {
		return nil
	}
	return &apipb.RetryProjectDetectionResponse{
		TaskCount: resp.TaskCount,
	}
}

func NewStartProjectDetectionResponse(resp *bizpb.StartProjectDetectionResponse) *apipb.StartProjectDetectionResponse {
	if resp == nil {
		return nil
	}
	return &apipb.StartProjectDetectionResponse{
		TaskCount: resp.TaskCount,
	}
}
