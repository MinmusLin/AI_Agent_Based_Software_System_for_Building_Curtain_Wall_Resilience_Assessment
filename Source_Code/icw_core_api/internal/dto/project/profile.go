package project

import (
	"icw_core_biz/pkg/dto/project"
)

type GetProjectProfileResponse struct {
	Project *Project `json:"project"`
}

func NewGetProjectProfileResponse(resp *project.GetProjectProfileResponse) *GetProjectProfileResponse {
	if resp == nil {
		return nil
	}
	return &GetProjectProfileResponse{
		Project: NewProject(resp.Project),
	}
}

type UpdateProjectProfileRequest struct {
	ProjectId           string `json:"project_id"`
	Name                string `json:"name"`
	BuildingName        string `json:"building_name"`
	BuildingLocation    string `json:"building_location"`
	BuiltYear           uint16 `json:"built_year"`
	BuildingDescription string `json:"building_description"`
	KnownIssues         string `json:"known_issues"`
	AssessmentGoal      string `json:"assessment_goal"`
}

type UpdateProjectProfileResponse struct {
	Project *Project `json:"project"`
}

func NewUpdateProjectProfileResponse(resp *project.UpdateProjectProfileResponse) *UpdateProjectProfileResponse {
	if resp == nil {
		return nil
	}
	return &UpdateProjectProfileResponse{
		Project: NewProject(resp.Project),
	}
}
