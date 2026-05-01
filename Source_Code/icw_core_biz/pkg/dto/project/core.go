package project

type Project struct {
	Id                  uint64
	Name                string
	BuildingName        string
	BuildingLocation    string
	BuiltYear           uint16
	BuildingDescription string
	KnownIssues         string
	AssessmentGoal      string
	Progress            uint8
	CreatedAt           string
	UpdatedAt           string
}

type ProjectListItem struct {
	Id               uint64
	Name             string
	BuildingName     string
	BuildingLocation string
	Progress         uint8
	CreatedAt        string
}

type AdvanceProjectRequest struct {
	UserId       uint64
	ProjectId    uint64
	FromProgress uint8
	ToProgress   uint8
}

type AdvanceProjectResponse struct{}

type CreateProjectRequest struct {
	UserId uint64
}

type CreateProjectResponse struct {
	Project *Project
}

type DeleteProjectRequest struct {
	UserId    uint64
	ProjectId uint64
}

type DeleteProjectResponse struct {
	ActiveProjects    []*ProjectListItem
	CompletedProjects []*ProjectListItem
}

type ListProjectsRequest struct {
	UserId uint64
}

type ListProjectsResponse struct {
	ActiveProjects    []*ProjectListItem
	CompletedProjects []*ProjectListItem
}
