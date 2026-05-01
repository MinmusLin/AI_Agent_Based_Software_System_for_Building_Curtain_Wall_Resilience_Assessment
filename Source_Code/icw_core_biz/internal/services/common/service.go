package common

import (
	"icw_core_biz/configs"
	"icw_core_biz/repositories/minio"
	"icw_core_biz/repositories/mysql"
	"icw_core_biz/repositories/redis"
	"icw_core_biz/repositories/smtp"
)

// Deps RPC Service 的公共依赖集合
type Deps struct {
	Config configs.Config
	MySQL  *mysql.Repository
	Redis  *redis.Repository
	SMTP   *smtp.Repository
	MinIO  *minio.Repository
}

func NewDeps(Config configs.Config, MySQL *mysql.Repository, Redis *redis.Repository, SMTP *smtp.Repository, MinIO *minio.Repository) *Deps {
	return &Deps{
		Config: Config,
		MySQL:  MySQL,
		Redis:  Redis,
		SMTP:   SMTP,
		MinIO:  MinIO,
	}
}

// BaseService 提供所有 RPC Service 共享的基础依赖
type BaseService struct {
	deps *Deps
}

func NewBaseService(deps *Deps) *BaseService {
	if deps == nil {
		deps = &Deps{}
	}
	return &BaseService{
		deps: deps,
	}
}

// Config 获取服务配置
func (s *BaseService) Config() configs.Config {
	return s.deps.Config
}

// MySQL 获取 MySQL 服务
func (s *BaseService) MySQL() *mysql.Repository {
	return s.deps.MySQL
}

// Redis 获取 Redis 服务
func (s *BaseService) Redis() *redis.Repository {
	return s.deps.Redis
}

// SMTP 获取 SMTP 服务
func (s *BaseService) SMTP() *smtp.Repository {
	return s.deps.SMTP
}

// MinIO 获取 MinIO 服务
func (s *BaseService) MinIO() *minio.Repository {
	return s.deps.MinIO
}
