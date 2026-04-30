package auth

import (
	"database/sql"

	goredis "github.com/redis/go-redis/v9"

	"icw_core_biz/configs"
	"icw_core_biz/internal/auth/utils"
	"icw_core_biz/repositories/mysql"
	"icw_core_biz/repositories/redis"
	"icw_core_biz/repositories/smtp"
)

// Service 登录鉴权 Service
type Service struct {
	cfg    configs.Config
	mysql  *mysql.MySQLRepository
	redis  *redis.RedisRepository
	smtp   *smtp.SMTPRepository
	tokens *utils.TokenManager
}

func NewService(cfg configs.Config, db *sql.DB, rdb *goredis.Client) *Service {
	return &Service{
		cfg:    cfg,
		mysql:  mysql.NewMySQLRepository(db),
		redis:  redis.NewRedisRepository(rdb),
		smtp:   smtp.NewSMTPRepository(cfg),
		tokens: utils.NewTokenManager(cfg.JWTSecret, cfg.AccessTokenTTL),
	}
}
