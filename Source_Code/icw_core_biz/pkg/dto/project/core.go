package project

import (
	"icw_core_biz/pkg/dto"
)

type Project struct {
	Id                  uint64
	Name                string
	BuildingName        string
	BuildingLocation    string
	BuiltYear           uint16
	BuildingDescription string
	KnownIssues         string
	AssessmentGoal      string
	ThumbnailURL        string
	Progress            uint8
	CreatedAt           string
	UpdatedAt           string
}

type ProjectListItem struct {
	Id               uint64
	Name             string
	BuildingName     string
	BuildingLocation string
	ThumbnailURL     string
	Progress         uint8
	CreatedAt        string
}

type CheckProjectAccessRequest struct {
	Meta      *dto.Meta
	UserId    uint64
	ProjectId uint64
}

type CheckProjectAccessResponse struct {
	ProjectId uint64
	Progress  uint8
	Status    string
}

type AdvanceProjectRequest struct {
	Meta         *dto.Meta
	UserId       uint64
	ProjectId    uint64
	FromProgress uint8
	ToProgress   uint8
}

type AdvanceProjectResponse struct{}

type CreateProjectRequest struct {
	Meta   *dto.Meta
	UserId uint64
}

type CreateProjectResponse struct {
	Project *Project
}

type DeleteProjectRequest struct {
	Meta      *dto.Meta
	UserId    uint64
	ProjectId uint64
}

type DeleteProjectResponse struct {
	ActiveProjects    []*ProjectListItem
	CompletedProjects []*ProjectListItem
}

type ListProjectsRequest struct {
	Meta   *dto.Meta
	UserId uint64
}

type ListProjectsResponse struct {
	ActiveProjects    []*ProjectListItem
	CompletedProjects []*ProjectListItem
}
