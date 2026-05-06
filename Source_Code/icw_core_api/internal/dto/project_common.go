package dto

import (
	"icw_common/gen/core/api"
	"icw_common/gen/core/biz"
	"icw_common/utils"
)

// NewProject 将 BIZ 项目结构体转换为 API 项目结构体
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

// NewProjectListItem 将 BIZ 项目列表项结构体转换为 API 项目列表项结构体
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

// NewProjectListItems 将 BIZ 项目列表项切片转换为 API 项目列表项切片
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

// NewProjectImages 将 BIZ 图像切片转换为 API 图像切片
func NewProjectImages(images []*bizpb.ProjectImage) []*apipb.ProjectImage {
	return images
}

// NewProjectGroup 将 BIZ 图像组结构体转换为 API 图像组结构体
func NewProjectGroup(group *bizpb.ProjectGroup) *apipb.ProjectGroup {
	if group == nil {
		return nil
	}
	return &apipb.ProjectGroup{
		Id:        utils.Encode(group.Id),
		Name:      group.Name,
		SortOrder: group.SortOrder,
		Images:    NewProjectImages(group.Images),
	}
}

// NewProjectGroups 将 BIZ 图像组切片转换为 API 图像组切片
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
