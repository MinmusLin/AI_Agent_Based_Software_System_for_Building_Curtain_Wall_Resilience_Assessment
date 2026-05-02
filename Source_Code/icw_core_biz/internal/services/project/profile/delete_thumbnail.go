package profile

import (
	"context"

	"icw_core_biz/internal/rpc_log"
	"icw_core_biz/internal/services/project/utils"
	"icw_core_biz/pkg/dto/project"
	"icw_core_biz/pkg/rpc_err"
	"icw_core_biz/repositories/minio"
)

// DeleteProjectThumbnail 删除项目缩略图
func (s *Service) DeleteProjectThumbnail(req *project.DeleteProjectThumbnailRequest, resp *project.DeleteProjectThumbnailResponse) (err error) {
	start := rpc_log.Start("ProjectProfileService.DeleteProjectThumbnail", req)
	defer func() {
		rpc_log.Finish("ProjectProfileService.DeleteProjectThumbnail", req, resp, start, err)
	}()

	if req == nil {
		return rpc_err.BadRequestDefault("request is nil")
	}
	ctx := context.Background()

	// 校验用户是否拥有项目访问权限
	if _, err := utils.ValidateProjectOwnership(ctx, s.MySQL(), req.UserId, req.ProjectId); err != nil {
		return err
	}

	// 生成项目缩略图对象 Key
	thumbnailKey, err := minio.GenProjectThumbnailKey(req.ProjectId)
	if err != nil {
		return rpc_err.BadRequestDefault(err.Error())
	}

	// 删除项目缩略图
	if err := s.MinIO().RemoveObject(ctx, thumbnailKey); err != nil {
		return err
	}

	return nil
}
