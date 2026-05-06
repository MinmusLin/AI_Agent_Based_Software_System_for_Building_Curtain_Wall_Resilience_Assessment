package profile

import (
	"context"

	"icw_common/gen/core/biz"
	"icw_common/rpc_err"
	"icw_core_biz/repositories/mysql"
)

// GetProjectProfile 获取项目基础信息
func (s *Service) GetProjectProfile(ctx context.Context, req *bizpb.GetProjectProfileRequest) (*bizpb.GetProjectProfileResponse, error) {
	resp := &bizpb.GetProjectProfileResponse{}
	err := s.CallRPC(ctx, req, func() error {
		return s.getProjectProfile(req, resp)
	})
	return resp, err
}

func (s *Service) getProjectProfile(req *bizpb.GetProjectProfileRequest, resp *bizpb.GetProjectProfileResponse) error {
	projectRecord, err := s.MySQL().FindProjectByIdAndUserId(s.Ctx(), req.UserId, req.ProjectId)
	if err != nil {
		return err
	}
	if projectRecord == nil {
		return rpc_err.BadRequest(rpc_err.DetailProjectNotAccessible, "project is not accessible")
	}

	// 获取项目缩略图
	resp.Project, err = mysql.ProjectRecordToDTOWithThumbnail(s.Ctx(), s.MinIO(), s.Redis(), projectRecord, s.Config().ProjectImageGetTTL)
	if err != nil {
		return err
	}

	return nil
}
