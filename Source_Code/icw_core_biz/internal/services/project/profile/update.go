package profile

import (
	"strings"

	"icw_core_biz/internal/services/project/utils"
	"icw_core_biz/pkg/dto/project"
	"icw_core_biz/pkg/rpc_err"
)

// UpdateProjectProfile 更新项目基础信息
func (s *Service) UpdateProjectProfile(req *project.UpdateProjectProfileRequest, resp *project.UpdateProjectProfileResponse) error {
	return s.CallRPC("ProjectProfileService.UpdateProjectProfile", req, resp, func() error {
		return s.updateProjectProfile(req, resp)
	})
}

func (s *Service) updateProjectProfile(req *project.UpdateProjectProfileRequest, resp *project.UpdateProjectProfileResponse) (err error) {
	projectRecord, err := s.MySQL().UpdateProjectProfile(
		s.Ctx,
		req.UserId,
		req.ProjectId,
		strings.TrimSpace(req.Name),
		strings.TrimSpace(req.BuildingName),
		strings.TrimSpace(req.BuildingLocation),
		req.BuiltYear,
		strings.TrimSpace(req.BuildingDescription),
		strings.TrimSpace(req.KnownIssues),
		strings.TrimSpace(req.AssessmentGoal),
	)
	if err != nil {
		return err
	}
	if projectRecord == nil {
		return rpc_err.BadRequestDefault("project profile can only be updated in progress 0 and active status")
	}

	// 获取项目缩略图
	resp.Project, err = utils.ProjectRecordToDTOWithThumbnail(s.Ctx, s.MinIO(), projectRecord, s.Config().ProjectThumbnailGetTTL)
	if err != nil {
		return err
	}

	return nil
}
