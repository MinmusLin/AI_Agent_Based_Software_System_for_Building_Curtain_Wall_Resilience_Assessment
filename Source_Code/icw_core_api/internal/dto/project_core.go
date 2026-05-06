package dto

import (
	"icw_common/gen/core/api"
	"icw_common/gen/core/biz"
)

func NewAdvanceProjectResponse(_ *bizpb.AdvanceProjectResponse) *apipb.AdvanceProjectResponse {
	return &apipb.AdvanceProjectResponse{}
}

func NewCreateProjectResponse(resp *bizpb.CreateProjectResponse) *apipb.CreateProjectResponse {
	if resp == nil {
		return nil
	}
	return &apipb.CreateProjectResponse{
		Project: NewProject(resp.Project),
	}
}

func NewDeleteProjectResponse(resp *bizpb.DeleteProjectResponse) *apipb.DeleteProjectResponse {
	if resp == nil {
		return nil
	}
	return &apipb.DeleteProjectResponse{
		ActiveProjects:    NewProjectListItems(resp.ActiveProjects),
		CompletedProjects: NewProjectListItems(resp.CompletedProjects),
	}
}

func NewListProjectsResponse(resp *bizpb.ListProjectsResponse) *apipb.ListProjectsResponse {
	if resp == nil {
		return nil
	}
	return &apipb.ListProjectsResponse{
		ActiveProjects:    NewProjectListItems(resp.ActiveProjects),
		CompletedProjects: NewProjectListItems(resp.CompletedProjects),
	}
}
