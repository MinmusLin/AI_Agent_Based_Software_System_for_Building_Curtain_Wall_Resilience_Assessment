package configs

import (
	"errors"
	"strings"
	"time"

	"icw_common/env"
)

// Config 服务配置
type Config struct {
	ActivityClassificationAddr       string        `env:"ICW_ACTIVITY_CLASSIFICATION_ADDR"`
	CoreBizAddr                      string        `env:"ICW_CORE_BIZ_ADDR"`
	ClassificationTaskMaxConcurrency int           `env:"CLASSIFICATION_TASK_MAX_CONCURRENCY"`
	AgentSecretToken                 string        `env:"AGENT_SECRET_TOKEN"`
	AgentBotId                       string        `env:"AGENT_BOT_ID"`
	AgentUserId                      string        `env:"AGENT_USER_ID"`
	AgentImageSize                   int           `env:"AGENT_IMAGE_SIZE"`
	AgentRequestTimeout              time.Duration `env:"AGENT_REQUEST_TIMEOUT_MINUTES"`
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
		{key: "AGENT_SECRET_TOKEN", value: cfg.AgentSecretToken},
		{key: "AGENT_BOT_ID", value: cfg.AgentBotId},
		{key: "AGENT_USER_ID", value: cfg.AgentUserId},
	}
	for _, item := range required {
		if strings.TrimSpace(item.value) == "" {
			problems = append(problems, item.key+" is required")
		}
	}
	if len(problems) > 0 {
		return errors.New(strings.Join(problems, "; "))
	}
	if cfg.ClassificationTaskMaxConcurrency <= 0 {
		return errors.New("CLASSIFICATION_TASK_MAX_CONCURRENCY must be greater than 0")
	}
	if cfg.AgentImageSize <= 0 {
		return errors.New("AGENT_IMAGE_SIZE must be greater than 0")
	}
	if cfg.AgentRequestTimeout <= 0 {
		return errors.New("AGENT_REQUEST_TIMEOUT_MINUTES must be greater than 0")
	}
	return nil
}

// Load 加载服务配置
func Load() (Config, error) {
	cfg := Config{
		ActivityClassificationAddr:       env.EnvString("ICW_ACTIVITY_CLASSIFICATION_ADDR"),
		CoreBizAddr:                      env.EnvString("ICW_CORE_BIZ_ADDR"),
		ClassificationTaskMaxConcurrency: env.EnvInt("CLASSIFICATION_TASK_MAX_CONCURRENCY"),
		AgentSecretToken:                 env.EnvString("AGENT_SECRET_TOKEN"),
		AgentBotId:                       env.EnvString("AGENT_BOT_ID"),
		AgentUserId:                      env.EnvString("AGENT_USER_ID"),
		AgentImageSize:                   env.EnvInt("AGENT_IMAGE_SIZE"),
		AgentRequestTimeout:              time.Duration(env.EnvInt("AGENT_REQUEST_TIMEOUT_MINUTES")) * time.Minute,
	}
	if err := cfg.Validate(); err != nil {
		return cfg, err
	}
	return cfg, nil
}
