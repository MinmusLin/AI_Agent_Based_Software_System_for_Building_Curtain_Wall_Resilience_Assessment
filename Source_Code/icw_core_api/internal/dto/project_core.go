package dto

import (
	"icw_common/gen/core/api"
	"icw_common/gen/core/biz"
)

func NewAdvanceProjectResponse(_ *bizpb.AdvanceProjectResponse) *apipb.AdvanceProjectResponse {
	return &apipb.AdvanceProjectResponse{}
}

func NewCreateProjectResponse(resp *bizpb.CreateProjectResponse) *apipb.CreateProjectResponse {
	if resp == nil {
		return nil
	}
	return &apipb.CreateProjectResponse{
		Project: NewProject(resp.Project),
	}
}

func NewDeleteProjectResponse(resp *bizpb.DeleteProjectResponse) *apipb.DeleteProjectResponse {
	if resp == nil {
		return nil
	}
	return &apipb.DeleteProjectResponse{
		ActiveProjects:    NewProjectListItems(resp.ActiveProjects),
		CompletedProjects: NewProjectListItems(resp.CompletedProjects),
	}
}

func NewGetProjectDashboardResponse(resp *bizpb.GetProjectDashboardResponse) *apipb.GetProjectDashboardResponse {
	if resp == nil {
		return nil
	}
	return &apipb.GetProjectDashboardResponse{
		ActiveProjectCount:          resp.ActiveProjectCount,
		CompletedProjectCount:       resp.CompletedProjectCount,
		TotalProjectCount:           resp.TotalProjectCount,
		UploadedImageCount:          resp.UploadedImageCount,
		ProjectGroupCount:           resp.ProjectGroupCount,
		MinioObjectCount:            resp.MinioObjectCount,
		DetectionTaskCount:          resp.DetectionTaskCount,
		CorrosionDetectionTaskCount: resp.CorrosionDetectionTaskCount,
		CrackDetectionTaskCount:     resp.CrackDetectionTaskCount,
		StainDetectionTaskCount:     resp.StainDetectionTaskCount,
		FlatnessDetectionTaskCount:  resp.FlatnessDetectionTaskCount,
		SpallingDetectionTaskCount:  resp.SpallingDetectionTaskCount,
		DetectionSummaryTaskCount:   resp.DetectionSummaryTaskCount,
		ReportTaskCount:             resp.ReportTaskCount,
		MinioBucketUsedBytes:        resp.MinioBucketUsedBytes,
		MinioBucketQuotaBytes:       resp.MinioBucketQuotaBytes,
		MinioBucketRemainingBytes:   resp.MinioBucketRemainingBytes,
		MinioStorageAvailable:       resp.MinioStorageAvailable,
	}
}

func NewListProjectsResponse(resp *bizpb.ListProjectsResponse) *apipb.ListProjectsResponse {
	if resp == nil {
		return nil
	}
	return &apipb.ListProjectsResponse{
		ActiveProjects:    NewProjectListItems(resp.ActiveProjects),
		CompletedProjects: NewProjectListItems(resp.CompletedProjects),
	}
}
