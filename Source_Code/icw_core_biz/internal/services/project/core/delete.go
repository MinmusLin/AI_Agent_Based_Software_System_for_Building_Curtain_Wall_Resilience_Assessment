package core

import (
	"icw_core_biz/internal/services/project/utils"
	"icw_core_biz/pkg/dto/project"
	"icw_core_biz/pkg/rpc_err"
)

// DeleteProject 删除项目
func (s *Service) DeleteProject(req *project.DeleteProjectRequest, resp *project.DeleteProjectResponse) error {
	return s.CallRPC("ProjectCoreService.DeleteProject", req, resp, func() error {
		return s.deleteProject(req, resp)
	})
}

func (s *Service) deleteProject(req *project.DeleteProjectRequest, resp *project.DeleteProjectResponse) (err error) {
	deleted, err := s.MySQL().DeleteProject(s.Ctx(), req.UserId, req.ProjectId)
	if err != nil {
		return err
	}
	if !deleted {
		return rpc_err.BadRequest(rpc_err.DetailProjectNotAccessible, "project is not accessible")
	}

	activeProjects, completedProjects, err := s.MySQL().ListProjects(s.Ctx(), req.UserId)
	if err != nil {
		return err
	}

	// 获取项目缩略图
	resp.ActiveProjects, err = utils.ProjectRecordsToListItemsDTOWithThumbnail(s.Ctx(), s.MinIO(), activeProjects, s.Config().ProjectThumbnailGetTTL)
	if err != nil {
		return err
	}
	resp.CompletedProjects, err = utils.ProjectRecordsToListItemsDTOWithThumbnail(s.Ctx(), s.MinIO(), completedProjects, s.Config().ProjectThumbnailGetTTL)
	if err != nil {
		return err
	}

	return nil
}
