package project

import (
	"strconv"

	"icw_core_biz/pkg/dto/project"
)

type ProjectGroup struct {
	Id        string          `json:"id"`
	Name      string          `json:"name"`
	SortOrder string          `json:"sort_order"`
	Images    []*ProjectImage `json:"images"`
}

func NewProjectGroup(group *project.ProjectGroup) *ProjectGroup {
	if group == nil {
		return nil
	}
	return &ProjectGroup{
		Id:        strconv.FormatUint(group.Id, 10),
		Name:      group.Name,
		SortOrder: group.SortOrder,
		Images:    NewProjectImages(group.Images),
	}
}

func NewProjectGroups(groups []*project.ProjectGroup) []*ProjectGroup {
	if groups == nil {
		return make([]*ProjectGroup, 0)
	}
	items := make([]*ProjectGroup, 0, len(groups))
	for _, group := range groups {
		if group == nil {
			continue
		}
		items = append(items, NewProjectGroup(group))
	}
	return items
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

func NewProjectImage(image *project.ProjectImage) *ProjectImage {
	if image == nil {
		return nil
	}
	return &ProjectImage{
		Uuid:         image.Uuid,
		FileName:     image.FileName,
		ContentType:  image.ContentType,
		SizeBytes:    image.SizeBytes,
		Width:        image.Width,
		Height:       image.Height,
		Metadata:     image.Metadata,
		Status:       image.Status,
		ThumbnailURL: image.ThumbnailURL,
		UploadedAt:   image.UploadedAt,
		CreatedAt:    image.CreatedAt,
	}
}

func NewProjectImages(images []*project.ProjectImage) []*ProjectImage {
	if images == nil {
		return make([]*ProjectImage, 0)
	}
	items := make([]*ProjectImage, 0, len(images))
	for _, image := range images {
		if image == nil {
			continue
		}
		items = append(items, NewProjectImage(image))
	}
	return items
}

type UploadProjectImageItem struct {
	FileName    string `json:"file_name"`
	ContentType string `json:"content_type"`
	SizeBytes   uint64 `json:"size_bytes"`
	Width       uint32 `json:"width"`
	Height      uint32 `json:"height"`
	Metadata    string `json:"metadata"`
}

type GetProjectAssetsResponse struct {
	Groups []*ProjectGroup `json:"groups"`
}

func NewGetProjectAssetsResponse(resp *project.GetProjectAssetsResponse) *GetProjectAssetsResponse {
	if resp == nil {
		return nil
	}
	return &GetProjectAssetsResponse{
		Groups: NewProjectGroups(resp.Groups),
	}
}

type CreateProjectGroupRequest struct {
	ProjectId string `json:"project_id"`
}

type CreateProjectGroupResponse struct {
	Group *ProjectGroup `json:"group"`
}

func NewCreateProjectGroupResponse(resp *project.CreateProjectGroupResponse) *CreateProjectGroupResponse {
	if resp == nil {
		return nil
	}
	return &CreateProjectGroupResponse{
		Group: NewProjectGroup(resp.Group),
	}
}

type DeleteProjectGroupRequest struct {
	ProjectId string `json:"project_id"`
	GroupId   string `json:"group_id"`
}

type DeleteProjectGroupResponse struct{}

func NewDeleteProjectGroupResponse(_ *project.DeleteProjectGroupResponse) *DeleteProjectGroupResponse {
	return &DeleteProjectGroupResponse{}
}

type MoveProjectGroupRequest struct {
	ProjectId       string `json:"project_id"`
	GroupId         string `json:"group_id"`
	PreviousGroupId string `json:"previous_group_id"`
	NextGroupId     string `json:"next_group_id"`
	MoveToFirst     bool   `json:"move_to_first"`
	MoveToLast      bool   `json:"move_to_last"`
}

type MoveProjectGroupResponse struct {
	Group *ProjectGroup `json:"group"`
}

func NewMoveProjectGroupResponse(resp *project.MoveProjectGroupResponse) *MoveProjectGroupResponse {
	if resp == nil {
		return nil
	}
	return &MoveProjectGroupResponse{
		Group: NewProjectGroup(resp.Group),
	}
}

type UpdateProjectGroupRequest struct {
	ProjectId string `json:"project_id"`
	GroupId   string `json:"group_id"`
	Name      string `json:"name"`
}

type UpdateProjectGroupResponse struct {
	Group *ProjectGroup `json:"group"`
}

func NewUpdateProjectGroupResponse(resp *project.UpdateProjectGroupResponse) *UpdateProjectGroupResponse {
	if resp == nil {
		return nil
	}
	return &UpdateProjectGroupResponse{
		Group: NewProjectGroup(resp.Group),
	}
}

type DeleteProjectImageRequest struct {
	ProjectId  string   `json:"project_id"`
	ImageUuids []string `json:"image_uuids"`
}

type DeleteProjectImageResponse struct{}

func NewDeleteProjectImageResponse(_ *project.DeleteProjectImageResponse) *DeleteProjectImageResponse {
	return &DeleteProjectImageResponse{}
}

type GetProjectImageOriginalResponse struct {
	OriginalURL string `json:"original_url"`
}

func NewGetProjectImageOriginalResponse(resp *project.GetProjectImageOriginalResponse) *GetProjectImageOriginalResponse {
	if resp == nil {
		return nil
	}
	return &GetProjectImageOriginalResponse{
		OriginalURL: resp.OriginalURL,
	}
}

type MoveProjectImageRequest struct {
	ProjectId     string   `json:"project_id"`
	ImageUuids    []string `json:"image_uuids"`
	TargetGroupId string   `json:"target_group_id"`
}

type MoveProjectImageResponse struct {
	Images []*ProjectImage `json:"images"`
}

func NewMoveProjectImageResponse(resp *project.MoveProjectImageResponse) *MoveProjectImageResponse {
	if resp == nil {
		return nil
	}
	return &MoveProjectImageResponse{
		Images: NewProjectImages(resp.Images),
	}
}

type ReportProjectImageRequest struct {
	ProjectId string                    `json:"project_id"`
	Images    []*ReportProjectImageItem `json:"images"`
}

type ReportProjectImageItem struct {
	ImageUuid string `json:"image_uuid"`
	Status    string `json:"status"`
}

type ReportProjectImageResponse struct {
	Images []*ProjectImage `json:"images"`
}

func NewReportProjectImageResponse(resp *project.ReportProjectImageResponse) *ReportProjectImageResponse {
	if resp == nil {
		return nil
	}
	return &ReportProjectImageResponse{
		Images: NewProjectImages(resp.Images),
	}
}

type UploadProjectImageRequest struct {
	ProjectId string                    `json:"project_id"`
	GroupId   string                    `json:"group_id"`
	Images    []*UploadProjectImageItem `json:"images"`
}

type UploadProjectImageResponse struct {
	Images []*UploadProjectImageResult `json:"images"`
}

type UploadProjectImageResult struct {
	Image              *ProjectImage `json:"image"`
	OriginalUploadURL  string        `json:"original_upload_url"`
	ThumbnailUploadURL string        `json:"thumbnail_upload_url"`
}

func NewUploadProjectImageResponse(resp *project.UploadProjectImageResponse) *UploadProjectImageResponse {
	if resp == nil {
		return nil
	}
	images := make([]*UploadProjectImageResult, 0, len(resp.Images))
	for _, image := range resp.Images {
		if image == nil {
			continue
		}
		images = append(images, &UploadProjectImageResult{
			Image:              NewProjectImage(image.Image),
			OriginalUploadURL:  image.OriginalUploadURL,
			ThumbnailUploadURL: image.ThumbnailUploadURL,
		})
	}
	return &UploadProjectImageResponse{
		Images: images,
	}
}
