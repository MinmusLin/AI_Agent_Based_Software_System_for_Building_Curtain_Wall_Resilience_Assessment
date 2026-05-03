package user

import (
	"icw_core_biz/internal/services/user/consts"
	"icw_core_biz/internal/services/user/utils"
	"icw_core_biz/pkg/dto"
	"icw_core_biz/repositories/minio"
)

// GetAvatar 获取用户头像
func (s *Service) GetAvatar(req *dto.GetAvatarRequest, resp *dto.GetAvatarResponse) error {
	return s.CallRPC("UserService.GetAvatar", req, resp, func() error {
		return s.getAvatar(req, resp)
	})
}

func (s *Service) getAvatar(req *dto.GetAvatarRequest, resp *dto.GetAvatarResponse) error {
	resp.AvatarType = consts.AvatarTypeNone

	// 对标准化邮箱地址做 SHA-256 哈希
	emailHash, err := utils.NormalizeEmailHash(req.Email)
	if err != nil {
		return err
	}

	// 获取用户自定义头像下载预签名 URL
	avatarURL, err := minio.PresignCustomAvatarURL(s.Ctx(), s.MinIO(), emailHash, s.Config().AvatarGetTTL)
	if err != nil {
		return err
	}
	if avatarURL != "" {
		resp.AvatarURL = avatarURL
		resp.AvatarType = consts.AvatarTypeCustom
		return nil
	}

	// 获取用户默认头像下载预签名 URL
	avatarURL, err = minio.PresignDefaultAvatarURL(s.Ctx(), s.MinIO(), emailHash, s.Config().AvatarGetTTL)
	if err != nil {
		return err
	}
	if avatarURL != "" {
		resp.AvatarURL = avatarURL
		resp.AvatarType = consts.AvatarTypeDefault
		return nil
	}

	// 生成用户默认头像
	if err := s.MinIO().PutObject(s.Ctx(), minio.GenDefaultAvatarKey(emailHash), consts.DefaultAvatarContentType, utils.BuildDefaultAvatarSVG(emailHash)); err != nil {
		return err
	}
	avatarURL, err = minio.PresignDefaultAvatarURL(s.Ctx(), s.MinIO(), emailHash, s.Config().AvatarGetTTL)
	if err != nil {
		return err
	}
	if avatarURL == "" {
		resp.AvatarURL = ""
		resp.AvatarType = consts.AvatarTypeNone
		return nil
	}

	resp.AvatarURL = avatarURL
	resp.AvatarType = consts.AvatarTypeDefault
	return nil
}
