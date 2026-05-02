package core

import (
	"context"

	"icw_core_biz/internal/rpc_log"
	"icw_core_biz/pkg/dto/project"
	"icw_core_biz/pkg/rpc_err"
	"icw_core_biz/utils"
)

// ListProjects 获取当前用户项目列表
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

	resp.ActiveProjects = utils.ProjectRecordsToListItemsDTO(activeProjects)
	resp.CompletedProjects = utils.ProjectRecordsToListItemsDTO(completedProjects)

	return nil
}
