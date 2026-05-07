package configs

import (
	"errors"
	"strings"

	"icw_common/env"
)

// Config 服务配置
type Config struct {
	ActivitySummaryAddr                string `env:"ICW_ACTIVITY_SUMMARY_ADDR"`
	CoreBizAddr                        string `env:"ICW_CORE_BIZ_ADDR"`
	DetectionSummaryTaskMaxConcurrency int    `env:"DETECTION_SUMMARY_TASK_MAX_CONCURRENCY"`
	ProjectSummaryTaskMaxConcurrency   int    `env:"PROJECT_SUMMARY_TASK_MAX_CONCURRENCY"`
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
	return nil
}

// Load 加载服务配置
func Load() (Config, error) {
	cfg := Config{
		ActivitySummaryAddr:                env.EnvString("ICW_ACTIVITY_SUMMARY_ADDR"),
		CoreBizAddr:                        env.EnvString("ICW_CORE_BIZ_ADDR"),
		DetectionSummaryTaskMaxConcurrency: env.EnvInt("DETECTION_SUMMARY_TASK_MAX_CONCURRENCY"),
		ProjectSummaryTaskMaxConcurrency:   env.EnvInt("PROJECT_SUMMARY_TASK_MAX_CONCURRENCY"),
	}
	if err := cfg.Validate(); err != nil {
		return cfg, err
	}
	return cfg, nil
}
