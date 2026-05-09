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
	"github.com/apache/rocketmq-client-go/v2/rlog"

	"icw_common/consts"
	"icw_common/gen/core/biz"
	"icw_common/utils"

	"icw_core_biz/configs"
)

// Producer RocketMQ 消息队列生产者服务
type Producer struct {
	producer rocketmq.Producer
	topic    string
}

// NewProducer 创建 RocketMQ 消息队列生产者服务
func NewProducer(cfg configs.Config) (*Producer, error) {
	rlog.SetLogLevel("fatal")
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
	return &Producer{
		producer: messageProducer,
		topic:    strings.TrimSpace(cfg.RocketMQProjectEventTopic),
	}, nil
}

// Shutdown 关闭 RocketMQ 消息队列生产者服务
func (r *Producer) Shutdown() error {
	if r == nil || r.producer == nil {
		return nil
	}
	return r.producer.Shutdown()
}

// producerSendSync 同步发送 RocketMQ 消息
func (r *Producer) producerSendSync(ctx context.Context, message *primitive.Message) (result *primitive.SendResult, err error) {
	if r == nil || r.producer == nil || message == nil {
		return nil, nil
	}

	start := time.Now()
	defer func() {
		if utils.IsEmptyError(err) {
			MQInfo("[PRODUCE|%s] %s %13v %s [%s] status=%d queue_offset=%s offset_msg_id=%s transaction_id=%s",
				r.topic,
				consts.LogColorBoldBlackOnWhite, time.Since(start), consts.LogColorReset,
				result.MsgID,
				result.Status,
				result.QueueOffset,
				result.OffsetMsgID,
				result.TransactionID,
			)
			return
		}
		MQError("[PRODUCE|%s] %s %13v %s [%s] status=%d queue_offset=%s offset_msg_id=%s transaction_id=%s err=%s",
			r.topic,
			consts.LogColorBoldBlackOnWhite, time.Since(start), consts.LogColorReset,
			result.MsgID,
			result.Status,
			result.QueueOffset,
			result.OffsetMsgID,
			result.TransactionID,
			utils.FormatErrorLog(err),
		)
	}()

	return r.producer.SendSync(ctx, message)
}

// PublishProjectImageStatusChangedEvent 发布项目图像状态变化事件
func (r *Producer) PublishProjectImageStatusChangedEvent(ctx context.Context, event *bizpb.ProjectImageStatusChangedEvent) (err error) {
	if r == nil || r.producer == nil || event == nil || event.Image == nil {
		return nil
	}

	body, err := json.Marshal(event)
	if err != nil {
		return err
	}

	message := primitive.NewMessage(r.topic, body)
	message.WithTag(consts.EventTagProjectImageStatusChanged)
	message.WithKeys([]string{fmt.Sprintf("%s:%d:%s", event.EventType, event.ProjectId, event.Image.Uuid)})

	_, err = r.producerSendSync(ctx, message)
	return err
}

// PublishProjectDetectionTaskStatusChangedEvent 发布项目图像检测任务状态变化事件
func (r *Producer) PublishProjectDetectionTaskStatusChangedEvent(ctx context.Context, event *bizpb.ProjectDetectionTaskStatusChangedEvent) (err error) {
	if r == nil || r.producer == nil || event == nil {
		return nil
	}

	body, err := json.Marshal(event)
	if err != nil {
		return err
	}

	message := primitive.NewMessage(r.topic, body)
	message.WithTag(consts.EventTagProjectDetectionTaskStatusChanged)
	message.WithKeys([]string{fmt.Sprintf("%s:%d:%s", event.EventType, event.ProjectId, event.MainTaskUuid)})

	_, err = r.producerSendSync(ctx, message)
	return err
}

// PublishProjectReportStatusChangedEvent 发布项目评估报告状态变化事件
func (r *Producer) PublishProjectReportStatusChangedEvent(ctx context.Context, event *bizpb.ProjectReportStatusChangedEvent) (err error) {
	if r == nil || r.producer == nil || event == nil {
		return nil
	}

	body, err := json.Marshal(event)
	if err != nil {
		return err
	}

	message := primitive.NewMessage(r.topic, body)
	message.WithTag(consts.EventTagProjectReportStatusChanged)
	message.WithKeys([]string{fmt.Sprintf("%s:%d:%s", event.EventType, event.ProjectId, event.ReportUuid)})

	_, err = r.producerSendSync(ctx, message)
	return err
}
