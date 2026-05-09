package dto

import (
	"icw_common/gen/core/api"
	"icw_common/gen/core/biz"
)

func NewGetProjectReportResponse(resp *bizpb.GetProjectReportResponse) *apipb.GetProjectReportResponse {
	if resp == nil {
		return nil
	}
	return &apipb.GetProjectReportResponse{
		Report: resp.Report,
	}
}
