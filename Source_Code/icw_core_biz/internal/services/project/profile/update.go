package profile

import (
	"icw_core_biz/internal/services/project/utils"
	"icw_core_biz/pkg/dto/project"
	"icw_core_biz/pkg/rpc_err"
	"icw_core_biz/repositories/mysql"
)

// UpdateProjectProfile 更新项目基础信息
func (s *Service) UpdateProjectProfile(req *project.UpdateProjectProfileRequest, resp *project.UpdateProjectProfileResponse) error {
	return s.CallRPC(req, resp, func() error {
		return s.updateProjectProfile(req, resp)
	})
}

func (s *Service) updateProjectProfile(req *project.UpdateProjectProfileRequest, resp *project.UpdateProjectProfileResponse) error {
	fields, err := utils.NormalizeProjectProfileFields(
		req.Name,
		req.BuildingName,
		req.BuildingLocation,
		req.BuildingDescription,
		req.KnownIssues,
		req.AssessmentGoal,
	)
	if err != nil {
		return err
	}

	projectRecord, err := s.MySQL().UpdateProjectProfile(
		s.Ctx(),
		req.UserId,
		req.ProjectId,
		fields.Name,
		fields.BuildingName,
		fields.BuildingLocation,
		req.BuiltYear,
		fields.BuildingDescription,
		fields.KnownIssues,
		fields.AssessmentGoal,
	)
	if err != nil {
		return err
	}
	if projectRecord == nil {
		return rpc_err.BadRequestDefault("project profile can only be updated in progress 0 and active status")
	}

	// 获取项目缩略图
	resp.Project, err = mysql.ProjectRecordToDTOWithThumbnail(s.Ctx(), s.MinIO(), s.Redis(), projectRecord, s.Config().ProjectImageGetTTL)
	if err != nil {
		return err
	}

	return nil
}
