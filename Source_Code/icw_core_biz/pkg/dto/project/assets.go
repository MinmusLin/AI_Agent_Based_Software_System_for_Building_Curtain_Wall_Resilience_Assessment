package project

import (
	"icw_core_biz/pkg/dto"
)

type ProjectGroup struct {
	Id        uint64
	Name      string
	SortOrder string
	Images    []*ProjectImage
}

type ProjectImage struct {
	Uuid         string `json:"uuid"`
	FileName     string `json:"file_name"`
	ContentType  string `json:"content_type"`
	SizeBytes    uint64 `json:"size_bytes"`
	Width        uint32 `json:"width"`
	Height       uint32 `json:"height"`
	Metadata     string `json:"metadata"`
	Status       string `json:"status"`
	ThumbnailURL string `json:"thumbnail_url"`
	UploadedAt   string `json:"uploaded_at"`
	CreatedAt    string `json:"created_at"`
}

type UploadProjectImageItem struct {
	FileName    string
	ContentType string
	SizeBytes   uint64
	Width       uint32
	Height      uint32
	Metadata    string
}

type UploadProjectImageResult struct {
	Image              *ProjectImage
	OriginalUploadURL  string
	ThumbnailUploadURL string
}

type GetProjectAssetsRequest struct {
	Meta      *dto.Meta
	UserId    uint64
	ProjectId uint64
}

type GetProjectAssetsResponse struct {
	Groups []*ProjectGroup
}

type CreateProjectGroupRequest struct {
	Meta      *dto.Meta
	UserId    uint64
	ProjectId uint64
}

type CreateProjectGroupResponse struct {
	Group *ProjectGroup
}

type DeleteProjectGroupRequest struct {
	Meta      *dto.Meta
	UserId    uint64
	ProjectId uint64
	GroupId   uint64
}

type DeleteProjectGroupResponse struct{}

type MoveProjectGroupRequest struct {
	Meta            *dto.Meta
	UserId          uint64
	ProjectId       uint64
	GroupId         uint64
	PreviousGroupId uint64
	NextGroupId     uint64
	MoveToFirst     bool
	MoveToLast      bool
}

type MoveProjectGroupResponse struct {
	Group *ProjectGroup
}

type UpdateProjectGroupRequest struct {
	Meta      *dto.Meta
	UserId    uint64
	ProjectId uint64
	GroupId   uint64
	Name      string
}

type UpdateProjectGroupResponse struct {
	Group *ProjectGroup
}

type DeleteProjectImageRequest struct {
	Meta       *dto.Meta
	UserId     uint64
	ProjectId  uint64
	ImageUuids []string
}

type DeleteProjectImageResponse struct{}

type GetProjectImageOriginalRequest struct {
	Meta      *dto.Meta
	UserId    uint64
	ProjectId uint64
	ImageUuid string
}

type GetProjectImageOriginalResponse struct {
	OriginalURL string
}

type MoveProjectImageRequest struct {
	Meta          *dto.Meta
	UserId        uint64
	ProjectId     uint64
	ImageUuids    []string
	TargetGroupId uint64
}

type MoveProjectImageResponse struct {
	Images []*ProjectImage
}

type ReportProjectImageRequest struct {
	Meta      *dto.Meta
	UserId    uint64
	ProjectId uint64
	ImageUuid string
	Status    string
}

type ReportProjectImageResponse struct {
	Image *ProjectImage
}

type UploadProjectImageRequest struct {
	Meta      *dto.Meta
	UserId    uint64
	ProjectId uint64
	GroupId   uint64
	Images    []*UploadProjectImageItem
}

type UploadProjectImageResponse struct {
	Images []*UploadProjectImageResult
}
