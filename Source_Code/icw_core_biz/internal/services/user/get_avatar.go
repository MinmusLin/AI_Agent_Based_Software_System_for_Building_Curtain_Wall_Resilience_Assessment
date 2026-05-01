package user

import (
	"context"

	"icw_core_biz/internal/rpc_log"
	"icw_core_biz/internal/services/user/consts"
	"icw_core_biz/internal/services/user/utils"
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
	resp.AvatarType = consts.AvatarTypeNone

	// 对标准化邮箱地址做 SHA-256 哈希
	emailHash, err := utils.NormalizeEmailHash(req.Email)
	if err != nil {
		return err
	}

	// 判断用户是否存在自定义头像
	customKey := utils.GenCustomAvatarKey(emailHash)
	customExists, err := s.MinIO().StatObject(ctx, customKey)
	if err != nil {
		return err
	}

	if customExists {
		// 返回用户自定义头像下载预签名 URL
		avatarURL, err := s.MinIO().PresignGetObject(ctx, customKey, s.Config().AvatarGetTTL)
		if err != nil {
			return err
		}
		resp.AvatarURL = avatarURL
		resp.AvatarType = consts.AvatarTypeCustom

		return nil
	}

	// 判断用户是否存在默认头像
	defaultKey := utils.GenDefaultAvatarKey(emailHash)
	defaultExists, err := s.MinIO().StatObject(ctx, defaultKey)
	if err != nil {
		return err
	}

	if !defaultExists {
		// 生成用户默认头像
		if err := s.MinIO().PutObject(ctx, defaultKey, consts.DefaultAvatarContentType, utils.BuildDefaultAvatarSVG(emailHash)); err != nil {
			return err
		}
	}

	// 返回用户默认头像下载预签名 URL
	avatarURL, err := s.MinIO().PresignGetObject(ctx, utils.GenDefaultAvatarKey(emailHash), s.Config().AvatarGetTTL)
	if err != nil {
		return err
	}
	resp.AvatarURL = avatarURL
	resp.AvatarType = consts.AvatarTypeDefault

	return nil
}
