package rocketmq

import (
	"context"
	"strings"
	"time"

	"github.com/apache/rocketmq-client-go/v2"
	"github.com/apache/rocketmq-client-go/v2/consumer"
	"github.com/apache/rocketmq-client-go/v2/primitive"
	"github.com/apache/rocketmq-client-go/v2/rlog"

	"icw_common/consts"
	"icw_common/utils"

	"icw_core_api/configs"
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
			Type: consumer.TAG,
			Expression: strings.Join([]string{
				consts.EventTagProjectImageStatusChanged,
				consts.EventTagProjectDetectionTaskStatusChanged,
			}, "||"),
		},
		func(_ context.Context, messages ...*primitive.MessageExt) (consumer.ConsumeResult, error) {
			for _, message := range messages {
				if message == nil {
					continue
				}
				start := time.Now()
				if err := dispatchProjectEvent(hub, message); err != nil {
					MQError("[CONSUME|%s] %s %13v %s msg_id=%s tag=%s err=%s",
						message.Topic,
						consts.LogColorBoldBlackOnWhite, time.Since(start), consts.LogColorReset,
						message.MsgId,
						message.GetTags(),
						utils.FormatErrorLog(err),
					)
					continue
				}
				MQInfo("[CONSUME|%s] %s %13v %s msg_id=%s tag=%s",
					message.Topic,
					consts.LogColorBoldBlackOnWhite, time.Since(start), consts.LogColorReset,
					message.MsgId,
					message.GetTags(),
				)
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
