package agent

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"mime/multipart"
	"net/http"
	"net/textproto"
	"strings"

	"icw_common/utils"
)

const (
	// CozeFileUploadURL Coze 文件上传接口地址
	CozeFileUploadURL = "https://api.coze.cn/v1/files/upload"
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
	// MessageObjectTypeText Coze 文本消息对象类型
	MessageObjectTypeText = "text"
	// MessageObjectTypeFile Coze 文件消息对象类型
	MessageObjectTypeFile = "file"
)

// Client 智能体客户端
type Client struct {
	token  string
	botId  string
	userId string
	client *http.Client
}

// File 智能体输入附件
type File struct {
	Name        string
	Data        []byte
	ContentType string
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
	if err := c.validate(); err != nil {
		return "", err
	}

	return c.streamChat(ctx, strings.TrimSpace(text), "")
}

// ChatWithFile 调用智能体并携带文件附件
func (c *Client) ChatWithFile(ctx context.Context, text string, file File) (string, error) {
	if err := c.validate(); err != nil {
		return "", err
	}

	fileId, err := c.uploadFile(ctx, file)
	if err != nil {
		return "", err
	}

	return c.streamChat(ctx, strings.TrimSpace(text), fileId)
}

func (c *Client) validate() error {
	if c == nil || c.client == nil {
		return errors.New("agent client is nil")
	}
	if c.token == "" {
		return errors.New("agent token is required")
	}
	if c.botId == "" {
		return errors.New("agent bot id is required")
	}
	return nil
}

// uploadFile 上传文件附件并返回文件 ID
func (c *Client) uploadFile(ctx context.Context, file File) (string, error) {
	if len(file.Data) == 0 {
		return "", errors.New("agent file data is required")
	}

	var body bytes.Buffer
	writer := multipart.NewWriter(&body)
	header := make(textproto.MIMEHeader)
	header.Set("Content-Disposition", `form-data; name="file"; filename="`+utils.FirstNotEmpty(file.Name, "source.json")+`"`)
	header.Set("Content-Type", utils.FirstNotEmpty(file.ContentType, "application/octet-stream"))

	part, err := writer.CreatePart(header)
	if err != nil {
		return "", err
	}
	if _, err := part.Write(file.Data); err != nil {
		return "", err
	}
	if err := writer.Close(); err != nil {
		return "", err
	}

	request, err := http.NewRequestWithContext(ctx, http.MethodPost, CozeFileUploadURL, &body)
	if err != nil {
		return "", err
	}
	request.Header.Set("Authorization", "Bearer "+c.token)
	request.Header.Set("Content-Type", writer.FormDataContentType())

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

	var uploadResp struct {
		Data struct {
			Id     string `json:"id"`
			FileId string `json:"file_id"`
		} `json:"data"`
		Id     string `json:"id"`
		FileId string `json:"file_id"`
	}
	if err := json.NewDecoder(response.Body).Decode(&uploadResp); err != nil {
		return "", err
	}

	fileId := utils.FirstNotEmpty(uploadResp.Data.Id, uploadResp.Data.FileId, uploadResp.Id, uploadResp.FileId)
	if fileId == "" {
		return "", errors.New("agent file id is required")
	}

	return fileId, nil
}

// streamChat 发起流式对话并拼接输出
func (c *Client) streamChat(ctx context.Context, text, fileId string) (string, error) {
	messageObjects := []map[string]string{{
		"type": MessageObjectTypeText,
		"text": text,
	}}
	if fileId != "" {
		messageObjects = append(messageObjects, map[string]string{
			"type":    MessageObjectTypeFile,
			"file_id": fileId,
		})
	}

	content, err := json.Marshal(messageObjects)
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
