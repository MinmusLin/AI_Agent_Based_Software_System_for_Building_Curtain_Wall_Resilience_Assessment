package user

import (
	"context"

	"icw_common/gen/core/biz"
	"icw_common/rpc/error"
	"icw_core_biz/internal/services/user/consts"
	"icw_core_biz/internal/services/user/utils"
	"icw_core_biz/repositories/minio"
)

// UploadAvatar 上传用户自定义头像
func (s *Service) UploadAvatar(ctx context.Context, req *bizpb.UploadAvatarRequest) (*bizpb.UploadAvatarResponse, error) {
	resp := &bizpb.UploadAvatarResponse{}
	err := s.CallRPC(ctx, req, func() error {
		return s.uploadAvatar(req, resp)
	})
	return resp, err
}

func (s *Service) uploadAvatar(req *bizpb.UploadAvatarRequest, resp *bizpb.UploadAvatarResponse) error {
	if req.ContentType != consts.CustomAvatarContentType {
		return rpc_error.BadRequest(rpc_error.DetailInvalidImageContentType, "image content type must be image/png")
	}

	// 对标准化邮箱地址做 SHA-256 哈希
	emailHash, err := utils.NormalizeEmailHash(req.Email)
	if err != nil {
		return err
	}

	// 返回用户自定义头像上传预签名 URL
	uploadURL, err := s.MinIO().PresignPutObject(s.Ctx(), minio.GenCustomAvatarKey(emailHash), s.Config().AvatarUploadTTL)
	if err != nil {
		return err
	}
	resp.UploadUrl = uploadURL

	return nil
}
