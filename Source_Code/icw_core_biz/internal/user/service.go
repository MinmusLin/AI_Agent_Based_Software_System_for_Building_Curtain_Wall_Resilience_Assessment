package user

import (
	"context"

	"icw_core_biz/configs"
	"icw_core_biz/internal/user/consts"
	"icw_core_biz/internal/user/utils"
	"icw_core_biz/repositories/minio"
)

// Service 用户业务服务
type Service struct {
	cfg   configs.Config
	minio *minio.Repository
}

func NewService(cfg configs.Config, minioRepo *minio.Repository) *Service {
	return &Service{
		cfg:   cfg,
		minio: minioRepo,
	}
}

// EnsureDefaultAvatar 保证用户存在默认头像
func (s *Service) EnsureDefaultAvatar(ctx context.Context, email string) error {
	emailHash, err := utils.NormalizeEmailHash(email)
	if err != nil {
		return err
	}
	key := utils.GenDefaultAvatarKey(emailHash)
	exists, err := s.minio.StatObject(ctx, key)
	if err != nil {
		return err
	}
	if exists {
		return nil
	}
	return s.minio.PutObject(ctx, key, consts.DefaultAvatarContentType, utils.BuildDefaultAvatarSVG(emailHash))
}
