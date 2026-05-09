package dto

import (
	"icw_common/gen/core/api"
	"icw_common/gen/core/biz"
	"icw_common/utils"
)

func NewProject(project *bizpb.Project) *apipb.Project {
	if project == nil {
		return nil
	}
	return &apipb.Project{
		Id:                  utils.Encode(project.Id),
		Name:                project.Name,
		BuildingName:        project.BuildingName,
		BuildingLocation:    project.BuildingLocation,
		BuiltYear:           project.BuiltYear,
		BuildingDescription: project.BuildingDescription,
		KnownIssues:         project.KnownIssues,
		AssessmentGoal:      project.AssessmentGoal,
		ThumbnailUrl:        project.ThumbnailUrl,
		Progress:            project.Progress,
		CreatedAt:           project.CreatedAt,
		UpdatedAt:           project.UpdatedAt,
	}
}

func NewProjectGroup(group *bizpb.ProjectGroup) *apipb.ProjectGroup {
	if group == nil {
		return nil
	}
	return &apipb.ProjectGroup{
		Id:        utils.Encode(group.Id),
		Name:      group.Name,
		SortOrder: group.SortOrder,
		Images:    group.Images,
	}
}

func NewProjectGroups(groups []*bizpb.ProjectGroup) []*apipb.ProjectGroup {
	if groups == nil {
		return make([]*apipb.ProjectGroup, 0)
	}
	items := make([]*apipb.ProjectGroup, 0, len(groups))
	for _, group := range groups {
		if group == nil {
			continue
		}
		items = append(items, NewProjectGroup(group))
	}
	return items
}

func NewProjectListItem(project *bizpb.ProjectListItem) *apipb.ProjectListItem {
	if project == nil {
		return nil
	}
	return &apipb.ProjectListItem{
		Id:               utils.Encode(project.Id),
		Name:             project.Name,
		BuildingName:     project.BuildingName,
		BuildingLocation: project.BuildingLocation,
		ThumbnailUrl:     project.ThumbnailUrl,
		Progress:         project.Progress,
		CreatedAt:        project.CreatedAt,
	}
}

func NewProjectListItems(projects []*bizpb.ProjectListItem) []*apipb.ProjectListItem {
	if projects == nil {
		return make([]*apipb.ProjectListItem, 0)
	}
	items := make([]*apipb.ProjectListItem, 0, len(projects))
	for _, item := range projects {
		if item == nil {
			continue
		}
		items = append(items, NewProjectListItem(item))
	}
	return items
}

func NewUploadProjectImageItem(image *apipb.UploadProjectImageItem) *bizpb.UploadProjectImageItem {
	if image == nil {
		return nil
	}
	return &bizpb.UploadProjectImageItem{
		FileName:    image.FileName,
		ContentType: image.ContentType,
		SizeBytes:   image.SizeBytes,
		Width:       image.Width,
		Height:      image.Height,
		Metadata:    image.Metadata,
	}
}
