package core

import (
	"context"

	"icw_core_biz/internal/rpc_log"
	"icw_core_biz/internal/services/project/consts"
	"icw_core_biz/internal/services/project/utils"
	"icw_core_biz/pkg/dto/project"
	"icw_core_biz/pkg/rpc_err"
)

// CreateProject 创建项目
func (s *Service) CreateProject(req *project.CreateProjectRequest, resp *project.CreateProjectResponse) (err error) {
	start := rpc_log.Start("ProjectCoreService.CreateProject", req)
	defer func() {
		rpc_log.Finish("ProjectCoreService.CreateProject", req, resp, start, err)
	}()

	if req == nil {
		return rpc_err.BadRequestDefault("request is nil")
	}
	ctx := context.Background()

	projectRecord, err := s.MySQL().CreateProject(ctx, req.UserId, consts.DefaultProjectName)
	if err != nil {
		return err
	}

	// 获取项目缩略图
	resp.Project, err = utils.ProjectRecordToDTOWithThumbnail(ctx, s.MinIO(), projectRecord, s.Config().ProjectThumbnailGetTTL)
	if err != nil {
		return err
	}

	return nil
}
