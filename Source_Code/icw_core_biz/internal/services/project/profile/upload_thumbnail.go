package profile

import (
	"context"

	"icw_core_biz/internal/rpc_log"
	"icw_core_biz/internal/services/project/consts"
	"icw_core_biz/internal/services/project/utils"
	"icw_core_biz/pkg/dto/project"
	"icw_core_biz/pkg/rpc_err"
	"icw_core_biz/repositories/minio"
)

// UploadProjectThumbnail 上传项目缩略图
func (s *Service) UploadProjectThumbnail(req *project.UploadProjectThumbnailRequest, resp *project.UploadProjectThumbnailResponse) (err error) {
	start := rpc_log.Start("ProjectProfileService.UploadProjectThumbnail", req)
	defer func() {
		rpc_log.Finish("ProjectProfileService.UploadProjectThumbnail", req, resp, start, err)
	}()

	if req == nil {
		return rpc_err.BadRequestDefault("request is nil")
	}
	if req.ContentType != consts.ThumbnailContentType {
		return rpc_err.BadRequest(rpc_err.DetailInvalidImageContentType, "image content type must be image/png")
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

	// 返回用户自定义头像上传预签名 URL
	uploadURL, err := s.MinIO().PresignPutObject(ctx, thumbnailKey, s.Config().ProjectThumbnailUploadTTL)
	if err != nil {
		return err
	}
	resp.UploadURL = uploadURL

	return nil
}
