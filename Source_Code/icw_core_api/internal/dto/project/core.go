package project

import (
	"icw_core_biz/pkg/dto/project"
	"icw_core_biz/utils"
)

type Project struct {
	Id                  string `json:"id"`
	Name                string `json:"name"`
	BuildingName        string `json:"building_name"`
	BuildingLocation    string `json:"building_location"`
	BuiltYear           uint16 `json:"built_year"`
	BuildingDescription string `json:"building_description"`
	KnownIssues         string `json:"known_issues"`
	AssessmentGoal      string `json:"assessment_goal"`
	ThumbnailURL        string `json:"thumbnail_url"`
	Progress            uint8  `json:"progress"`
	CreatedAt           string `json:"created_at"`
	UpdatedAt           string `json:"updated_at"`
}

func NewProject(project *project.Project) *Project {
	if project == nil {
		return nil
	}
	return &Project{
		Id:                  utils.Encode(project.Id),
		Name:                project.Name,
		BuildingName:        project.BuildingName,
		BuildingLocation:    project.BuildingLocation,
		BuiltYear:           project.BuiltYear,
		BuildingDescription: project.BuildingDescription,
		KnownIssues:         project.KnownIssues,
		AssessmentGoal:      project.AssessmentGoal,
		ThumbnailURL:        project.ThumbnailURL,
		Progress:            project.Progress,
		CreatedAt:           project.CreatedAt,
		UpdatedAt:           project.UpdatedAt,
	}
}

type ProjectListItem struct {
	Id               string `json:"id"`
	Name             string `json:"name"`
	BuildingName     string `json:"building_name"`
	BuildingLocation string `json:"building_location"`
	ThumbnailURL     string `json:"thumbnail_url"`
	Progress         uint8  `json:"progress"`
	CreatedAt        string `json:"created_at"`
}

func NewProjectListItem(project *project.ProjectListItem) *ProjectListItem {
	if project == nil {
		return nil
	}
	return &ProjectListItem{
		Id:               utils.Encode(project.Id),
		Name:             project.Name,
		BuildingName:     project.BuildingName,
		BuildingLocation: project.BuildingLocation,
		ThumbnailURL:     project.ThumbnailURL,
		Progress:         project.Progress,
		CreatedAt:        project.CreatedAt,
	}
}

func NewProjectListItems(projects []*project.ProjectListItem) []*ProjectListItem {
	if projects == nil {
		return make([]*ProjectListItem, 0)
	}
	items := make([]*ProjectListItem, 0, len(projects))
	for _, item := range projects {
		if item == nil {
			continue
		}
		items = append(items, NewProjectListItem(item))
	}
	return items
}

type AdvanceProjectRequest struct {
	ProjectId    string `json:"project_id"`
	FromProgress uint8  `json:"from_progress"`
	ToProgress   uint8  `json:"to_progress"`
}

type AdvanceProjectResponse struct{}

func NewAdvanceProjectResponse(_ *project.AdvanceProjectResponse) *AdvanceProjectResponse {
	return &AdvanceProjectResponse{}
}

type CreateProjectResponse struct {
	Project *Project `json:"project"`
}

func NewCreateProjectResponse(resp *project.CreateProjectResponse) *CreateProjectResponse {
	if resp == nil {
		return nil
	}
	return &CreateProjectResponse{
		Project: NewProject(resp.Project),
	}
}

type DeleteProjectRequest struct {
	ProjectId string `json:"project_id"`
}

type DeleteProjectResponse struct {
	ActiveProjects    []*ProjectListItem `json:"active_projects"`
	CompletedProjects []*ProjectListItem `json:"completed_projects"`
}

func NewDeleteProjectResponse(resp *project.DeleteProjectResponse) *DeleteProjectResponse {
	if resp == nil {
		return nil
	}
	return &DeleteProjectResponse{
		ActiveProjects:    NewProjectListItems(resp.ActiveProjects),
		CompletedProjects: NewProjectListItems(resp.CompletedProjects),
	}
}

type ListProjectsResponse struct {
	ActiveProjects    []*ProjectListItem `json:"active_projects"`
	CompletedProjects []*ProjectListItem `json:"completed_projects"`
}

func NewListProjectsResponse(resp *project.ListProjectsResponse) *ListProjectsResponse {
	if resp == nil {
		return nil
	}
	return &ListProjectsResponse{
		ActiveProjects:    NewProjectListItems(resp.ActiveProjects),
		CompletedProjects: NewProjectListItems(resp.CompletedProjects),
	}
}
