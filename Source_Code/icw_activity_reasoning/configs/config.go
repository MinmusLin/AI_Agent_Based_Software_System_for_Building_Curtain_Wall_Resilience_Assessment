package configs

import (
	"errors"
	"strings"
	"time"

	"icw_core_biz/configs"
)

// Config 服务配置
type Config struct {
	ActivityReasoningAddr       string        `env:"ICW_ACTIVITY_REASONING_ADDR"`
	CoreBizAddr                 string        `env:"ICW_CORE_BIZ_ADDR"`
	PythonBin                   string        `env:"PYTHON_BIN"`
	ReasoningWorkDir            string        `env:"REASONING_WORK_DIR"`
	ReasoningTaskCodes          []string      `env:"REASONING_TASK_CODES"`
	ReasoningTaskMaxConcurrency int           `env:"REASONING_TASK_MAX_CONCURRENCY"`
	ReasoningTaskTimeout        time.Duration `env:"REASONING_TASK_TIMEOUT_MINUTES"`
	ArtifactDownloadTimeout     time.Duration `env:"ARTIFACT_DOWNLOAD_TIMEOUT_MINUTES"`
	ArtifactUploadTimeout       time.Duration `env:"ARTIFACT_UPLOAD_TIMEOUT_MINUTES"`
}

// Validate 校验服务配置
func (cfg *Config) Validate() error {
	var problems []string
	required := []struct {
		key   string
		value string
	}{
		{key: "ICW_ACTIVITY_REASONING_ADDR", value: cfg.ActivityReasoningAddr},
		{key: "ICW_CORE_BIZ_ADDR", value: cfg.CoreBizAddr},
		{key: "PYTHON_BIN", value: cfg.PythonBin},
		{key: "REASONING_WORK_DIR", value: cfg.ReasoningWorkDir},
	}
	for _, item := range required {
		if strings.TrimSpace(item.value) == "" {
			problems = append(problems, item.key+" is required")
		}
	}
	if len(cfg.ReasoningTaskCodes) == 0 {
		problems = append(problems, "REASONING_TASK_CODES is required")
	}
	if cfg.ReasoningTaskMaxConcurrency <= 0 {
		problems = append(problems, "REASONING_TASK_MAX_CONCURRENCY must be greater than 0")
	}
	if cfg.ReasoningTaskTimeout <= 0 {
		problems = append(problems, "REASONING_TASK_TIMEOUT_MINUTES must be greater than 0")
	}
	if cfg.ArtifactDownloadTimeout <= 0 {
		problems = append(problems, "ARTIFACT_DOWNLOAD_TIMEOUT_MINUTES must be greater than 0")
	}
	if cfg.ArtifactUploadTimeout <= 0 {
		problems = append(problems, "ARTIFACT_UPLOAD_TIMEOUT_MINUTES must be greater than 0")
	}
	if len(problems) > 0 {
		return errors.New(strings.Join(problems, "; "))
	}
	return nil
}

// Load 加载服务配置
func Load() (Config, error) {
	cfg := Config{
		ActivityReasoningAddr:       configs.EnvString("ICW_ACTIVITY_REASONING_ADDR"),
		CoreBizAddr:                 configs.EnvString("ICW_CORE_BIZ_ADDR"),
		PythonBin:                   configs.EnvString("PYTHON_BIN"),
		ReasoningWorkDir:            configs.EnvString("REASONING_WORK_DIR"),
		ReasoningTaskCodes:          configs.EnvStringSlice("REASONING_TASK_CODES"),
		ReasoningTaskMaxConcurrency: configs.EnvInt("REASONING_TASK_MAX_CONCURRENCY"),
		ReasoningTaskTimeout:        time.Duration(configs.EnvInt("REASONING_TASK_TIMEOUT_MINUTES")) * time.Minute,
		ArtifactDownloadTimeout:     time.Duration(configs.EnvInt("ARTIFACT_DOWNLOAD_TIMEOUT_MINUTES")) * time.Minute,
		ArtifactUploadTimeout:       time.Duration(configs.EnvInt("ARTIFACT_UPLOAD_TIMEOUT_MINUTES")) * time.Minute,
	}
	if err := cfg.Validate(); err != nil {
		return cfg, err
	}
	return cfg, nil
}
