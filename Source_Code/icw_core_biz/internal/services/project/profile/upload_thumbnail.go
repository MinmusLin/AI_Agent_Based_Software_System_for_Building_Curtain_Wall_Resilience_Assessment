package profile

import (
	"context"

	"icw_common/gen/core/biz"
	"icw_common/rpc_err"
	"icw_core_biz/internal/services/project/consts"
	"icw_core_biz/repositories/minio"
)

// UploadProjectThumbnail 上传项目缩略图
func (s *Service) UploadProjectThumbnail(ctx context.Context, req *bizpb.UploadProjectThumbnailRequest) (*bizpb.UploadProjectThumbnailResponse, error) {
	resp := &bizpb.UploadProjectThumbnailResponse{}
	err := s.CallRPC(ctx, req, func() error {
		return s.uploadProjectThumbnail(req, resp)
	})
	return resp, err
}

func (s *Service) uploadProjectThumbnail(req *bizpb.UploadProjectThumbnailRequest, resp *bizpb.UploadProjectThumbnailResponse) error {
	if req.ContentType != consts.ThumbnailContentType {
		return rpc_err.BadRequest(rpc_err.DetailInvalidImageContentType, "image content type must be image/png")
	}

	// 生成项目缩略图对象 Key
	thumbnailKey, err := minio.GenProjectThumbnailKey(req.ProjectId)
	if err != nil {
		return rpc_err.BadRequestDefault(err.Error())
	}

	// 返回用户自定义头像上传预签名 URL
	uploadURL, err := s.MinIO().PresignPutObject(s.Ctx(), thumbnailKey, s.Config().ProjectImageUploadTTL)
	if err != nil {
		return err
	}
	resp.UploadUrl = uploadURL

	return nil
}
