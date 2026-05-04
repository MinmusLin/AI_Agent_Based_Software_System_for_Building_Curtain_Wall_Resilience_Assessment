package rocketmq

import (
	"context"
	"encoding/json"

	"github.com/apache/rocketmq-client-go/v2"
	"github.com/apache/rocketmq-client-go/v2/consumer"
	"github.com/apache/rocketmq-client-go/v2/primitive"
	"github.com/apache/rocketmq-client-go/v2/rlog"

	"icw_core_api/configs"
	"icw_core_api/internal/dto"
	"icw_core_api/internal/socket"
	"icw_core_biz/pkg/dto/project"
)

// Consumer RocketMQ 消息队列消费者
type Consumer struct {
	consumer rocketmq.PushConsumer
}

// NewConsumer 创建 RocketMQ SDK 消费者
func NewConsumer(cfg configs.Config, hub *socket.Hub) (*Consumer, error) {
	rlog.SetLogLevel("fatal")
	messageConsumer, err := rocketmq.NewPushConsumer(
		consumer.WithGroupName(cfg.RocketMQConsumerGroup),
		consumer.WithNameServer([]string{cfg.RocketMQNamesrvAddr}),
	)
	if err != nil {
		return nil, err
	}

	if err := messageConsumer.Subscribe(
		cfg.RocketMQProjectEventTopic,
		consumer.MessageSelector{
			Type:       consumer.TAG,
			Expression: project.EventTagProjectImageStatusChanged,
		},
		func(_ context.Context, messages ...*primitive.MessageExt) (consumer.ConsumeResult, error) {
			for _, message := range messages {
				dispatchProjectImageStatusChangedEvent(hub, message)
			}
			return consumer.ConsumeSuccess, nil
		},
	); err != nil {
		return nil, err
	}

	return &Consumer{
		consumer: messageConsumer,
	}, nil
}

// Start 启动 RocketMQ 消费者
func (c *Consumer) Start() error {
	if c == nil || c.consumer == nil {
		return nil
	}
	if err := c.consumer.Start(); err != nil {
		return err
	}
	return nil
}

// Close 关闭 RocketMQ 消费者
func (c *Consumer) Close() error {
	if c == nil || c.consumer == nil {
		return nil
	}
	if err := c.consumer.Shutdown(); err != nil {
		return err
	}
	return nil
}

func dispatchProjectImageStatusChangedEvent(hub *socket.Hub, message *primitive.MessageExt) {
	if hub == nil || message == nil {
		return
	}
	var event project.ProjectImageStatusChangedEvent
	if err := json.Unmarshal(message.Body, &event); err != nil {
		return
	}
	if event.EventType != project.EventTypeProjectImageStatusChanged {
		return
	}
	socketMessage := dto.NewProjectImageStatusChangedMessage(&event)
	messageBytes, err := json.Marshal(socketMessage)
	if err != nil {
		return
	}
	hub.BroadcastProject(event.ProjectId, messageBytes)
}
