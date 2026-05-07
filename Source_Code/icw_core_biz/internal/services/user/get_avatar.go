package user

import (
	"context"

	"icw_common/gen/core/biz"

	"icw_core_biz/internal/services/user/consts"
	"icw_core_biz/internal/services/user/utils"
	"icw_core_biz/repositories/minio"
)

// GetAvatar 获取用户头像
func (s *Service) GetAvatar(ctx context.Context, req *bizpb.GetAvatarRequest) (*bizpb.GetAvatarResponse, error) {
	resp := &bizpb.GetAvatarResponse{}
	err := s.CallRPC(ctx, req, func() error {
		return s.getAvatar(req, resp)
	})
	return resp, err
}

func (s *Service) getAvatar(req *bizpb.GetAvatarRequest, resp *bizpb.GetAvatarResponse) error {
	resp.AvatarType = consts.AvatarTypeNone

	// 对标准化邮箱地址做 SHA-256 哈希
	emailHash, err := utils.NormalizeEmailHash(req.Email)
	if err != nil {
		return err
	}

	// 获取用户自定义头像下载预签名 URL
	avatarURL, err := minio.PresignCustomAvatarURL(s.Ctx(), s.MinIO(), s.Redis(), emailHash, s.Config().AvatarGetTTL)
	if err != nil {
		return err
	}
	if avatarURL != "" {
		resp.AvatarUrl = avatarURL
		resp.AvatarType = consts.AvatarTypeCustom
		return nil
	}

	// 获取用户默认头像下载预签名 URL
	avatarURL, err = minio.PresignDefaultAvatarURL(s.Ctx(), s.MinIO(), s.Redis(), emailHash, s.Config().AvatarGetTTL)
	if err != nil {
		return err
	}
	if avatarURL != "" {
		resp.AvatarUrl = avatarURL
		resp.AvatarType = consts.AvatarTypeDefault
		return nil
	}

	// 生成用户默认头像
	if err := s.MinIO().PutObject(s.Ctx(), minio.GenDefaultAvatarKey(emailHash), consts.DefaultAvatarContentType, utils.BuildDefaultAvatarSVG(emailHash)); err != nil {
		return err
	}
	avatarURL, err = minio.PresignDefaultAvatarURL(s.Ctx(), s.MinIO(), s.Redis(), emailHash, s.Config().AvatarGetTTL)
	if err != nil {
		return err
	}
	if avatarURL == "" {
		resp.AvatarUrl = ""
		resp.AvatarType = consts.AvatarTypeNone
		return nil
	}

	resp.AvatarUrl = avatarURL
	resp.AvatarType = consts.AvatarTypeDefault
	return nil
}
