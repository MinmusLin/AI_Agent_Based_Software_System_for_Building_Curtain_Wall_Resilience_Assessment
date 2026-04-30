package auth

import (
	"database/sql"

	"github.com/redis/go-redis/v9"

	"icw_core_biz/configs"
	"icw_core_biz/internal/auth/utils"
	"icw_core_biz/repositories"
)

// Service 登录鉴权 Service
type Service struct {
	cfg    configs.Config
	mysql  *repositories.MySQLRepository
	redis  *repositories.RedisRepository
	smtp   *repositories.SMTPRepository
	tokens *utils.TokenManager
}

func NewService(cfg configs.Config, db *sql.DB, rdb *redis.Client) *Service {
	return &Service{
		cfg:    cfg,
		mysql:  repositories.NewMySQLRepository(db),
		redis:  repositories.NewRedisRepository(rdb),
		smtp:   repositories.NewSMTPRepository(cfg),
		tokens: utils.NewTokenManager(cfg.JWTSecret, cfg.AccessTokenTTL),
	}
}
