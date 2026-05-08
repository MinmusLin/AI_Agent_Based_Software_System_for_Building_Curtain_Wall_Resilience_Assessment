package agent

import (
	"bufio"
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"io"
	"mime/multipart"
	"net/http"
	"net/textproto"
	"strings"
)

const (
	CozeFileUploadURL = "https://api.coze.cn/v1/files/upload"
	CozeChatURL       = "https://api.coze.cn/v3/chat"
	SSEDataPrefix     = "data:"
	SSEEventPrefix    = "event:"
	SSEEventDone      = "done"
	ChatEventDelta    = "conversation.message.delta"
	ChatEventComplete = "conversation.message.completed"
	ChatEventFailed   = "conversation.chat.failed"
	MessageTypeAnswer = "answer"
)

// Client Agent 客户端
type Client struct {
	token  string
	botId  string
	userId string
	client *http.Client
}

// Message Agent 输入消息
type Message struct {
	Text        string
	Image       []byte
	ContentType string
}

// NewClient 创建 Agent 客户端
func NewClient(token, botId, userId string) *Client {
	return &Client{
		token:  token,
		botId:  botId,
		userId: userId,
		client: http.DefaultClient,
	}
}

// Chat 调用智能体并返回模型输出
func (c *Client) Chat(ctx context.Context, message Message) (string, error) {
	if c == nil || c.client == nil {
		return "", errors.New("agent client is nil")
	}
	if c.token == "" {
		return "", errors.New("agent token is required")
	}
	if c.botId == "" {
		return "", errors.New("agent bot id is required")
	}

	fileId := ""
	if len(message.Image) > 0 {
		uploadedFileId, err := c.uploadFile(ctx, message.Image, message.ContentType)
		if err != nil {
			return "", err
		}
		fileId = uploadedFileId
	}

	return c.streamChat(ctx, strings.TrimSpace(message.Text), fileId)
}

// uploadFile 上传图像附件并返回文件 ID
func (c *Client) uploadFile(ctx context.Context, image []byte, contentType string) (string, error) {
	var body bytes.Buffer

	writer := multipart.NewWriter(&body)
	header := make(textproto.MIMEHeader)
	header.Set("Content-Disposition", `form-data; name="file"; filename="classification.png"`)
	header.Set("Content-Type", firstNotEmpty(contentType, "application/octet-stream"))

	part, err := writer.CreatePart(header)
	if err != nil {
		return "", err
	}
	if _, err := part.Write(image); err != nil {
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

	fileId := firstNotEmpty(uploadResp.Data.Id, uploadResp.Data.FileId, uploadResp.Id, uploadResp.FileId)
	if fileId == "" {
		return "", errors.New("agent file id is empty")
	}

	return fileId, nil
}

// streamChat 发起流式对话并拼接输出
func (c *Client) streamChat(ctx context.Context, text, fileId string) (string, error) {
	if fileId == "" {
		return "", errors.New("file id is required")
	}

	content, err := json.Marshal([]map[string]string{{
		"type": "text",
		"text": text,
	}, {
		"type":    "image",
		"file_id": fileId,
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

// readChatStream 读取 SSE 响应
func readChatStream(reader io.Reader) (string, error) {
	scanner := bufio.NewScanner(reader)
	scanner.Buffer(make([]byte, 1024), 2<<20)

	event := ""
	var builder strings.Builder

	for scanner.Scan() {
		line := scanner.Text()
		if strings.HasPrefix(line, SSEEventPrefix) {
			event = strings.TrimSpace(strings.TrimPrefix(line, SSEEventPrefix))
			continue
		}
		if !strings.HasPrefix(line, SSEDataPrefix) {
			continue
		}
		data := strings.TrimSpace(strings.TrimPrefix(line, SSEDataPrefix))
		if data == "" || data == SSEEventDone {
			continue
		}
		content, completedContent, err := parseChatEvent(event, data)
		if err != nil {
			return "", err
		}
		if completedContent != "" {
			builder.Reset()
			builder.WriteString(completedContent)
			continue
		}
		builder.WriteString(content)
	}
	if err := scanner.Err(); err != nil {
		return "", err
	}
	output := strings.TrimSpace(builder.String())
	if output == "" {
		return "", errors.New("agent output is empty")
	}
	return output, nil
}

// parseChatEvent 解析单条 SSE 事件数据
func parseChatEvent(event, data string) (string, string, error) {
	var textPayload string
	if err := json.Unmarshal([]byte(data), &textPayload); err == nil {
		switch event {
		case ChatEventDelta:
			return textPayload, "", nil
		case ChatEventComplete:
			return "", textPayload, nil
		case ChatEventFailed:
			return "", "", errors.New(firstNotEmpty(textPayload, "agent chat failed"))
		default:
			return "", "", nil
		}
	}

	var payload struct {
		Type      string `json:"type"`
		Content   string `json:"content"`
		LastError struct {
			Code int    `json:"code"`
			Msg  string `json:"msg"`
		} `json:"last_error"`
	}
	if err := json.Unmarshal([]byte(data), &payload); err != nil {
		return "", "", err
	}
	if event == ChatEventFailed || payload.LastError.Code != 0 {
		return "", "", errors.New(firstNotEmpty(payload.LastError.Msg, "agent chat failed"))
	}
	if payload.Type != "" && payload.Type != MessageTypeAnswer {
		return "", "", nil
	}
	switch event {
	case ChatEventDelta:
		return payload.Content, "", nil
	case ChatEventComplete:
		return "", payload.Content, nil
	default:
		return "", "", nil
	}
}

// readAgentError 读取智能体接口错误响应
func readAgentError(reader io.Reader, status string) error {
	data, err := io.ReadAll(io.LimitReader(reader, 4096))
	if err != nil {
		return err
	}
	msg := strings.TrimSpace(string(data))
	if msg == "" {
		msg = status
	}
	return errors.New(msg)
}

// firstNotEmpty 返回第一个非空字符串
func firstNotEmpty(values ...string) string {
	for _, value := range values {
		value = strings.TrimSpace(value)
		if value != "" {
			return value
		}
	}
	return ""
}
