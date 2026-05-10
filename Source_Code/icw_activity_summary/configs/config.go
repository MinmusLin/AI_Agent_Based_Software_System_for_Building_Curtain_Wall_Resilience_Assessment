package configs

import (
	"errors"
	"strings"
	"time"

	"icw_common/env"
)

// Config 服务配置
type Config struct {
	ActivitySummaryAddr                string        `env:"ICW_ACTIVITY_SUMMARY_ADDR"`
	CoreBizAddr                        string        `env:"ICW_CORE_BIZ_ADDR"`
	DetectionSummaryTaskMaxConcurrency int           `env:"DETECTION_SUMMARY_TASK_MAX_CONCURRENCY"`
	DetectionSummaryAgentSecretToken   string        `env:"DETECTION_SUMMARY_AGENT_SECRET_TOKEN"`
	DetectionSummaryAgentBotId         string        `env:"DETECTION_SUMMARY_AGENT_BOT_ID"`
	DetectionSummaryAgentUserId        string        `env:"DETECTION_SUMMARY_AGENT_USER_ID"`
	ProjectSummaryTaskMaxConcurrency   int           `env:"PROJECT_SUMMARY_TASK_MAX_CONCURRENCY"`
	ProjectSummaryAgentSecretToken     string        `env:"PROJECT_SUMMARY_AGENT_SECRET_TOKEN"`
	ProjectSummaryAgentBotId           string        `env:"PROJECT_SUMMARY_AGENT_BOT_ID"`
	ProjectSummaryAgentUserId          string        `env:"PROJECT_SUMMARY_AGENT_USER_ID"`
	AgentRequestTimeout                time.Duration `env:"AGENT_REQUEST_TIMEOUT_MINUTES"`
}

// Validate 校验服务配置
func (cfg *Config) Validate() error {
	var problems []string
	required := []struct {
		key   string
		value string
	}{
		{key: "ICW_ACTIVITY_SUMMARY_ADDR", value: cfg.ActivitySummaryAddr},
		{key: "ICW_CORE_BIZ_ADDR", value: cfg.CoreBizAddr},
		{key: "DETECTION_SUMMARY_AGENT_SECRET_TOKEN", value: cfg.DetectionSummaryAgentSecretToken},
		{key: "DETECTION_SUMMARY_AGENT_BOT_ID", value: cfg.DetectionSummaryAgentBotId},
		{key: "DETECTION_SUMMARY_AGENT_USER_ID", value: cfg.DetectionSummaryAgentUserId},
		{key: "PROJECT_SUMMARY_AGENT_SECRET_TOKEN", value: cfg.ProjectSummaryAgentSecretToken},
		{key: "PROJECT_SUMMARY_AGENT_BOT_ID", value: cfg.ProjectSummaryAgentBotId},
		{key: "PROJECT_SUMMARY_AGENT_USER_ID", value: cfg.ProjectSummaryAgentUserId},
	}
	for _, item := range required {
		if strings.TrimSpace(item.value) == "" {
			problems = append(problems, item.key+" is required")
		}
	}
	if len(problems) > 0 {
		return errors.New(strings.Join(problems, "; "))
	}
	if cfg.DetectionSummaryTaskMaxConcurrency <= 0 {
		return errors.New("DETECTION_SUMMARY_TASK_MAX_CONCURRENCY must be greater than 0")
	}
	if cfg.ProjectSummaryTaskMaxConcurrency <= 0 {
		return errors.New("PROJECT_SUMMARY_TASK_MAX_CONCURRENCY must be greater than 0")
	}
	if cfg.AgentRequestTimeout <= 0 {
		return errors.New("AGENT_REQUEST_TIMEOUT_MINUTES must be greater than 0")
	}
	return nil
}

// Load 加载服务配置
func Load() (Config, error) {
	cfg := Config{
		ActivitySummaryAddr:                env.EnvString("ICW_ACTIVITY_SUMMARY_ADDR"),
		CoreBizAddr:                        env.EnvString("ICW_CORE_BIZ_ADDR"),
		DetectionSummaryTaskMaxConcurrency: env.EnvInt("DETECTION_SUMMARY_TASK_MAX_CONCURRENCY"),
		DetectionSummaryAgentSecretToken:   env.EnvString("DETECTION_SUMMARY_AGENT_SECRET_TOKEN"),
		DetectionSummaryAgentBotId:         env.EnvString("DETECTION_SUMMARY_AGENT_BOT_ID"),
		DetectionSummaryAgentUserId:        env.EnvString("DETECTION_SUMMARY_AGENT_USER_ID"),
		ProjectSummaryTaskMaxConcurrency:   env.EnvInt("PROJECT_SUMMARY_TASK_MAX_CONCURRENCY"),
		ProjectSummaryAgentSecretToken:     env.EnvString("PROJECT_SUMMARY_AGENT_SECRET_TOKEN"),
		ProjectSummaryAgentBotId:           env.EnvString("PROJECT_SUMMARY_AGENT_BOT_ID"),
		ProjectSummaryAgentUserId:          env.EnvString("PROJECT_SUMMARY_AGENT_USER_ID"),
		AgentRequestTimeout:                time.Duration(env.EnvInt("AGENT_REQUEST_TIMEOUT_MINUTES")) * time.Minute,
	}
	if err := cfg.Validate(); err != nil {
		return cfg, err
	}
	return cfg, nil
}
