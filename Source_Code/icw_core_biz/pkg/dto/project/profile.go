package project

type GetProjectProfileRequest struct {
	UserId    uint64
	ProjectId uint64
}

type GetProjectProfileResponse struct {
	Project *Project
}

type UpdateProjectProfileRequest struct {
	UserId              uint64
	ProjectId           uint64
	Name                string
	BuildingName        string
	BuildingLocation    string
	BuiltYear           uint16
	BuildingDescription string
	KnownIssues         string
	AssessmentGoal      string
}

type UpdateProjectProfileResponse struct {
	Project *Project
}
