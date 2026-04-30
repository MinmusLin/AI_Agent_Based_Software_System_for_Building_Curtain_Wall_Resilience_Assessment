package configs

import (
	"bufio"
	"errors"
	"log"
	"os"
	"strings"
)

// Config 服务配置
type Config struct {
	CoreApiAddr string
	CoreBizAddr string
}

// Validate 校验服务配置
func (cfg *Config) Validate() error {
	var problems []string
	required := []struct {
		key   string
		value string
	}{
		{key: "ICW_CORE_API_ADDR", value: cfg.CoreApiAddr},
		{key: "ICW_CORE_BIZ_ADDR", value: cfg.CoreBizAddr},
	}
	for _, item := range required {
		if strings.TrimSpace(item.value) == "" {
			problems = append(problems, item.key+" is required")
		}
	}
	if len(problems) > 0 {
		return errors.New(strings.Join(problems, "; "))
	}
	return nil
}

// Load 加载服务配置
func Load() (Config, error) {
	cfg := Config{
		CoreApiAddr: env("ICW_CORE_API_ADDR"),
		CoreBizAddr: env("ICW_CORE_BIZ_ADDR"),
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
			_ = os.Setenv(key, value)
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
