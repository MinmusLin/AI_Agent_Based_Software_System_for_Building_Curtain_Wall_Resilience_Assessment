package dto

import (
	"icw_common/gen/core/api"
	"icw_common/gen/core/biz"
)

func NewDeleteProjectThumbnailResponse(_ *bizpb.DeleteProjectThumbnailResponse) *apipb.DeleteProjectThumbnailResponse {
	return &apipb.DeleteProjectThumbnailResponse{}
}

func NewGetProjectProfileResponse(resp *bizpb.GetProjectProfileResponse) *apipb.GetProjectProfileResponse {
	if resp == nil {
		return nil
	}
	return &apipb.GetProjectProfileResponse{
		Project: NewProject(resp.Project),
	}
}

func NewGetProjectThumbnailResponse(resp *bizpb.GetProjectThumbnailResponse) *apipb.GetProjectThumbnailResponse {
	if resp == nil {
		return nil
	}
	return &apipb.GetProjectThumbnailResponse{
		ThumbnailUrl: resp.ThumbnailUrl,
	}
}

func NewUpdateProjectProfileResponse(resp *bizpb.UpdateProjectProfileResponse) *apipb.UpdateProjectProfileResponse {
	if resp == nil {
		return nil
	}
	return &apipb.UpdateProjectProfileResponse{
		Project: NewProject(resp.Project),
	}
}

func NewUploadProjectThumbnailResponse(resp *bizpb.UploadProjectThumbnailResponse) *apipb.UploadProjectThumbnailResponse {
	if resp == nil {
		return nil
	}
	return &apipb.UploadProjectThumbnailResponse{
		UploadUrl: resp.UploadUrl,
	}
}
