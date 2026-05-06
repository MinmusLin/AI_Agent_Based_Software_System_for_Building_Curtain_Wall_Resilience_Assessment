package dto

import (
	"icw_common/gen/core/api"
	"icw_common/gen/core/biz"
)

func NewCreateProjectGroupResponse(resp *bizpb.CreateProjectGroupResponse) *apipb.CreateProjectGroupResponse {
	if resp == nil {
		return nil
	}
	return &apipb.CreateProjectGroupResponse{
		Group: NewProjectGroup(resp.Group),
	}
}

func NewDeleteProjectGroupResponse(_ *bizpb.DeleteProjectGroupResponse) *apipb.DeleteProjectGroupResponse {
	return &apipb.DeleteProjectGroupResponse{}
}

func NewDeleteProjectImageResponse(_ *bizpb.DeleteProjectImageResponse) *apipb.DeleteProjectImageResponse {
	return &apipb.DeleteProjectImageResponse{}
}

func NewGetProjectAssetsResponse(resp *bizpb.GetProjectAssetsResponse) *apipb.GetProjectAssetsResponse {
	if resp == nil {
		return nil
	}
	return &apipb.GetProjectAssetsResponse{
		Groups: NewProjectGroups(resp.Groups),
	}
}

func NewGetProjectImageOriginalResponse(resp *bizpb.GetProjectImageOriginalResponse) *apipb.GetProjectImageOriginalResponse {
	if resp == nil {
		return nil
	}
	return &apipb.GetProjectImageOriginalResponse{
		OriginalUrl: resp.OriginalUrl,
	}
}

func NewMoveProjectGroupResponse(resp *bizpb.MoveProjectGroupResponse) *apipb.MoveProjectGroupResponse {
	if resp == nil {
		return nil
	}
	return &apipb.MoveProjectGroupResponse{
		Group: NewProjectGroup(resp.Group),
	}
}

func NewMoveProjectImageResponse(resp *bizpb.MoveProjectImageResponse) *apipb.MoveProjectImageResponse {
	if resp == nil {
		return nil
	}
	return &apipb.MoveProjectImageResponse{
		Images: resp.Images,
	}
}

func NewReportProjectImageResponse(_ *bizpb.ReportProjectImageResponse) *apipb.ReportProjectImageResponse {
	return &apipb.ReportProjectImageResponse{}
}

func NewUpdateProjectGroupResponse(resp *bizpb.UpdateProjectGroupResponse) *apipb.UpdateProjectGroupResponse {
	if resp == nil {
		return nil
	}
	return &apipb.UpdateProjectGroupResponse{
		Group: NewProjectGroup(resp.Group),
	}
}

func NewUploadProjectImageResponse(resp *bizpb.UploadProjectImageResponse) *apipb.UploadProjectImageResponse {
	if resp == nil {
		return nil
	}
	return &apipb.UploadProjectImageResponse{
		Images: resp.Images,
	}
}
