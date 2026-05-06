package dto

import (
	"icw_common/gen/core/api"
	"icw_common/gen/core/biz"
)

// NewAdvanceProjectResponse 将 BIZ 项目进度流转响应转换为 API 项目进度流转响应
func NewAdvanceProjectResponse(_ *bizpb.AdvanceProjectResponse) *apipb.AdvanceProjectResponse {
	return &apipb.AdvanceProjectResponse{}
}

// NewCreateProjectResponse 将 BIZ 创建项目响应转换为 API 创建项目响应
func NewCreateProjectResponse(resp *bizpb.CreateProjectResponse) *apipb.CreateProjectResponse {
	if resp == nil {
		return nil
	}
	return &apipb.CreateProjectResponse{
		Project: NewProject(resp.Project),
	}
}

// NewDeleteProjectResponse 将 BIZ 删除项目响应转换为 API 删除项目响应
func NewDeleteProjectResponse(resp *bizpb.DeleteProjectResponse) *apipb.DeleteProjectResponse {
	if resp == nil {
		return nil
	}
	return &apipb.DeleteProjectResponse{
		ActiveProjects:    NewProjectListItems(resp.ActiveProjects),
		CompletedProjects: NewProjectListItems(resp.CompletedProjects),
	}
}

// NewListProjectsResponse 将 BIZ 项目列表响应转换为 API 项目列表响应
func NewListProjectsResponse(resp *bizpb.ListProjectsResponse) *apipb.ListProjectsResponse {
	if resp == nil {
		return nil
	}
	return &apipb.ListProjectsResponse{
		ActiveProjects:    NewProjectListItems(resp.ActiveProjects),
		CompletedProjects: NewProjectListItems(resp.CompletedProjects),
	}
}
