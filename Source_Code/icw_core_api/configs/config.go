package configs

import (
	"errors"
	"strings"

	"icw_core_biz/configs"
)

// Config 服务配置
type Config struct {
	GinMode                   string `env:"GIN_MODE"`
	CoreApiAddr               string `env:"ICW_CORE_API_ADDR"`
	CoreBizAddr               string `env:"ICW_CORE_BIZ_ADDR"`
	RocketMQNamesrvAddr       string `env:"ROCKETMQ_NAMESRV_ADDR"`
	RocketMQProjectEventTopic string `env:"ROCKETMQ_PROJECT_EVENT_TOPIC"`
	RocketMQConsumerGroup     string `env:"ROCKETMQ_CONSUMER_GROUP"`
}

// Validate 校验服务配置
func (cfg *Config) Validate() error {
	var problems []string
	required := []struct {
		key   string
		value string
	}{
		{key: "GIN_MODE", value: cfg.GinMode},
		{key: "ICW_CORE_API_ADDR", value: cfg.CoreApiAddr},
		{key: "ICW_CORE_BIZ_ADDR", value: cfg.CoreBizAddr},
		{key: "ROCKETMQ_NAMESRV_ADDR", value: cfg.RocketMQNamesrvAddr},
		{key: "ROCKETMQ_PROJECT_EVENT_TOPIC", value: cfg.RocketMQProjectEventTopic},
		{key: "ROCKETMQ_CONSUMER_GROUP", value: cfg.RocketMQConsumerGroup},
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
		GinMode:                   configs.EnvString("GIN_MODE"),
		CoreApiAddr:               configs.EnvString("ICW_CORE_API_ADDR"),
		CoreBizAddr:               configs.EnvString("ICW_CORE_BIZ_ADDR"),
		RocketMQNamesrvAddr:       configs.EnvString("ROCKETMQ_NAMESRV_ADDR"),
		RocketMQProjectEventTopic: configs.EnvString("ROCKETMQ_PROJECT_EVENT_TOPIC"),
		RocketMQConsumerGroup:     configs.EnvString("ROCKETMQ_CONSUMER_GROUP"),
	}
	if err := cfg.Validate(); err != nil {
		return cfg, err
	}
	return cfg, nil
}
