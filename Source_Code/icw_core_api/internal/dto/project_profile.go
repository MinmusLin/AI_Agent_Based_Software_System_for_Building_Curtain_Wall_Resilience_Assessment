package dto

import (
	"icw_common/gen/core/api"
	"icw_common/gen/core/biz"
)

// NewGetProjectProfileResponse 将 BIZ 项目基础信息响应转换为 API 项目基础信息响应
func NewGetProjectProfileResponse(resp *bizpb.GetProjectProfileResponse) *apipb.GetProjectProfileResponse {
	if resp == nil {
		return nil
	}
	return &apipb.GetProjectProfileResponse{
		Project: NewProject(resp.Project),
	}
}

// NewGetProjectThumbnailResponse 将 BIZ 项目缩略图响应转换为 API 项目缩略图响应
func NewGetProjectThumbnailResponse(resp *bizpb.GetProjectThumbnailResponse) *apipb.GetProjectThumbnailResponse {
	if resp == nil {
		return nil
	}
	return &apipb.GetProjectThumbnailResponse{
		ThumbnailUrl: resp.ThumbnailUrl,
	}
}

// NewUpdateProjectProfileResponse 将 BIZ 更新项目基础信息响应转换为 API 更新项目基础信息响应
func NewUpdateProjectProfileResponse(resp *bizpb.UpdateProjectProfileResponse) *apipb.UpdateProjectProfileResponse {
	if resp == nil {
		return nil
	}
	return &apipb.UpdateProjectProfileResponse{
		Project: NewProject(resp.Project),
	}
}

// NewUploadProjectThumbnailResponse 将 BIZ 上传项目缩略图响应转换为 API 上传项目缩略图响应
func NewUploadProjectThumbnailResponse(resp *bizpb.UploadProjectThumbnailResponse) *apipb.UploadProjectThumbnailResponse {
	if resp == nil {
		return nil
	}
	return &apipb.UploadProjectThumbnailResponse{
		UploadUrl: resp.UploadUrl,
	}
}

// NewDeleteProjectThumbnailResponse 将 BIZ 删除项目缩略图响应转换为 API 删除项目缩略图响应
func NewDeleteProjectThumbnailResponse(_ *bizpb.DeleteProjectThumbnailResponse) *apipb.DeleteProjectThumbnailResponse {
	return &apipb.DeleteProjectThumbnailResponse{}
}
