package configs

import (
	"bufio"
	"errors"
	"log"
	"os"
	"strconv"
	"strings"
	"time"
)

// Config 服务配置
type Config struct {
	CoreBizAddr               string
	MySQLDSN                  string
	RedisAddr                 string
	RedisPassword             string
	RedisDB                   int
	SMTPHost                  string
	SMTPPort                  int
	SMTPPassword              string
	SMTPFromName              string
	SMTPFromEmail             string
	MinIOEndpoint             string
	MinIOAccessKey            string
	MinIOAccessSecret         string
	MinIOBucket               string
	JWTSecret                 string
	EmailCodeSecret           string
	EmailCodeTTL              time.Duration
	LoginFailTTL              time.Duration
	AccessTokenTTL            time.Duration
	RefreshTokenTTL           time.Duration
	AvatarGetTTL              time.Duration
	AvatarUploadTTL           time.Duration
	ProjectThumbnailGetTTL    time.Duration
	ProjectThumbnailUploadTTL time.Duration
}

// Validate 校验服务配置
func (cfg *Config) Validate() error {
	var problems []string
	required := []struct {
		key   string
		value string
	}{
		{key: "ICW_CORE_BIZ_ADDR", value: cfg.CoreBizAddr},
		{key: "MYSQL_DSN", value: cfg.MySQLDSN},
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
		{key: "JWT_SECRET", value: cfg.JWTSecret},
		{key: "EMAIL_CODE_SECRET", value: cfg.EmailCodeSecret},
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
	if cfg.ProjectThumbnailGetTTL <= 0 {
		problems = append(problems, "PROJECT_THUMBNAIL_GET_TTL_MINUTES must be greater than 0")
	}
	if cfg.ProjectThumbnailUploadTTL <= 0 {
		problems = append(problems, "PROJECT_THUMBNAIL_UPLOAD_TTL_MINUTES must be greater than 0")
	}
	if len(problems) > 0 {
		return errors.New(strings.Join(problems, "; "))
	}
	return nil
}

// Load 加载服务配置
func Load() (Config, error) {
	cfg := Config{
		CoreBizAddr:               env("ICW_CORE_BIZ_ADDR"),
		MySQLDSN:                  env("MYSQL_DSN"),
		RedisAddr:                 env("REDIS_ADDR"),
		RedisPassword:             env("REDIS_PASSWORD"),
		RedisDB:                   envInt("REDIS_DB"),
		SMTPHost:                  env("SMTP_HOST"),
		SMTPPort:                  envInt("SMTP_PORT"),
		SMTPPassword:              env("SMTP_PASSWORD"),
		SMTPFromName:              env("SMTP_FROM_NAME"),
		SMTPFromEmail:             env("SMTP_FROM_EMAIL"),
		MinIOEndpoint:             env("MINIO_ENDPOINT"),
		MinIOAccessKey:            env("MINIO_ACCESS_KEY"),
		MinIOAccessSecret:         env("MINIO_ACCESS_SECRET"),
		MinIOBucket:               env("MINIO_BUCKET"),
		JWTSecret:                 env("JWT_SECRET"),
		EmailCodeSecret:           env("EMAIL_CODE_SECRET"),
		EmailCodeTTL:              time.Duration(envInt("EMAIL_CODE_TTL_MINUTES")) * time.Minute,
		LoginFailTTL:              time.Duration(envInt("LOGIN_FAIL_TTL_MINUTES")) * time.Minute,
		AccessTokenTTL:            time.Duration(envInt("ACCESS_TOKEN_TTL_MINUTES")) * time.Minute,
		RefreshTokenTTL:           time.Duration(envInt("REFRESH_TOKEN_TTL_MINUTES")) * time.Minute,
		AvatarGetTTL:              time.Duration(envInt("AVATAR_GET_TTL_MINUTES")) * time.Minute,
		AvatarUploadTTL:           time.Duration(envInt("AVATAR_UPLOAD_TTL_MINUTES")) * time.Minute,
		ProjectThumbnailGetTTL:    time.Duration(envInt("PROJECT_THUMBNAIL_GET_TTL_MINUTES")) * time.Minute,
		ProjectThumbnailUploadTTL: time.Duration(envInt("PROJECT_THUMBNAIL_UPLOAD_TTL_MINUTES")) * time.Minute,
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
		log.Fatalf("Failed to open .env file: %v", err)
	}
	defer func(file *os.File) {
		_ = file.Close()
	}(file)

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
				log.Fatalf("Failed to set environment variable %s=%s: %v", key, value, err)
			}
		}
	}

	if err := scanner.Err(); err != nil {
		log.Fatalf("Failed to read .env file: %v", err)
	}
}

// env 获取环境变量（String 类型）
func env(key string) string {
	if value := strings.TrimSpace(os.Getenv(key)); value != "" {
		return value
	}
	return ""
}

// envInt 获取环境变量（Int 类型）
func envInt(key string) int {
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
