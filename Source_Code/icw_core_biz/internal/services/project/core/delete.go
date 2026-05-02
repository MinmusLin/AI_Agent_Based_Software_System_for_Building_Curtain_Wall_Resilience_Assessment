package core

import (
	"context"

	"icw_core_biz/internal/rpc_log"
	projectUtils "icw_core_biz/internal/services/project/utils"
	"icw_core_biz/pkg/dto/project"
	"icw_core_biz/pkg/rpc_err"
	"icw_core_biz/utils"
)

// DeleteProject 软删除项目并返回最新项目列表
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
	if _, err := projectUtils.ValidateProjectOwnership(ctx, s.MySQL(), req.UserId, req.ProjectId); err != nil {
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

	resp.ActiveProjects = utils.ProjectRecordsToListItemsDTO(activeProjects)
	resp.CompletedProjects = utils.ProjectRecordsToListItemsDTO(completedProjects)

	return nil
}
