package profile

import (
	"icw_core_biz/internal/services/project/utils"
	"icw_core_biz/pkg/dto/project"
	"icw_core_biz/pkg/rpc_err"
)

// GetProjectProfile 获取项目基础信息
func (s *Service) GetProjectProfile(req *project.GetProjectProfileRequest, resp *project.GetProjectProfileResponse) error {
	return s.CallRPC("ProjectProfileService.GetProjectProfile", req, resp, func() error {
		return s.getProjectProfile(req, resp)
	})
}

func (s *Service) getProjectProfile(req *project.GetProjectProfileRequest, resp *project.GetProjectProfileResponse) (err error) {
	projectRecord, err := s.MySQL().FindProjectByIdAndUserId(s.Ctx, req.UserId, req.ProjectId)
	if err != nil {
		return err
	}
	if projectRecord == nil {
		return rpc_err.BadRequest(rpc_err.DetailProjectNotAccessible, "project is not accessible")
	}

	// 获取项目缩略图
	resp.Project, err = utils.ProjectRecordToDTOWithThumbnail(s.Ctx, s.MinIO(), projectRecord, s.Config().ProjectThumbnailGetTTL)
	if err != nil {
		return err
	}

	return nil
}
