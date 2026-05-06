package project_core

import (
	"context"

	"icw_common/gen/core/biz"
	"icw_core_api/rpc/common"
	"icw_core_api/rpc/icw_core_biz"
)

// AdvanceProject 项目进度流转
func AdvanceProject(ctx context.Context, client *icw_core_biz.Client, req *bizpb.AdvanceProjectRequest, resp *bizpb.AdvanceProjectResponse) error {
	return common.CallGRPC[bizpb.AdvanceProjectRequest, bizpb.AdvanceProjectResponse](ctx, client, req, resp, client.ProjectCore().AdvanceProject)
}

// CheckProjectAccess 校验项目访问权限
func CheckProjectAccess(ctx context.Context, client *icw_core_biz.Client, req *bizpb.CheckProjectAccessRequest, resp *bizpb.CheckProjectAccessResponse) error {
	return common.CallGRPC[bizpb.CheckProjectAccessRequest, bizpb.CheckProjectAccessResponse](ctx, client, req, resp, client.ProjectCore().CheckProjectAccess)
}

// CreateProject 创建项目
func CreateProject(ctx context.Context, client *icw_core_biz.Client, req *bizpb.CreateProjectRequest, resp *bizpb.CreateProjectResponse) error {
	return common.CallGRPC[bizpb.CreateProjectRequest, bizpb.CreateProjectResponse](ctx, client, req, resp, client.ProjectCore().CreateProject)
}

// DeleteProject 删除项目
func DeleteProject(ctx context.Context, client *icw_core_biz.Client, req *bizpb.DeleteProjectRequest, resp *bizpb.DeleteProjectResponse) error {
	return common.CallGRPC[bizpb.DeleteProjectRequest, bizpb.DeleteProjectResponse](ctx, client, req, resp, client.ProjectCore().DeleteProject)
}

// ListProjects 获取项目列表
func ListProjects(ctx context.Context, client *icw_core_biz.Client, req *bizpb.ListProjectsRequest, resp *bizpb.ListProjectsResponse) error {
	return common.CallGRPC[bizpb.ListProjectsRequest, bizpb.ListProjectsResponse](ctx, client, req, resp, client.ProjectCore().ListProjects)
}
