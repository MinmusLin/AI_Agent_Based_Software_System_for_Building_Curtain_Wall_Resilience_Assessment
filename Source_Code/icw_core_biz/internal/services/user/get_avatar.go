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

	// 判断用户是否存在自定义头像
	customKey := minio.GenCustomAvatarKey(emailHash)
	customExists, err := s.MinIO().StatObject(s.Ctx(), customKey)
	if err != nil {
		return err
	}

	if customExists {
		// 返回用户自定义头像下载预签名 URL
		avatarURL, err := s.MinIO().PresignGetObject(s.Ctx(), customKey, s.Config().AvatarGetTTL)
		if err != nil {
			return err
		}
		resp.AvatarURL = avatarURL
		resp.AvatarType = consts.AvatarTypeCustom

		return nil
	}

	// 判断用户是否存在默认头像
	defaultKey := minio.GenDefaultAvatarKey(emailHash)
	defaultExists, err := s.MinIO().StatObject(s.Ctx(), defaultKey)
	if err != nil {
		return err
	}

	if !defaultExists {
		// 生成用户默认头像
		if err := s.MinIO().PutObject(s.Ctx(), defaultKey, consts.DefaultAvatarContentType, utils.BuildDefaultAvatarSVG(emailHash)); err != nil {
			return err
		}
	}

	// 返回用户默认头像下载预签名 URL
	avatarURL, err := s.MinIO().PresignGetObject(s.Ctx(), minio.GenDefaultAvatarKey(emailHash), s.Config().AvatarGetTTL)
	if err != nil {
		return err
	}
	resp.AvatarURL = avatarURL
	resp.AvatarType = consts.AvatarTypeDefault

	return nil
}
