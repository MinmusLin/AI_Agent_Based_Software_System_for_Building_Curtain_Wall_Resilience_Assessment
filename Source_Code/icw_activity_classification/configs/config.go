package configs

import (
	"errors"
	"strings"

	"icw_common/env"
)

// Config 服务配置
type Config struct {
	ActivityClassificationAddr string `env:"ICW_ACTIVITY_CLASSIFICATION_ADDR"`
	CoreBizAddr                string `env:"ICW_CORE_BIZ_ADDR"`
}

// Validate 校验服务配置
func (cfg *Config) Validate() error {
	var problems []string
	required := []struct {
		key   string
		value string
	}{
		{key: "ICW_ACTIVITY_CLASSIFICATION_ADDR", value: cfg.ActivityClassificationAddr},
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
		ActivityClassificationAddr: env.EnvString("ICW_ACTIVITY_CLASSIFICATION_ADDR"),
		CoreBizAddr:                env.EnvString("ICW_CORE_BIZ_ADDR"),
	}
	if err := cfg.Validate(); err != nil {
		return cfg, err
	}
	return cfg, nil
}
