package core

import (
	"context"

	"icw_common/gen/core/biz"
	"icw_common/rpc/error"

	"icw_core_biz/repositories/mysql/model"
)

// DeleteProject 删除项目
func (s *Service) DeleteProject(ctx context.Context, req *bizpb.DeleteProjectRequest) (*bizpb.DeleteProjectResponse, error) {
	resp := &bizpb.DeleteProjectResponse{}
	err := s.CallRPC(ctx, req, func() error {
		return s.deleteProject(req, resp)
	})
	return resp, err
}

func (s *Service) deleteProject(req *bizpb.DeleteProjectRequest, resp *bizpb.DeleteProjectResponse) error {
	deleted, err := s.MySQL().DeleteProject(s.Ctx(), req.UserId, req.ProjectId)
	if err != nil {
		return err
	}
	if !deleted {
		return rpc_error.BadRequest(rpc_error.DetailProjectNotAccessible, "project is not accessible")
	}

	activeProjects, completedProjects, err := s.MySQL().ListProjects(s.Ctx(), req.UserId)
	if err != nil {
		return err
	}

	// 获取项目缩略图
	resp.ActiveProjects, err = model.ProjectRecordsToListItemsDTOWithThumbnail(s.Ctx(), s.MinIO(), s.Redis(), activeProjects, s.Config().ProjectImageGetTTL)
	if err != nil {
		return err
	}
	resp.CompletedProjects, err = model.ProjectRecordsToListItemsDTOWithThumbnail(s.Ctx(), s.MinIO(), s.Redis(), completedProjects, s.Config().ProjectImageGetTTL)
	if err != nil {
		return err
	}

	return nil
}
