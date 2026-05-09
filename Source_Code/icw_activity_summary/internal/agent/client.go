package agent

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"strings"
)

const (
	// CozeChatURL Coze 对话接口地址
	CozeChatURL = "https://api.coze.cn/v3/chat"
	// SSEDataPrefix SSE 数据行前缀
	SSEDataPrefix = "data:"
	// SSEEventPrefix SSE 事件行前缀
	SSEEventPrefix = "event:"
	// SSEEventDone SSE 结束事件标识
	SSEEventDone = "done"
	// ChatEventDelta Coze 对话增量消息事件
	ChatEventDelta = "conversation.message.delta"
	// ChatEventComplete Coze 对话完成消息事件
	ChatEventComplete = "conversation.message.completed"
	// ChatEventFailed Coze 对话失败事件
	ChatEventFailed = "conversation.chat.failed"
	// MessageTypeAnswer Coze 智能体回答消息类型
	MessageTypeAnswer = "answer"
)

// Client 智能体客户端
type Client struct {
	token  string
	botId  string
	userId string
	client *http.Client
}

// NewClient 创建智能体客户端
func NewClient(token, botId, userId string) *Client {
	return &Client{
		token:  token,
		botId:  botId,
		userId: userId,
		client: http.DefaultClient,
	}
}

// Chat 调用智能体并返回模型输出
func (c *Client) Chat(ctx context.Context, text string) (string, error) {
	if c == nil || c.client == nil {
		return "", errors.New("agent client is nil")
	}
	if c.token == "" {
		return "", errors.New("agent token is required")
	}
	if c.botId == "" {
		return "", errors.New("agent bot id is required")
	}

	return c.streamChat(ctx, strings.TrimSpace(text))
}

// streamChat 发起流式对话并拼接输出
func (c *Client) streamChat(ctx context.Context, text string) (string, error) {
	content, err := json.Marshal([]map[string]string{{
		"type": "text",
		"text": text,
	}})
	if err != nil {
		return "", err
	}

	requestBody, err := json.Marshal(map[string]interface{}{
		"bot_id":  c.botId,
		"user_id": c.userId,
		"stream":  true,
		"additional_messages": []map[string]string{
			{
				"role":         "user",
				"content":      string(content),
				"content_type": "object_string",
			},
		},
	})
	if err != nil {
		return "", err
	}

	request, err := http.NewRequestWithContext(ctx, http.MethodPost, CozeChatURL, bytes.NewReader(requestBody))
	if err != nil {
		return "", err
	}

	request.Header.Set("Authorization", "Bearer "+c.token)
	request.Header.Set("Content-Type", "application/json")

	response, err := c.client.Do(request)
	if err != nil {
		return "", err
	}

	defer func() {
		_ = response.Body.Close()
	}()

	if response.StatusCode < http.StatusOK || response.StatusCode >= http.StatusMultipleChoices {
		return "", readAgentError(response.Body, response.Status)
	}

	return readChatStream(response.Body)
}
