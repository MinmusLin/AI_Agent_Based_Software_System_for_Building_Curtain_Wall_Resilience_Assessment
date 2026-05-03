package core

import (
	"icw_core_biz/pkg/dto/project"
	"icw_core_biz/pkg/rpc_err"
	"icw_core_biz/repositories/mysql"
)

// DeleteProject 删除项目
func (s *Service) DeleteProject(req *project.DeleteProjectRequest, resp *project.DeleteProjectResponse) error {
	return s.CallRPC("ProjectCoreService.DeleteProject", req, resp, func() error {
		return s.deleteProject(req, resp)
	})
}

func (s *Service) deleteProject(req *project.DeleteProjectRequest, resp *project.DeleteProjectResponse) error {
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
	resp.ActiveProjects, err = mysql.ProjectRecordsToListItemsDTOWithThumbnail(s.Ctx(), s.MinIO(), s.Redis(), activeProjects, s.Config().ProjectImageGetTTL)
	if err != nil {
		return err
	}
	resp.CompletedProjects, err = mysql.ProjectRecordsToListItemsDTOWithThumbnail(s.Ctx(), s.MinIO(), s.Redis(), completedProjects, s.Config().ProjectImageGetTTL)
	if err != nil {
		return err
	}

	return nil
}
