package dto

import (
	"icw_common/gen/core/api"
	"icw_common/gen/core/biz"
)

// NewGetProjectAssetsResponse 将 BIZ 项目图像列表响应转换为 API 项目图像列表响应
func NewGetProjectAssetsResponse(resp *bizpb.GetProjectAssetsResponse) *apipb.GetProjectAssetsResponse {
	if resp == nil {
		return nil
	}
	return &apipb.GetProjectAssetsResponse{
		Groups: NewProjectGroups(resp.Groups),
	}
}

// NewCreateProjectGroupResponse 将 BIZ 创建图像组响应转换为 API 创建图像组响应
func NewCreateProjectGroupResponse(resp *bizpb.CreateProjectGroupResponse) *apipb.CreateProjectGroupResponse {
	if resp == nil {
		return nil
	}
	return &apipb.CreateProjectGroupResponse{
		Group: NewProjectGroup(resp.Group),
	}
}

// NewDeleteProjectGroupResponse 将 BIZ 删除图像组响应转换为 API 删除图像组响应
func NewDeleteProjectGroupResponse(_ *bizpb.DeleteProjectGroupResponse) *apipb.DeleteProjectGroupResponse {
	return &apipb.DeleteProjectGroupResponse{}
}

// NewMoveProjectGroupResponse 将 BIZ 移动图像组响应转换为 API 移动图像组响应
func NewMoveProjectGroupResponse(resp *bizpb.MoveProjectGroupResponse) *apipb.MoveProjectGroupResponse {
	if resp == nil {
		return nil
	}
	return &apipb.MoveProjectGroupResponse{
		Group: NewProjectGroup(resp.Group),
	}
}

// NewUpdateProjectGroupResponse 将 BIZ 更新图像组响应转换为 API 更新图像组响应
func NewUpdateProjectGroupResponse(resp *bizpb.UpdateProjectGroupResponse) *apipb.UpdateProjectGroupResponse {
	if resp == nil {
		return nil
	}
	return &apipb.UpdateProjectGroupResponse{
		Group: NewProjectGroup(resp.Group),
	}
}

// NewDeleteProjectImageResponse 将 BIZ 删除图像响应转换为 API 删除图像响应
func NewDeleteProjectImageResponse(_ *bizpb.DeleteProjectImageResponse) *apipb.DeleteProjectImageResponse {
	return &apipb.DeleteProjectImageResponse{}
}

// NewGetProjectImageOriginalResponse 将 BIZ 获取原图响应转换为 API 获取原图响应
func NewGetProjectImageOriginalResponse(resp *bizpb.GetProjectImageOriginalResponse) *apipb.GetProjectImageOriginalResponse {
	if resp == nil {
		return nil
	}
	return &apipb.GetProjectImageOriginalResponse{
		OriginalUrl: resp.OriginalUrl,
	}
}

// NewMoveProjectImageResponse 将 BIZ 移动图像响应转换为 API 移动图像响应
func NewMoveProjectImageResponse(resp *bizpb.MoveProjectImageResponse) *apipb.MoveProjectImageResponse {
	if resp == nil {
		return nil
	}
	return &apipb.MoveProjectImageResponse{
		Images: NewProjectImages(resp.Images),
	}
}

// NewReportProjectImageResponse 将 BIZ 上报图像响应转换为 API 上报图像响应
func NewReportProjectImageResponse(_ *bizpb.ReportProjectImageResponse) *apipb.ReportProjectImageResponse {
	return &apipb.ReportProjectImageResponse{}
}

// NewUploadProjectImageItem 将 API 上传图像项转换为 BIZ 上传图像项
func NewUploadProjectImageItem(image *apipb.UploadProjectImageItem) *bizpb.UploadProjectImageItem {
	if image == nil {
		return nil
	}
	return &bizpb.UploadProjectImageItem{
		FileName:    image.FileName,
		ContentType: image.ContentType,
		SizeBytes:   image.SizeBytes,
		Width:       image.Width,
		Height:      image.Height,
		Metadata:    image.Metadata,
	}
}

// NewUploadProjectImageResults 将 BIZ 上传图像结果切片转换为 API 上传图像结果切片
func NewUploadProjectImageResults(images []*bizpb.UploadProjectImageResult) []*apipb.UploadProjectImageResult {
	return images
}

// NewUploadProjectImageResponse 将 BIZ 上传图像响应转换为 API 上传图像响应
func NewUploadProjectImageResponse(resp *bizpb.UploadProjectImageResponse) *apipb.UploadProjectImageResponse {
	if resp == nil {
		return nil
	}
	return &apipb.UploadProjectImageResponse{
		Images: NewUploadProjectImageResults(resp.Images),
	}
}
