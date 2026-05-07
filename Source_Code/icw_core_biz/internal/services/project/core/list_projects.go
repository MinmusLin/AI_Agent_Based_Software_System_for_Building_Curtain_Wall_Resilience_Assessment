package core

import (
	"context"

	"icw_common/gen/core/biz"

	"icw_core_biz/repositories/mysql/model"
)

// ListProjects 获取项目列表
func (s *Service) ListProjects(ctx context.Context, req *bizpb.ListProjectsRequest) (*bizpb.ListProjectsResponse, error) {
	resp := &bizpb.ListProjectsResponse{}
	err := s.CallRPC(ctx, req, func() error {
		return s.listProjects(req, resp)
	})
	return resp, err
}

func (s *Service) listProjects(req *bizpb.ListProjectsRequest, resp *bizpb.ListProjectsResponse) error {
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
