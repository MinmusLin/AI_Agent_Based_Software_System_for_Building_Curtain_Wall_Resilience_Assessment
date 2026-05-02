package profile

import (
	"context"

	"icw_core_biz/internal/rpc_log"
	"icw_core_biz/internal/services/project/utils"
	"icw_core_biz/pkg/dto/project"
	"icw_core_biz/pkg/rpc_err"
)

// GetProjectThumbnail 获取项目缩略图
func (s *Service) GetProjectThumbnail(req *project.GetProjectThumbnailRequest, resp *project.GetProjectThumbnailResponse) (err error) {
	start := rpc_log.Start("ProjectProfileService.GetProjectThumbnail", req)
	defer func() {
		rpc_log.Finish("ProjectProfileService.GetProjectThumbnail", req, resp, start, err)
	}()

	if req == nil {
		return rpc_err.BadRequestDefault("request is nil")
	}
	ctx := context.Background()

	// 校验用户是否拥有项目访问权限
	if _, err := utils.ValidateProjectOwnership(ctx, s.MySQL(), req.UserId, req.ProjectId); err != nil {
		return err
	}

	// 获取项目缩略图下载预签名 URL
	thumbnailURL, err := utils.PresignProjectThumbnailURL(ctx, s.MinIO(), req.ProjectId, s.Config().ProjectThumbnailGetTTL)
	if err != nil {
		return err
	}
	resp.ThumbnailURL = thumbnailURL

	return nil
}
