package core

import (
	"context"

	"icw_common/gen/core/biz"
	"icw_core_biz/internal/services/project/consts"
	"icw_core_biz/repositories/mysql"
)

// CreateProject 创建项目
func (s *Service) CreateProject(ctx context.Context, req *bizpb.CreateProjectRequest) (*bizpb.CreateProjectResponse, error) {
	resp := &bizpb.CreateProjectResponse{}
	err := s.CallRPC(ctx, req, resp, func() error {
		return s.createProject(req, resp)
	})
	return resp, err
}

func (s *Service) createProject(req *bizpb.CreateProjectRequest, resp *bizpb.CreateProjectResponse) error {
	projectRecord, err := s.MySQL().CreateProject(s.Ctx(), req.UserId, consts.DefaultProjectName)
	if err != nil {
		return err
	}

	// 获取项目缩略图
	resp.Project, err = mysql.ProjectRecordToDTOWithThumbnail(s.Ctx(), s.MinIO(), s.Redis(), projectRecord, s.Config().ProjectImageGetTTL)
	if err != nil {
		return err
	}

	return nil
}
