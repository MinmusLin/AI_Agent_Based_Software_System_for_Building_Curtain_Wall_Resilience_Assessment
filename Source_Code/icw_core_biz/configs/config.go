package configs

import (
	"bufio"
	"errors"
	"os"
	"strconv"
	"strings"
	"time"

	"icw_core_biz/consts"
	"icw_core_biz/utils"
)

// Config 服务配置
type Config struct {
	CoreBizAddr                string        `env:"ICW_CORE_BIZ_ADDR"`
	ActivityClassificationAddr string        `env:"ICW_ACTIVITY_CLASSIFICATION_ADDR"`
	ActivityReasoningAddr      string        `env:"ICW_ACTIVITY_REASONING_ADDR"`
	ActivitySummaryAddr        string        `env:"ICW_ACTIVITY_SUMMARY_ADDR"`
	MySQLUsername              string        `env:"MYSQL_USERNAME"`
	MySQLPassword              string        `env:"MYSQL_PASSWORD"`
	MySQLAddr                  string        `env:"MYSQL_ADDR"`
	MySQLDatabase              string        `env:"MYSQL_DATABASE"`
	RedisAddr                  string        `env:"REDIS_ADDR"`
	RedisPassword              string        `env:"REDIS_PASSWORD"`
	RedisDB                    int           `env:"REDIS_DB"`
	SMTPHost                   string        `env:"SMTP_HOST"`
	SMTPPort                   int           `env:"SMTP_PORT"`
	SMTPPassword               string        `env:"SMTP_PASSWORD"`
	SMTPFromName               string        `env:"SMTP_FROM_NAME"`
	SMTPFromEmail              string        `env:"SMTP_FROM_EMAIL"`
	MinIOEndpoint              string        `env:"MINIO_ENDPOINT"`
	MinIOAccessKey             string        `env:"MINIO_ACCESS_KEY"`
	MinIOAccessSecret          string        `env:"MINIO_ACCESS_SECRET"`
	MinIOBucket                string        `env:"MINIO_BUCKET"`
	RocketMQNamesrvAddr        string        `env:"ROCKETMQ_NAMESRV_ADDR"`
	RocketMQProjectEventTopic  string        `env:"ROCKETMQ_PROJECT_EVENT_TOPIC"`
	JWTSecret                  string        `env:"JWT_SECRET"`
	EmailCodeSecret            string        `env:"EMAIL_CODE_SECRET"`
	EmailCodeTTL               time.Duration `env:"EMAIL_CODE_TTL_MINUTES"`
	LoginFailTTL               time.Duration `env:"LOGIN_FAIL_TTL_MINUTES"`
	AccessTokenTTL             time.Duration `env:"ACCESS_TOKEN_TTL_MINUTES"`
	RefreshTokenTTL            time.Duration `env:"REFRESH_TOKEN_TTL_MINUTES"`
	AvatarGetTTL               time.Duration `env:"AVATAR_GET_TTL_MINUTES"`
	AvatarUploadTTL            time.Duration `env:"AVATAR_UPLOAD_TTL_MINUTES"`
	ProjectImageGetTTL         time.Duration `env:"PROJECT_IMAGE_GET_TTL_MINUTES"`
	ProjectImageUploadTTL      time.Duration `env:"PROJECT_IMAGE_UPLOAD_TTL_MINUTES"`
	SocketTicketTTL            time.Duration `env:"SOCKET_TICKET_TTL_MINUTES"`
	ProjectImagePendingTimeout time.Duration `env:"PROJECT_IMAGE_PENDING_TIMEOUT_MINUTES"`
	PendingImageTimeoutJobCron string        `env:"PENDING_IMAGE_TIMEOUT_JOB_CRON"`
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
	if len(problems) > 0 {
		return errors.New(strings.Join(problems, "; "))
	}
	return nil
}

// Load 加载服务配置
func Load() (Config, error) {
	cfg := Config{
		CoreBizAddr:                EnvString("ICW_CORE_BIZ_ADDR"),
		ActivityClassificationAddr: EnvString("ICW_ACTIVITY_CLASSIFICATION_ADDR"),
		ActivityReasoningAddr:      EnvString("ICW_ACTIVITY_REASONING_ADDR"),
		ActivitySummaryAddr:        EnvString("ICW_ACTIVITY_SUMMARY_ADDR"),
		MySQLUsername:              EnvString("MYSQL_USERNAME"),
		MySQLPassword:              EnvString("MYSQL_PASSWORD"),
		MySQLAddr:                  EnvString("MYSQL_ADDR"),
		MySQLDatabase:              EnvString("MYSQL_DATABASE"),
		RedisAddr:                  EnvString("REDIS_ADDR"),
		RedisPassword:              EnvString("REDIS_PASSWORD"),
		RedisDB:                    EnvInt("REDIS_DB"),
		SMTPHost:                   EnvString("SMTP_HOST"),
		SMTPPort:                   EnvInt("SMTP_PORT"),
		SMTPPassword:               EnvString("SMTP_PASSWORD"),
		SMTPFromName:               EnvString("SMTP_FROM_NAME"),
		SMTPFromEmail:              EnvString("SMTP_FROM_EMAIL"),
		MinIOEndpoint:              EnvString("MINIO_ENDPOINT"),
		MinIOAccessKey:             EnvString("MINIO_ACCESS_KEY"),
		MinIOAccessSecret:          EnvString("MINIO_ACCESS_SECRET"),
		MinIOBucket:                EnvString("MINIO_BUCKET"),
		RocketMQNamesrvAddr:        EnvString("ROCKETMQ_NAMESRV_ADDR"),
		RocketMQProjectEventTopic:  EnvString("ROCKETMQ_PROJECT_EVENT_TOPIC"),
		JWTSecret:                  EnvString("JWT_SECRET"),
		EmailCodeSecret:            EnvString("EMAIL_CODE_SECRET"),
		EmailCodeTTL:               time.Duration(EnvInt("EMAIL_CODE_TTL_MINUTES")) * time.Minute,
		LoginFailTTL:               time.Duration(EnvInt("LOGIN_FAIL_TTL_MINUTES")) * time.Minute,
		AccessTokenTTL:             time.Duration(EnvInt("ACCESS_TOKEN_TTL_MINUTES")) * time.Minute,
		RefreshTokenTTL:            time.Duration(EnvInt("REFRESH_TOKEN_TTL_MINUTES")) * time.Minute,
		AvatarGetTTL:               time.Duration(EnvInt("AVATAR_GET_TTL_MINUTES")) * time.Minute,
		AvatarUploadTTL:            time.Duration(EnvInt("AVATAR_UPLOAD_TTL_MINUTES")) * time.Minute,
		ProjectImageGetTTL:         time.Duration(EnvInt("PROJECT_IMAGE_GET_TTL_MINUTES")) * time.Minute,
		ProjectImageUploadTTL:      time.Duration(EnvInt("PROJECT_IMAGE_UPLOAD_TTL_MINUTES")) * time.Minute,
		SocketTicketTTL:            time.Duration(EnvInt("SOCKET_TICKET_TTL_MINUTES")) * time.Minute,
		ProjectImagePendingTimeout: time.Duration(EnvInt("PROJECT_IMAGE_PENDING_TIMEOUT_MINUTES")) * time.Minute,
		PendingImageTimeoutJobCron: EnvString("PENDING_IMAGE_TIMEOUT_JOB_CRON"),
	}
	if err := cfg.Validate(); err != nil {
		return cfg, err
	}
	return cfg, nil
}

// LoadDotEnv 从 .env 文件加载环境变量
func LoadDotEnv(path string) {
	file, err := os.Open(path)
	if err != nil {
		utils.LogFatal(consts.LogScopeInit, "Failed to open .env file: %v", err)
	}
	defer func() {
		_ = file.Close()
	}()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		key, value, ok := strings.Cut(line, "=")
		if !ok {
			continue
		}
		key = strings.TrimSpace(key)
		value = strings.Trim(strings.TrimSpace(value), `"'`)
		if key == "" || value == "" {
			continue
		}
		if _, exists := os.LookupEnv(key); !exists {
			if err := os.Setenv(key, value); err != nil {
				utils.LogFatal(consts.LogScopeInit, "Failed to set environment variable %s=%s: %v", key, value, err)
			}
		}
	}

	if err := scanner.Err(); err != nil {
		utils.LogFatal(consts.LogScopeInit, "Failed to read .env file: %v", err)
	}
}

// EnvString 获取环境变量（String 类型）
func EnvString(key string) string {
	if value := strings.TrimSpace(os.Getenv(key)); value != "" {
		return value
	}
	return ""
}

// EnvInt 获取环境变量（Int 类型）
func EnvInt(key string) int {
	value := strings.TrimSpace(os.Getenv(key))
	if value == "" {
		return 0
	}
	parsedInt, err := strconv.Atoi(value)
	if err != nil {
		return 0
	}
	return parsedInt
}
