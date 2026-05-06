package project_core

import (
	"context"

	"icw_common/gen/core/biz"
	"icw_common/rpc"
	"icw_core_api/rpc/icw_core_biz"
	"icw_core_api/utils"
)

// AdvanceProject 项目进度流转
func AdvanceProject(ctx context.Context, client *icw_core_biz.Client, req *bizpb.AdvanceProjectRequest, resp *bizpb.AdvanceProjectResponse) error {
	return rpc.CallGRPC[bizpb.AdvanceProjectRequest, bizpb.AdvanceProjectResponse](ctx, client, req, resp, client.ProjectCore().AdvanceProject, rpc.WithRequestIdResolver(utils.GetXRequestId))
}

// CheckProjectAccess 校验项目访问权限
func CheckProjectAccess(ctx context.Context, client *icw_core_biz.Client, req *bizpb.CheckProjectAccessRequest, resp *bizpb.CheckProjectAccessResponse) error {
	return rpc.CallGRPC[bizpb.CheckProjectAccessRequest, bizpb.CheckProjectAccessResponse](ctx, client, req, resp, client.ProjectCore().CheckProjectAccess, rpc.WithRequestIdResolver(utils.GetXRequestId))
}

// CreateProject 创建项目
func CreateProject(ctx context.Context, client *icw_core_biz.Client, req *bizpb.CreateProjectRequest, resp *bizpb.CreateProjectResponse) error {
	return rpc.CallGRPC[bizpb.CreateProjectRequest, bizpb.CreateProjectResponse](ctx, client, req, resp, client.ProjectCore().CreateProject, rpc.WithRequestIdResolver(utils.GetXRequestId))
}

// DeleteProject 删除项目
func DeleteProject(ctx context.Context, client *icw_core_biz.Client, req *bizpb.DeleteProjectRequest, resp *bizpb.DeleteProjectResponse) error {
	return rpc.CallGRPC[bizpb.DeleteProjectRequest, bizpb.DeleteProjectResponse](ctx, client, req, resp, client.ProjectCore().DeleteProject, rpc.WithRequestIdResolver(utils.GetXRequestId))
}

// ListProjects 获取项目列表
func ListProjects(ctx context.Context, client *icw_core_biz.Client, req *bizpb.ListProjectsRequest, resp *bizpb.ListProjectsResponse) error {
	return rpc.CallGRPC[bizpb.ListProjectsRequest, bizpb.ListProjectsResponse](ctx, client, req, resp, client.ProjectCore().ListProjects, rpc.WithRequestIdResolver(utils.GetXRequestId))
}
