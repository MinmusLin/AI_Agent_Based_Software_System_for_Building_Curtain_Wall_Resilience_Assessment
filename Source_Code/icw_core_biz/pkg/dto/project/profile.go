package project

type GetProjectProfileRequest struct {
	UserId    uint64
	ProjectId uint64
}

type GetProjectProfileResponse struct {
	Project *Project
}

type GetProjectThumbnailRequest struct {
	UserId    uint64
	ProjectId uint64
}

type GetProjectThumbnailResponse struct {
	ThumbnailURL string
}

type UploadProjectThumbnailRequest struct {
	UserId      uint64
	ProjectId   uint64
	ContentType string
}

type UploadProjectThumbnailResponse struct {
	UploadURL string
}

type DeleteProjectThumbnailRequest struct {
	UserId    uint64
	ProjectId uint64
}

type DeleteProjectThumbnailResponse struct{}

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
