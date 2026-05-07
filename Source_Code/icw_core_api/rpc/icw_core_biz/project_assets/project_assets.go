package project_assets

import (
	"context"

	"icw_common/gen/core/biz"
	"icw_common/rpc"

	"icw_core_api/rpc/icw_core_biz"
	"icw_core_api/utils"
)

// CreateProjectGroup 创建图像组
func CreateProjectGroup(ctx context.Context, client *icw_core_biz.Client, req *bizpb.CreateProjectGroupRequest, resp *bizpb.CreateProjectGroupResponse) error {
	return rpc.CallGRPC[bizpb.CreateProjectGroupRequest, bizpb.CreateProjectGroupResponse](ctx, client, req, resp, client.ProjectAssets().CreateProjectGroup, rpc.WithRequestIdResolver(utils.GetXRequestId))
}

// DeleteProjectGroup 删除图像组
func DeleteProjectGroup(ctx context.Context, client *icw_core_biz.Client, req *bizpb.DeleteProjectGroupRequest, resp *bizpb.DeleteProjectGroupResponse) error {
	return rpc.CallGRPC[bizpb.DeleteProjectGroupRequest, bizpb.DeleteProjectGroupResponse](ctx, client, req, resp, client.ProjectAssets().DeleteProjectGroup, rpc.WithRequestIdResolver(utils.GetXRequestId))
}

// DeleteProjectImage 删除图像
func DeleteProjectImage(ctx context.Context, client *icw_core_biz.Client, req *bizpb.DeleteProjectImageRequest, resp *bizpb.DeleteProjectImageResponse) error {
	return rpc.CallGRPC[bizpb.DeleteProjectImageRequest, bizpb.DeleteProjectImageResponse](ctx, client, req, resp, client.ProjectAssets().DeleteProjectImage, rpc.WithRequestIdResolver(utils.GetXRequestId))
}

// GetProjectAssets 获取项目图像列表
func GetProjectAssets(ctx context.Context, client *icw_core_biz.Client, req *bizpb.GetProjectAssetsRequest, resp *bizpb.GetProjectAssetsResponse) error {
	return rpc.CallGRPC[bizpb.GetProjectAssetsRequest, bizpb.GetProjectAssetsResponse](ctx, client, req, resp, client.ProjectAssets().GetProjectAssets, rpc.WithRequestIdResolver(utils.GetXRequestId))
}

// GetProjectImageOriginal 获取原图
func GetProjectImageOriginal(ctx context.Context, client *icw_core_biz.Client, req *bizpb.GetProjectImageOriginalRequest, resp *bizpb.GetProjectImageOriginalResponse) error {
	return rpc.CallGRPC[bizpb.GetProjectImageOriginalRequest, bizpb.GetProjectImageOriginalResponse](ctx, client, req, resp, client.ProjectAssets().GetProjectImageOriginal, rpc.WithRequestIdResolver(utils.GetXRequestId))
}

// MoveProjectGroup 移动图像组
func MoveProjectGroup(ctx context.Context, client *icw_core_biz.Client, req *bizpb.MoveProjectGroupRequest, resp *bizpb.MoveProjectGroupResponse) error {
	return rpc.CallGRPC[bizpb.MoveProjectGroupRequest, bizpb.MoveProjectGroupResponse](ctx, client, req, resp, client.ProjectAssets().MoveProjectGroup, rpc.WithRequestIdResolver(utils.GetXRequestId))
}

// MoveProjectImage 移动图像
func MoveProjectImage(ctx context.Context, client *icw_core_biz.Client, req *bizpb.MoveProjectImageRequest, resp *bizpb.MoveProjectImageResponse) error {
	return rpc.CallGRPC[bizpb.MoveProjectImageRequest, bizpb.MoveProjectImageResponse](ctx, client, req, resp, client.ProjectAssets().MoveProjectImage, rpc.WithRequestIdResolver(utils.GetXRequestId))
}

// ReportProjectImage 上报图像
func ReportProjectImage(ctx context.Context, client *icw_core_biz.Client, req *bizpb.ReportProjectImageRequest, resp *bizpb.ReportProjectImageResponse) error {
	return rpc.CallGRPC[bizpb.ReportProjectImageRequest, bizpb.ReportProjectImageResponse](ctx, client, req, resp, client.ProjectAssets().ReportProjectImage, rpc.WithRequestIdResolver(utils.GetXRequestId))
}

// UpdateProjectGroup 更新图像组
func UpdateProjectGroup(ctx context.Context, client *icw_core_biz.Client, req *bizpb.UpdateProjectGroupRequest, resp *bizpb.UpdateProjectGroupResponse) error {
	return rpc.CallGRPC[bizpb.UpdateProjectGroupRequest, bizpb.UpdateProjectGroupResponse](ctx, client, req, resp, client.ProjectAssets().UpdateProjectGroup, rpc.WithRequestIdResolver(utils.GetXRequestId))
}

// UploadProjectImage 上传图像
func UploadProjectImage(ctx context.Context, client *icw_core_biz.Client, req *bizpb.UploadProjectImageRequest, resp *bizpb.UploadProjectImageResponse) error {
	return rpc.CallGRPC[bizpb.UploadProjectImageRequest, bizpb.UploadProjectImageResponse](ctx, client, req, resp, client.ProjectAssets().UploadProjectImage, rpc.WithRequestIdResolver(utils.GetXRequestId))
}
