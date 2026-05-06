package rocketmq

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/apache/rocketmq-client-go/v2"
	"github.com/apache/rocketmq-client-go/v2/consumer"
	"github.com/apache/rocketmq-client-go/v2/primitive"
	"github.com/apache/rocketmq-client-go/v2/rlog"

	"icw_common/consts"
	"icw_common/gen/core/biz"
	"icw_common/utils"
	"icw_core_api/configs"
	"icw_core_api/internal/dto"
	"icw_core_api/internal/socket"
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
			Expression: consts.EventTagProjectImageStatusChanged,
		},
		func(_ context.Context, messages ...*primitive.MessageExt) (consumer.ConsumeResult, error) {
			for _, message := range messages {
				if message == nil {
					continue
				}
				start := time.Now()
				if err := dispatchProjectImageStatusChangedEvent(hub, message); err != nil {
					MQError("[CONSUME|%s] cost=%s tag=%s msg_id=%s err=%s", message.Topic, time.Since(start), message.GetTags(), message.MsgId, utils.FormatErrorLog(err))
					continue
				}
				MQInfo("[CONSUME|%s] cost=%s tag=%s msg_id=%s", message.Topic, time.Since(start), message.GetTags(), message.MsgId)
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

// dispatchProjectImageStatusChangedEvent 分发项目图像状态变化事件
func dispatchProjectImageStatusChangedEvent(hub *socket.Hub, message *primitive.MessageExt) error {
	if hub == nil {
		return errors.New("websocket hub is nil")
	}
	var event bizpb.ProjectImageStatusChangedEvent
	if err := json.Unmarshal(message.Body, &event); err != nil {
		return err
	}
	if event.EventType != consts.EventTypeProjectImageStatusChanged {
		return fmt.Errorf("unexpected event type: %s", event.EventType)
	}
	socketMessage := dto.NewProjectImageStatusChangedMessage(&event)
	messageBytes, err := json.Marshal(socketMessage)
	if err != nil {
		return err
	}
	hub.BroadcastProject(event.ProjectId, event.ProjectCode, messageBytes)
	return nil
}
