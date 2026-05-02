package core

import (
	"context"

	"icw_core_biz/internal/rpc_log"
	"icw_core_biz/internal/services/project/consts"
	projectUtils "icw_core_biz/internal/services/project/utils"
	"icw_core_biz/pkg/dto/project"
	"icw_core_biz/pkg/rpc_err"
	"icw_core_biz/utils"
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
	thumbnailURL, err := projectUtils.PresignProjectThumbnailURL(ctx, s.MinIO(), projectRecord.Id, s.Config().ProjectThumbnailGetTTL)
	if err != nil {
		return err
	}

	resp.Project = utils.ProjectRecordToDTO(projectRecord)
	resp.Project.ThumbnailURL = thumbnailURL

	return nil
}
