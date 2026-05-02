package profile

import (
	"context"

	"icw_core_biz/internal/rpc_log"
	projectUtils "icw_core_biz/internal/services/project/utils"
	"icw_core_biz/pkg/dto/project"
	"icw_core_biz/pkg/rpc_err"
	"icw_core_biz/utils"
)

// GetProjectProfile 获取项目详情
func (s *Service) GetProjectProfile(req *project.GetProjectProfileRequest, resp *project.GetProjectProfileResponse) (err error) {
	start := rpc_log.Start("ProjectProfileService.GetProjectProfile", req)
	defer func() {
		rpc_log.Finish("ProjectProfileService.GetProjectProfile", req, resp, start, err)
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

	resp.Project = utils.ProjectRecordToDTO(projectRecord)

	return nil
}
