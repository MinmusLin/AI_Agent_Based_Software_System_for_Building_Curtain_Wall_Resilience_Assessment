package profile

import (
	"icw_core_biz/internal/services/project/consts"
	"icw_core_biz/pkg/dto/project"
	"icw_core_biz/pkg/rpc_err"
	"icw_core_biz/repositories/minio"
)

// UploadProjectThumbnail 上传项目缩略图
func (s *Service) UploadProjectThumbnail(req *project.UploadProjectThumbnailRequest, resp *project.UploadProjectThumbnailResponse) error {
	return s.CallRPC("ProjectProfileService.UploadProjectThumbnail", req, resp, func() error {
		return s.uploadProjectThumbnail(req, resp)
	})
}

func (s *Service) uploadProjectThumbnail(req *project.UploadProjectThumbnailRequest, resp *project.UploadProjectThumbnailResponse) (err error) {
	if req.ContentType != consts.ThumbnailContentType {
		return rpc_err.BadRequest(rpc_err.DetailInvalidImageContentType, "image content type must be image/png")
	}

	// 生成项目缩略图对象 Key
	thumbnailKey, err := minio.GenProjectThumbnailKey(req.ProjectId)
	if err != nil {
		return rpc_err.BadRequestDefault(err.Error())
	}

	// 返回用户自定义头像上传预签名 URL
	uploadURL, err := s.MinIO().PresignPutObject(s.Ctx(), thumbnailKey, s.Config().ProjectThumbnailUploadTTL)
	if err != nil {
		return err
	}
	resp.UploadURL = uploadURL

	return nil
}
