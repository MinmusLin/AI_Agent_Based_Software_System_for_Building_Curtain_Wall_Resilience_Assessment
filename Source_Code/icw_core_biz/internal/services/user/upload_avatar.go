package user

import (
	"context"

	"icw_core_biz/internal/rpc_log"
	"icw_core_biz/internal/services/user/consts"
	"icw_core_biz/internal/services/user/utils"
	"icw_core_biz/pkg/dto"
	"icw_core_biz/pkg/rpc_err"
)

// UploadAvatar 上传用户自定义头像
func (s *Service) UploadAvatar(req *dto.UploadAvatarRequest, resp *dto.UploadAvatarResponse) (err error) {
	start := rpc_log.Start("UserService.UploadAvatar", req)
	defer func() {
		rpc_log.Finish("UserService.UploadAvatar", req, resp, start, err)
	}()

	if req == nil {
		return rpc_err.BadRequestDefault("request is nil")
	}
	if req.ContentType != consts.CustomAvatarContentType {
		return rpc_err.BadRequest(rpc_err.DetailInvalidAvatarContentType, "avatar content type must be image/png")
	}
	ctx := context.Background()

	// 对标准化邮箱地址做 SHA-256 哈希
	emailHash, err := utils.NormalizeEmailHash(req.Email)
	if err != nil {
		return err
	}

	// 返回用户自定义头像上传预签名 URL
	uploadURL, err := s.MinIO().PresignPutObject(ctx, utils.GenCustomAvatarKey(emailHash), s.Config().AvatarUploadTTL)
	if err != nil {
		return err
	}
	resp.UploadURL = uploadURL

	return nil
}
