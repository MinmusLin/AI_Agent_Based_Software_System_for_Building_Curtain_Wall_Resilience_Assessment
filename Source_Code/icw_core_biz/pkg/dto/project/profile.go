package project

import (
	"icw_core_biz/pkg/dto"
)

type GetProjectProfileRequest struct {
	Meta      *dto.Meta
	UserId    uint64
	ProjectId uint64
}

type GetProjectProfileResponse struct {
	Project *Project
}

type GetProjectThumbnailRequest struct {
	Meta      *dto.Meta
	UserId    uint64
	ProjectId uint64
}

type GetProjectThumbnailResponse struct {
	ThumbnailURL string
}

type UploadProjectThumbnailRequest struct {
	Meta        *dto.Meta
	UserId      uint64
	ProjectId   uint64
	ContentType string
}

type UploadProjectThumbnailResponse struct {
	UploadURL string
}

type DeleteProjectThumbnailRequest struct {
	Meta      *dto.Meta
	UserId    uint64
	ProjectId uint64
}

type DeleteProjectThumbnailResponse struct{}

type UpdateProjectProfileRequest struct {
	Meta                *dto.Meta
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
