package profile

import (
	"context"

	"icw_common/gen/core/biz"
	"icw_common/rpc/error"
	"icw_core_biz/internal/services/project/utils"
	"icw_core_biz/repositories/mysql"
)

// UpdateProjectProfile 更新项目基础信息
func (s *Service) UpdateProjectProfile(ctx context.Context, req *bizpb.UpdateProjectProfileRequest) (*bizpb.UpdateProjectProfileResponse, error) {
	resp := &bizpb.UpdateProjectProfileResponse{}
	err := s.CallRPC(ctx, req, func() error {
		return s.updateProjectProfile(req, resp)
	})
	return resp, err
}

func (s *Service) updateProjectProfile(req *bizpb.UpdateProjectProfileRequest, resp *bizpb.UpdateProjectProfileResponse) error {
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
		uint16(req.BuiltYear),
		fields.BuildingDescription,
		fields.KnownIssues,
		fields.AssessmentGoal,
	)
	if err != nil {
		return err
	}
	if projectRecord == nil {
		return rpc_error.BadRequestDefault("project profile can only be updated in progress 0 and active status")
	}

	// 获取项目缩略图
	resp.Project, err = mysql.ProjectRecordToDTOWithThumbnail(s.Ctx(), s.MinIO(), s.Redis(), projectRecord, s.Config().ProjectImageGetTTL)
	if err != nil {
		return err
	}

	return nil
}
