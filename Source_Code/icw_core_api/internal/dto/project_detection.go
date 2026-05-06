package dto

import (
	"icw_common/gen/core/api"
	"icw_common/gen/core/biz"
)

func NewStartProjectDetectionResponse(resp *bizpb.StartProjectDetectionResponse) *apipb.StartProjectDetectionResponse {
	if resp == nil {
		return nil
	}
	return &apipb.StartProjectDetectionResponse{
		TaskCount: resp.TaskCount,
	}
}
