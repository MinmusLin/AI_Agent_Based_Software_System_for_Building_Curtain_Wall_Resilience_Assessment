package core

import (
	"icw_core_biz/internal/services/project/utils"
	"icw_core_biz/pkg/dto/project"
)

// ListProjects 获取项目列表
func (s *Service) ListProjects(req *project.ListProjectsRequest, resp *project.ListProjectsResponse) error {
	return s.CallRPC("ProjectCoreService.ListProjects", req, resp, func() error {
		return s.listProjects(req, resp)
	})
}

func (s *Service) listProjects(req *project.ListProjectsRequest, resp *project.ListProjectsResponse) (err error) {
	activeProjects, completedProjects, err := s.MySQL().ListProjects(s.Ctx, req.UserId)
	if err != nil {
		return err
	}

	// 获取项目缩略图
	resp.ActiveProjects, err = utils.ProjectRecordsToListItemsDTOWithThumbnail(s.Ctx, s.MinIO(), activeProjects, s.Config().ProjectThumbnailGetTTL)
	if err != nil {
		return err
	}
	resp.CompletedProjects, err = utils.ProjectRecordsToListItemsDTOWithThumbnail(s.Ctx, s.MinIO(), completedProjects, s.Config().ProjectThumbnailGetTTL)
	if err != nil {
		return err
	}

	return nil
}
