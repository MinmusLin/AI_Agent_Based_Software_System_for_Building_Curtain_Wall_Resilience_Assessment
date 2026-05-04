package rocketmq

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/apache/rocketmq-client-go/v2"
	"github.com/apache/rocketmq-client-go/v2/primitive"
	"github.com/apache/rocketmq-client-go/v2/producer"

	"icw_core_biz/configs"
	"icw_core_biz/pkg/dto/project"
	"icw_core_biz/utils"
)

// Repository RocketMQ 消息队列生产者服务
type Repository struct {
	producer rocketmq.Producer
	topic    string
}

func NewRepository(producer rocketmq.Producer, topic string) *Repository {
	return &Repository{
		producer: producer,
		topic:    strings.TrimSpace(topic),
	}
}

// NewProducer 创建 RocketMQ SDK 生产者
func NewProducer(cfg configs.Config) (rocketmq.Producer, error) {
	messageProducer, err := rocketmq.NewProducer(
		producer.WithNameServer([]string{cfg.RocketMQNamesrvAddr}),
		producer.WithRetry(3),
	)
	if err != nil {
		return nil, err
	}
	if err := messageProducer.Start(); err != nil {
		return nil, err
	}
	return messageProducer, nil
}

// ProducerSendSync 同步发送 RocketMQ 消息
func (r *Repository) ProducerSendSync(ctx context.Context, message *primitive.Message) (result *primitive.SendResult, err error) {
	if r == nil || r.producer == nil || message == nil {
		return nil, nil
	}

	start := time.Now()
	defer func() {
		if utils.IsEmptyError(err) {
			MQSuccess("[%s] cost=%s message=%s result=%s", r.topic, time.Since(start), utils.JSONF(message), utils.JSONF(result))
			return
		}
		MQError("[%s] cost=%s message=%s result=%s err=%s", r.topic, time.Since(start), utils.JSONF(message), utils.JSONF(result), utils.FormatErrorLog(err))
	}()

	return r.producer.SendSync(ctx, message)
}

// PublishProjectImageStatusChangedEvent 发布项目图像状态变化事件
func (r *Repository) PublishProjectImageStatusChangedEvent(ctx context.Context, event *project.ProjectImageStatusChangedEvent) (err error) {
	if r == nil || r.producer == nil || event == nil || event.Image == nil {
		return nil
	}

	body, err := json.Marshal(event)
	if err != nil {
		return err
	}

	message := primitive.NewMessage(r.topic, body)
	message.WithTag(project.EventTagProjectImageStatusChanged)
	message.WithKeys([]string{fmt.Sprintf("%s:%d:%s", event.EventType, event.ProjectId, event.Image.Uuid)})

	_, err = r.ProducerSendSync(ctx, message)
	return err
}
