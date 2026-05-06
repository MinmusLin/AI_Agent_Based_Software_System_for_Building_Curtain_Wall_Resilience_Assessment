package project_profile

import (
	"context"

	"icw_common/gen/core/biz"
	"icw_core_api/rpc/common"
	"icw_core_api/rpc/icw_core_biz"
)

// DeleteProjectThumbnail 删除项目缩略图
func DeleteProjectThumbnail(ctx context.Context, client *icw_core_biz.Client, req *bizpb.DeleteProjectThumbnailRequest, resp *bizpb.DeleteProjectThumbnailResponse) error {
	return common.CallGRPC[bizpb.DeleteProjectThumbnailRequest, bizpb.DeleteProjectThumbnailResponse](ctx, client, req, resp, client.ProjectProfile().DeleteProjectThumbnail)
}

// GetProjectProfile 获取项目基础信息
func GetProjectProfile(ctx context.Context, client *icw_core_biz.Client, req *bizpb.GetProjectProfileRequest, resp *bizpb.GetProjectProfileResponse) error {
	return common.CallGRPC[bizpb.GetProjectProfileRequest, bizpb.GetProjectProfileResponse](ctx, client, req, resp, client.ProjectProfile().GetProjectProfile)
}

// GetProjectThumbnail 获取项目缩略图
func GetProjectThumbnail(ctx context.Context, client *icw_core_biz.Client, req *bizpb.GetProjectThumbnailRequest, resp *bizpb.GetProjectThumbnailResponse) error {
	return common.CallGRPC[bizpb.GetProjectThumbnailRequest, bizpb.GetProjectThumbnailResponse](ctx, client, req, resp, client.ProjectProfile().GetProjectThumbnail)
}

// UpdateProjectProfile 更新项目基础信息
func UpdateProjectProfile(ctx context.Context, client *icw_core_biz.Client, req *bizpb.UpdateProjectProfileRequest, resp *bizpb.UpdateProjectProfileResponse) error {
	return common.CallGRPC[bizpb.UpdateProjectProfileRequest, bizpb.UpdateProjectProfileResponse](ctx, client, req, resp, client.ProjectProfile().UpdateProjectProfile)
}

// UploadProjectThumbnail 上传项目缩略图
func UploadProjectThumbnail(ctx context.Context, client *icw_core_biz.Client, req *bizpb.UploadProjectThumbnailRequest, resp *bizpb.UploadProjectThumbnailResponse) error {
	return common.CallGRPC[bizpb.UploadProjectThumbnailRequest, bizpb.UploadProjectThumbnailResponse](ctx, client, req, resp, client.ProjectProfile().UploadProjectThumbnail)
}
