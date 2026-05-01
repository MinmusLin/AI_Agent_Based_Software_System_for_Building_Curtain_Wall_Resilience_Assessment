package user

import (
	"context"
	"icw_core_biz/internal/user/consts"

	"icw_core_biz/internal/rpc_log"
	"icw_core_biz/internal/user/utils"
	"icw_core_biz/pkg/dto"
	"icw_core_biz/pkg/rpc_err"
)

// GetAvatar 获取用户头像
func (s *Service) GetAvatar(req *dto.GetAvatarRequest, resp *dto.GetAvatarResponse) (err error) {
	start := rpc_log.Start("UserService.GetAvatar", req)
	defer func() {
		rpc_log.Finish("UserService.GetAvatar", req, resp, start, err)
	}()

	if req == nil {
		return rpc_err.BadRequestDefault("request is nil")
	}
	ctx := context.Background()

	// 对标准化邮箱地址做 SHA-256 哈希
	emailHash, err := utils.NormalizeEmailHash(req.Email)
	if err != nil {
		return err
	}

	// 判断用户是否存在自定义头像
	customKey := utils.GenCustomAvatarKey(emailHash)
	customExists, err := s.minio.StatObject(ctx, customKey)
	if err != nil {
		return err
	}

	if customExists {
		// 返回用户自定义头像下载预签名 URL
		avatarURL, err := s.minio.PresignGetObject(ctx, customKey, s.cfg.AvatarGetTTL)
		if err != nil {
			return err
		}
		resp.AvatarURL = avatarURL
		return nil
	}

	// 判断用户是否存在默认头像
	defaultKey := utils.GenDefaultAvatarKey(emailHash)
	defaultExists, err := s.minio.StatObject(ctx, defaultKey)
	if err != nil {
		return err
	}

	if defaultExists {
		// 返回用户默认头像下载预签名 URL
		avatarURL, err := s.minio.PresignGetObject(ctx, utils.GenDefaultAvatarKey(emailHash), s.cfg.AvatarGetTTL)
		if err != nil {
			return err
		}
		resp.AvatarURL = avatarURL
		return nil
	}

	// 生成用户默认头像
	return s.minio.PutObject(ctx, defaultKey, consts.DefaultAvatarContentType, utils.BuildDefaultAvatarSVG(emailHash))
}
