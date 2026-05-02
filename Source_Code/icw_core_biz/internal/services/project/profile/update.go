package profile

import (
	"context"
	"strings"

	"icw_core_biz/internal/rpc_log"
	"icw_core_biz/internal/services/project/consts"
	projectUtils "icw_core_biz/internal/services/project/utils"
	"icw_core_biz/pkg/dto/project"
	"icw_core_biz/pkg/rpc_err"
	"icw_core_biz/utils"
)

// UpdateProjectProfile 更新项目基础信息
func (s *Service) UpdateProjectProfile(req *project.UpdateProjectProfileRequest, resp *project.UpdateProjectProfileResponse) (err error) {
	start := rpc_log.Start("ProjectProfileService.UpdateProjectProfile", req)
	defer func() {
		rpc_log.Finish("ProjectProfileService.UpdateProjectProfile", req, resp, start, err)
	}()

	if req == nil {
		return rpc_err.BadRequestDefault("request is nil")
	}
	ctx := context.Background()

	// 校验用户是否拥有项目访问权限
	projectRecord, err := projectUtils.ValidateProjectOwnership(ctx, s.MySQL(), req.UserId, req.ProjectId)
	if err != nil {
		return err
	}

	// 只有项目状态为进行中且项目进度为项目基础信息阶段时，可以更新项目基础信息
	if projectRecord.Progress != consts.ProjectProgressInitializationFinished || projectRecord.Status != consts.ProjectStatusActive {
		return rpc_err.BadRequestDefault("project profile can only be updated in progress 0 and active status")
	}

	projectRecord, err = s.MySQL().UpdateProjectProfile(
		ctx,
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
	thumbnailURL, err := projectUtils.PresignProjectThumbnailURL(ctx, s.MinIO(), projectRecord.Id, s.Config().ProjectThumbnailGetTTL)
	if err != nil {
		return err
	}

	resp.Project = utils.ProjectRecordToDTO(projectRecord)
	resp.Project.ThumbnailURL = thumbnailURL

	return nil
}
