package user

import (
	"icw_core_biz/configs"
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
