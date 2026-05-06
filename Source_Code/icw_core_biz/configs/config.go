package configs

import (
	"errors"
	"strings"
	"time"

	"icw_common/env"
)

// Config 服务配置
type Config struct {
	CoreBizAddr                   string        `env:"ICW_CORE_BIZ_ADDR"`
	ActivityClassificationAddr    string        `env:"ICW_ACTIVITY_CLASSIFICATION_ADDR"`
	ActivityReasoningAddr         string        `env:"ICW_ACTIVITY_REASONING_ADDR"`
	ActivitySummaryAddr           string        `env:"ICW_ACTIVITY_SUMMARY_ADDR"`
	MySQLUsername                 string        `env:"MYSQL_USERNAME"`
	MySQLPassword                 string        `env:"MYSQL_PASSWORD"`
	MySQLAddr                     string        `env:"MYSQL_ADDR"`
	MySQLDatabase                 string        `env:"MYSQL_DATABASE"`
	RedisAddr                     string        `env:"REDIS_ADDR"`
	RedisPassword                 string        `env:"REDIS_PASSWORD"`
	RedisDB                       int           `env:"REDIS_DB"`
	SMTPHost                      string        `env:"SMTP_HOST"`
	SMTPPort                      int           `env:"SMTP_PORT"`
	SMTPPassword                  string        `env:"SMTP_PASSWORD"`
	SMTPFromName                  string        `env:"SMTP_FROM_NAME"`
	SMTPFromEmail                 string        `env:"SMTP_FROM_EMAIL"`
	MinIOEndpoint                 string        `env:"MINIO_ENDPOINT"`
	MinIOAccessKey                string        `env:"MINIO_ACCESS_KEY"`
	MinIOAccessSecret             string        `env:"MINIO_ACCESS_SECRET"`
	MinIOBucket                   string        `env:"MINIO_BUCKET"`
	RocketMQNamesrvAddr           string        `env:"ROCKETMQ_NAMESRV_ADDR"`
	RocketMQProjectEventTopic     string        `env:"ROCKETMQ_PROJECT_EVENT_TOPIC"`
	JWTSecret                     string        `env:"JWT_SECRET"`
	EmailCodeSecret               string        `env:"EMAIL_CODE_SECRET"`
	EmailCodeTTL                  time.Duration `env:"EMAIL_CODE_TTL_MINUTES"`
	LoginFailTTL                  time.Duration `env:"LOGIN_FAIL_TTL_MINUTES"`
	AccessTokenTTL                time.Duration `env:"ACCESS_TOKEN_TTL_MINUTES"`
	RefreshTokenTTL               time.Duration `env:"REFRESH_TOKEN_TTL_MINUTES"`
	AvatarGetTTL                  time.Duration `env:"AVATAR_GET_TTL_MINUTES"`
	AvatarUploadTTL               time.Duration `env:"AVATAR_UPLOAD_TTL_MINUTES"`
	ProjectImageGetTTL            time.Duration `env:"PROJECT_IMAGE_GET_TTL_MINUTES"`
	ProjectImageUploadTTL         time.Duration `env:"PROJECT_IMAGE_UPLOAD_TTL_MINUTES"`
	SocketTicketTTL               time.Duration `env:"SOCKET_TICKET_TTL_MINUTES"`
	ProjectImagePendingTimeout    time.Duration `env:"PROJECT_IMAGE_PENDING_TIMEOUT_MINUTES"`
	PendingImageTimeoutJobCron    string        `env:"PENDING_IMAGE_TIMEOUT_JOB_CRON"`
	DetectionWorkerMaxConcurrency int           `env:"DETECTION_WORKER_MAX_CONCURRENCY"`
}

// Validate 校验服务配置
func (cfg *Config) Validate() error {
	var problems []string
	required := []struct {
		key   string
		value string
	}{
		{key: "ICW_CORE_BIZ_ADDR", value: cfg.CoreBizAddr},
		{key: "ICW_ACTIVITY_CLASSIFICATION_ADDR", value: cfg.ActivityClassificationAddr},
		{key: "ICW_ACTIVITY_REASONING_ADDR", value: cfg.ActivityReasoningAddr},
		{key: "ICW_ACTIVITY_SUMMARY_ADDR", value: cfg.ActivitySummaryAddr},
		{key: "MYSQL_USERNAME", value: cfg.MySQLUsername},
		{key: "MYSQL_PASSWORD", value: cfg.MySQLPassword},
		{key: "MYSQL_ADDR", value: cfg.MySQLAddr},
		{key: "MYSQL_DATABASE", value: cfg.MySQLDatabase},
		{key: "REDIS_ADDR", value: cfg.RedisAddr},
		{key: "REDIS_PASSWORD", value: cfg.RedisPassword},
		{key: "SMTP_HOST", value: cfg.SMTPHost},
		{key: "SMTP_PASSWORD", value: cfg.SMTPPassword},
		{key: "SMTP_FROM_NAME", value: cfg.SMTPFromName},
		{key: "SMTP_FROM_EMAIL", value: cfg.SMTPFromEmail},
		{key: "MINIO_ENDPOINT", value: cfg.MinIOEndpoint},
		{key: "MINIO_ACCESS_KEY", value: cfg.MinIOAccessKey},
		{key: "MINIO_ACCESS_SECRET", value: cfg.MinIOAccessSecret},
		{key: "MINIO_BUCKET", value: cfg.MinIOBucket},
		{key: "ROCKETMQ_NAMESRV_ADDR", value: cfg.RocketMQNamesrvAddr},
		{key: "ROCKETMQ_PROJECT_EVENT_TOPIC", value: cfg.RocketMQProjectEventTopic},
		{key: "JWT_SECRET", value: cfg.JWTSecret},
		{key: "EMAIL_CODE_SECRET", value: cfg.EmailCodeSecret},
		{key: "PENDING_IMAGE_TIMEOUT_JOB_CRON", value: cfg.PendingImageTimeoutJobCron},
	}
	for _, item := range required {
		if strings.TrimSpace(item.value) == "" {
			problems = append(problems, item.key+" is required")
		}
	}
	if cfg.RedisDB < 0 {
		problems = append(problems, "REDIS_DB must be greater than or equal to 0")
	}
	if cfg.SMTPPort <= 0 {
		problems = append(problems, "SMTP_PORT must be greater than 0")
	}
	if cfg.EmailCodeTTL <= 0 {
		problems = append(problems, "EMAIL_CODE_TTL_MINUTES must be greater than 0")
	}
	if cfg.LoginFailTTL <= 0 {
		problems = append(problems, "LOGIN_FAIL_TTL_MINUTES must be greater than 0")
	}
	if cfg.AccessTokenTTL <= 0 {
		problems = append(problems, "ACCESS_TOKEN_TTL_MINUTES must be greater than 0")
	}
	if cfg.RefreshTokenTTL <= 0 {
		problems = append(problems, "REFRESH_TOKEN_TTL_MINUTES must be greater than 0")
	}
	if cfg.AvatarGetTTL <= 0 {
		problems = append(problems, "AVATAR_GET_TTL_MINUTES must be greater than 0")
	}
	if cfg.AvatarUploadTTL <= 0 {
		problems = append(problems, "AVATAR_UPLOAD_TTL_MINUTES must be greater than 0")
	}
	if cfg.ProjectImageGetTTL <= 0 {
		problems = append(problems, "PROJECT_IMAGE_GET_TTL_MINUTES must be greater than 0")
	}
	if cfg.ProjectImageUploadTTL <= 0 {
		problems = append(problems, "PROJECT_IMAGE_UPLOAD_TTL_MINUTES must be greater than 0")
	}
	if cfg.SocketTicketTTL <= 0 {
		problems = append(problems, "SOCKET_TICKET_TTL_MINUTES must be greater than 0")
	}
	if cfg.ProjectImagePendingTimeout <= 0 {
		problems = append(problems, "PROJECT_IMAGE_PENDING_TIMEOUT_MINUTES must be greater than 0")
	}
	if cfg.DetectionWorkerMaxConcurrency <= 0 {
		problems = append(problems, "DETECTION_WORKER_MAX_CONCURRENCY must be greater than 0")
	}
	if len(problems) > 0 {
		return errors.New(strings.Join(problems, "; "))
	}
	return nil
}

// Load 加载服务配置
func Load() (Config, error) {
	cfg := Config{
		CoreBizAddr:                   env.EnvString("ICW_CORE_BIZ_ADDR"),
		ActivityClassificationAddr:    env.EnvString("ICW_ACTIVITY_CLASSIFICATION_ADDR"),
		ActivityReasoningAddr:         env.EnvString("ICW_ACTIVITY_REASONING_ADDR"),
		ActivitySummaryAddr:           env.EnvString("ICW_ACTIVITY_SUMMARY_ADDR"),
		MySQLUsername:                 env.EnvString("MYSQL_USERNAME"),
		MySQLPassword:                 env.EnvString("MYSQL_PASSWORD"),
		MySQLAddr:                     env.EnvString("MYSQL_ADDR"),
		MySQLDatabase:                 env.EnvString("MYSQL_DATABASE"),
		RedisAddr:                     env.EnvString("REDIS_ADDR"),
		RedisPassword:                 env.EnvString("REDIS_PASSWORD"),
		RedisDB:                       env.EnvInt("REDIS_DB"),
		SMTPHost:                      env.EnvString("SMTP_HOST"),
		SMTPPort:                      env.EnvInt("SMTP_PORT"),
		SMTPPassword:                  env.EnvString("SMTP_PASSWORD"),
		SMTPFromName:                  env.EnvString("SMTP_FROM_NAME"),
		SMTPFromEmail:                 env.EnvString("SMTP_FROM_EMAIL"),
		MinIOEndpoint:                 env.EnvString("MINIO_ENDPOINT"),
		MinIOAccessKey:                env.EnvString("MINIO_ACCESS_KEY"),
		MinIOAccessSecret:             env.EnvString("MINIO_ACCESS_SECRET"),
		MinIOBucket:                   env.EnvString("MINIO_BUCKET"),
		RocketMQNamesrvAddr:           env.EnvString("ROCKETMQ_NAMESRV_ADDR"),
		RocketMQProjectEventTopic:     env.EnvString("ROCKETMQ_PROJECT_EVENT_TOPIC"),
		JWTSecret:                     env.EnvString("JWT_SECRET"),
		EmailCodeSecret:               env.EnvString("EMAIL_CODE_SECRET"),
		EmailCodeTTL:                  time.Duration(env.EnvInt("EMAIL_CODE_TTL_MINUTES")) * time.Minute,
		LoginFailTTL:                  time.Duration(env.EnvInt("LOGIN_FAIL_TTL_MINUTES")) * time.Minute,
		AccessTokenTTL:                time.Duration(env.EnvInt("ACCESS_TOKEN_TTL_MINUTES")) * time.Minute,
		RefreshTokenTTL:               time.Duration(env.EnvInt("REFRESH_TOKEN_TTL_MINUTES")) * time.Minute,
		AvatarGetTTL:                  time.Duration(env.EnvInt("AVATAR_GET_TTL_MINUTES")) * time.Minute,
		AvatarUploadTTL:               time.Duration(env.EnvInt("AVATAR_UPLOAD_TTL_MINUTES")) * time.Minute,
		ProjectImageGetTTL:            time.Duration(env.EnvInt("PROJECT_IMAGE_GET_TTL_MINUTES")) * time.Minute,
		ProjectImageUploadTTL:         time.Duration(env.EnvInt("PROJECT_IMAGE_UPLOAD_TTL_MINUTES")) * time.Minute,
		SocketTicketTTL:               time.Duration(env.EnvInt("SOCKET_TICKET_TTL_MINUTES")) * time.Minute,
		ProjectImagePendingTimeout:    time.Duration(env.EnvInt("PROJECT_IMAGE_PENDING_TIMEOUT_MINUTES")) * time.Minute,
		PendingImageTimeoutJobCron:    env.EnvString("PENDING_IMAGE_TIMEOUT_JOB_CRON"),
		DetectionWorkerMaxConcurrency: env.EnvInt("DETECTION_WORKER_MAX_CONCURRENCY"),
	}
	if err := cfg.Validate(); err != nil {
		return cfg, err
	}
	return cfg, nil
}
