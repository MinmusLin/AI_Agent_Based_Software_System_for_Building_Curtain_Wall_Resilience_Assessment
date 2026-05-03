package core

import (
	"icw_core_biz/internal/services/project/consts"
	"icw_core_biz/pkg/dto/project"
	"icw_core_biz/repositories/mysql"
)

// CreateProject 创建项目
func (s *Service) CreateProject(req *project.CreateProjectRequest, resp *project.CreateProjectResponse) error {
	return s.CallRPC("ProjectCoreService.CreateProject", req, resp, func() error {
		return s.createProject(req, resp)
	})
}

func (s *Service) createProject(req *project.CreateProjectRequest, resp *project.CreateProjectResponse) error {
	projectRecord, err := s.MySQL().CreateProject(s.Ctx(), req.UserId, consts.DefaultProjectName)
	if err != nil {
		return err
	}

	// 获取项目缩略图
	resp.Project, err = mysql.ProjectRecordToDTOWithThumbnail(s.Ctx(), s.MinIO(), projectRecord, s.Config().ProjectThumbnailGetTTL)
	if err != nil {
		return err
	}

	return nil
}
