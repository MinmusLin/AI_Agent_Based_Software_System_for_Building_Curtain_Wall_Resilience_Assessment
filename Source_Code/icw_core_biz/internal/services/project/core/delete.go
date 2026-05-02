package core

import (
	"context"

	"icw_core_biz/internal/rpc_log"
	"icw_core_biz/internal/services/project/utils"
	"icw_core_biz/pkg/dto/project"
	"icw_core_biz/pkg/rpc_err"
)

// DeleteProject 删除项目
func (s *Service) DeleteProject(req *project.DeleteProjectRequest, resp *project.DeleteProjectResponse) (err error) {
	start := rpc_log.Start("ProjectCoreService.DeleteProject", req)
	defer func() {
		rpc_log.Finish("ProjectCoreService.DeleteProject", req, resp, start, err)
	}()

	if req == nil {
		return rpc_err.BadRequestDefault("request is nil")
	}
	ctx := context.Background()

	// 校验用户是否拥有项目访问权限
	if _, err := utils.ValidateProjectOwnership(ctx, s.MySQL(), req.UserId, req.ProjectId); err != nil {
		return err
	}

	deleted, err := s.MySQL().DeleteProject(ctx, req.UserId, req.ProjectId)
	if err != nil {
		return err
	}
	if !deleted {
		return rpc_err.BadRequest(rpc_err.DetailProjectNotAccessible, "project is not accessible")
	}

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
