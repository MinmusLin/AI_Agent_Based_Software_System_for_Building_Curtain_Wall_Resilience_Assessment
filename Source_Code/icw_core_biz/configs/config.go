package configs

import (
	"bufio"
	"log"
	"os"
	"strconv"
	"strings"
	"time"
)

// Config 服务配置
type Config struct {
	CoreBizAddr     string
	MySQLDSN        string
	RedisAddr       string
	RedisPassword   string
	RedisDB         int
	SMTPHost        string
	SMTPPort        int
	SMTPPassword    string
	SMTPFromName    string
	SMTPFromEmail   string
	EmailCodeSecret string
	EmailCodeTTL    time.Duration
	JWTSecret       string
	AccessTokenTTL  time.Duration
	RefreshTokenTTL time.Duration
}

// Load 加载服务配置
func Load() Config {
	return Config{
		CoreBizAddr:     env("ICW_CORE_BIZ_ADDR", "127.0.0.1:8001"),
		MySQLDSN:        env("MYSQL_DSN", ""),
		RedisAddr:       env("REDIS_ADDR", "127.0.0.1:6379"),
		RedisPassword:   env("REDIS_PASSWORD", ""),
		RedisDB:         envInt("REDIS_DB", 0),
		SMTPHost:        env("SMTP_HOST", ""),
		SMTPPort:        envInt("SMTP_PORT", 0),
		SMTPPassword:    env("SMTP_PASSWORD", ""),
		SMTPFromName:    env("SMTP_FROM_NAME", ""),
		SMTPFromEmail:   env("SMTP_FROM_EMAIL", ""),
		EmailCodeSecret: env("EMAIL_CODE_SECRET", ""),
		EmailCodeTTL:    time.Duration(envInt("EMAIL_CODE_TTL_MINUTES", 1)) * time.Minute,
		JWTSecret:       env("JWT_SECRET", ""),
		AccessTokenTTL:  time.Duration(envInt("ACCESS_TOKEN_TTL_MINUTES", 1)) * time.Minute,
		RefreshTokenTTL: time.Duration(envInt("REFRESH_TOKEN_TTL_MINUTES", 1)) * time.Minute,
	}
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
			_ = os.Setenv(key, value)
		}
	}
}

// env 获取环境变量（String 类型）
func env(key, fallback string) string {
	if value := strings.TrimSpace(os.Getenv(key)); value != "" {
		return value
	}
	return fallback
}

// envInt 获取环境变量（Int 类型）
func envInt(key string, fallback int) int {
	value := strings.TrimSpace(os.Getenv(key))
	if value == "" {
		return fallback
	}
	parsedInt, err := strconv.Atoi(value)
	if err != nil {
		return fallback
	}
	return parsedInt
}
