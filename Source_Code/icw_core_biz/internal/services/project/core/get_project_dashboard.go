package core

import (
	"context"

	"icw_common/gen/core/biz"

	"icw_core_biz/internal/services/common"
)

// GetProjectDashboard 获取项目工作台统计
func (s *Service) GetProjectDashboard(ctx context.Context, req *bizpb.GetProjectDashboardRequest) (*bizpb.GetProjectDashboardResponse, error) {
	resp := &bizpb.GetProjectDashboardResponse{}
	err := s.CallRPC(req, func() error {
		return s.getProjectDashboard(req, resp)
	})
	return resp, err
}

func (s *Service) getProjectDashboard(req *bizpb.GetProjectDashboardRequest, resp *bizpb.GetProjectDashboardResponse) error {
	stats, err := s.MySQL().GetProjectDashboardStats(s.Ctx(), req.UserId)
	if err != nil {
		return err
	}

	resp.ActiveProjectCount = stats.ActiveProjectCount
	resp.CompletedProjectCount = stats.CompletedProjectCount
	resp.TotalProjectCount = stats.TotalProjectCount
	resp.UploadedImageCount = stats.UploadedImageCount
	resp.ProjectGroupCount = stats.ProjectGroupCount
	resp.DetectionTaskCount = stats.DetectionTaskCount
	resp.CorrosionDetectionTaskCount = stats.CorrosionDetectionTaskCount
	resp.CrackDetectionTaskCount = stats.CrackDetectionTaskCount
	resp.StainDetectionTaskCount = stats.StainDetectionTaskCount
	resp.FlatnessDetectionTaskCount = stats.FlatnessDetectionTaskCount
	resp.SpallingDetectionTaskCount = stats.SpallingDetectionTaskCount
	resp.DetectionSummaryTaskCount = stats.DetectionSummaryTaskCount
	resp.ReportTaskCount = stats.ReportTaskCount

	minioRepo := s.MinIO()
	if minioRepo == nil {
		common.RpcWarn("Get MinIO bucket stats failed, err: minio repository is nil")
		return nil
	}
	bucketStats, err := minioRepo.BucketStats(s.Ctx())
	if err != nil {
		common.RpcWarn("Get MinIO bucket stats failed, err: %v", err)
		return nil
	}

	resp.MinioStorageAvailable = true
	resp.MinioObjectCount = bucketStats.ObjectCount
	resp.MinioBucketUsedBytes = bucketStats.UsedBytes
	resp.MinioBucketQuotaBytes = bucketStats.QuotaBytes
	resp.MinioBucketRemainingBytes = bucketStats.RemainingBytes

	return nil
}
