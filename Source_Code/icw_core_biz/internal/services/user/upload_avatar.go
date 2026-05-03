package user

import (
	"icw_core_biz/internal/services/user/consts"
	"icw_core_biz/internal/services/user/utils"
	"icw_core_biz/pkg/dto"
	"icw_core_biz/pkg/rpc_err"
	"icw_core_biz/repositories/minio"
)

// UploadAvatar 上传用户自定义头像
func (s *Service) UploadAvatar(req *dto.UploadAvatarRequest, resp *dto.UploadAvatarResponse) error {
	return s.CallRPC("UserService.UploadAvatar", req, resp, func() error {
		return s.uploadAvatar(req, resp)
	})
}

func (s *Service) uploadAvatar(req *dto.UploadAvatarRequest, resp *dto.UploadAvatarResponse) (err error) {
	if req.ContentType != consts.CustomAvatarContentType {
		return rpc_err.BadRequest(rpc_err.DetailInvalidImageContentType, "image content type must be image/png")
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
	resp.UploadURL = uploadURL

	return nil
}
