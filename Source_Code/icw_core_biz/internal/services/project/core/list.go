package core

import (
	"context"

	"icw_core_biz/internal/rpc_log"
	"icw_core_biz/internal/services/project/utils"
	"icw_core_biz/pkg/dto/project"
	"icw_core_biz/pkg/rpc_err"
)

// ListProjects 获取项目列表
func (s *Service) ListProjects(req *project.ListProjectsRequest, resp *project.ListProjectsResponse) (err error) {
	start := rpc_log.Start("ProjectCoreService.ListProjects", req)
	defer func() {
		rpc_log.Finish("ProjectCoreService.ListProjects", req, resp, start, err)
	}()

	if req == nil {
		return rpc_err.BadRequestDefault("request is nil")
	}
	ctx := context.Background()

	activeProjects, completedProjects, err := s.MySQL().ListProjects(ctx, req.UserId)
	if err != nil {
		return err
	}

	// 获取项目缩略图
	resp.ActiveProjects, err = utils.ProjectRecordsToListItemsDTOWithThumbnail(ctx, s.MinIO(), activeProjects, s.Config().ProjectThumbnailGetTTL)
	if err != nil {
		return err
	}
	resp.CompletedProjects, err = utils.ProjectRecordsToListItemsDTOWithThumbnail(ctx, s.MinIO(), completedProjects, s.Config().ProjectThumbnailGetTTL)
	if err != nil {
		return err
	}

	return nil
}
