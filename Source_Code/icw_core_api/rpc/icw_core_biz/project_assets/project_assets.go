package project_assets

import (
	"context"

	"icw_common/gen/core/biz"
	"icw_core_api/rpc/common"
	"icw_core_api/rpc/icw_core_biz"
)

// CreateProjectGroup 创建图像组
func CreateProjectGroup(ctx context.Context, client *icw_core_biz.Client, req *bizpb.CreateProjectGroupRequest, resp *bizpb.CreateProjectGroupResponse) error {
	return common.CallGRPC[bizpb.CreateProjectGroupRequest, bizpb.CreateProjectGroupResponse](ctx, client, req, resp, client.ProjectAssets().CreateProjectGroup)
}

// DeleteProjectGroup 删除图像组
func DeleteProjectGroup(ctx context.Context, client *icw_core_biz.Client, req *bizpb.DeleteProjectGroupRequest, resp *bizpb.DeleteProjectGroupResponse) error {
	return common.CallGRPC[bizpb.DeleteProjectGroupRequest, bizpb.DeleteProjectGroupResponse](ctx, client, req, resp, client.ProjectAssets().DeleteProjectGroup)
}

// DeleteProjectImage 删除图像
func DeleteProjectImage(ctx context.Context, client *icw_core_biz.Client, req *bizpb.DeleteProjectImageRequest, resp *bizpb.DeleteProjectImageResponse) error {
	return common.CallGRPC[bizpb.DeleteProjectImageRequest, bizpb.DeleteProjectImageResponse](ctx, client, req, resp, client.ProjectAssets().DeleteProjectImage)
}

// GetProjectAssets 获取项目图像列表
func GetProjectAssets(ctx context.Context, client *icw_core_biz.Client, req *bizpb.GetProjectAssetsRequest, resp *bizpb.GetProjectAssetsResponse) error {
	return common.CallGRPC[bizpb.GetProjectAssetsRequest, bizpb.GetProjectAssetsResponse](ctx, client, req, resp, client.ProjectAssets().GetProjectAssets)
}

// GetProjectImageOriginal 获取原图
func GetProjectImageOriginal(ctx context.Context, client *icw_core_biz.Client, req *bizpb.GetProjectImageOriginalRequest, resp *bizpb.GetProjectImageOriginalResponse) error {
	return common.CallGRPC[bizpb.GetProjectImageOriginalRequest, bizpb.GetProjectImageOriginalResponse](ctx, client, req, resp, client.ProjectAssets().GetProjectImageOriginal)
}

// MoveProjectGroup 移动图像组
func MoveProjectGroup(ctx context.Context, client *icw_core_biz.Client, req *bizpb.MoveProjectGroupRequest, resp *bizpb.MoveProjectGroupResponse) error {
	return common.CallGRPC[bizpb.MoveProjectGroupRequest, bizpb.MoveProjectGroupResponse](ctx, client, req, resp, client.ProjectAssets().MoveProjectGroup)
}

// MoveProjectImage 移动图像
func MoveProjectImage(ctx context.Context, client *icw_core_biz.Client, req *bizpb.MoveProjectImageRequest, resp *bizpb.MoveProjectImageResponse) error {
	return common.CallGRPC[bizpb.MoveProjectImageRequest, bizpb.MoveProjectImageResponse](ctx, client, req, resp, client.ProjectAssets().MoveProjectImage)
}

// ReportProjectImage 上报图像
func ReportProjectImage(ctx context.Context, client *icw_core_biz.Client, req *bizpb.ReportProjectImageRequest, resp *bizpb.ReportProjectImageResponse) error {
	return common.CallGRPC[bizpb.ReportProjectImageRequest, bizpb.ReportProjectImageResponse](ctx, client, req, resp, client.ProjectAssets().ReportProjectImage)
}

// UpdateProjectGroup 更新图像组
func UpdateProjectGroup(ctx context.Context, client *icw_core_biz.Client, req *bizpb.UpdateProjectGroupRequest, resp *bizpb.UpdateProjectGroupResponse) error {
	return common.CallGRPC[bizpb.UpdateProjectGroupRequest, bizpb.UpdateProjectGroupResponse](ctx, client, req, resp, client.ProjectAssets().UpdateProjectGroup)
}

// UploadProjectImage 上传图像
func UploadProjectImage(ctx context.Context, client *icw_core_biz.Client, req *bizpb.UploadProjectImageRequest, resp *bizpb.UploadProjectImageResponse) error {
	return common.CallGRPC[bizpb.UploadProjectImageRequest, bizpb.UploadProjectImageResponse](ctx, client, req, resp, client.ProjectAssets().UploadProjectImage)
}
