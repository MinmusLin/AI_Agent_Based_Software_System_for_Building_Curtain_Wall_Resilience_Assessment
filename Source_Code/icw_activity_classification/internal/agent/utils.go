package agent

import (
	"bufio"
	"encoding/json"
	"errors"
	"io"
	"strings"

	"icw_common/utils"
)

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
		return "", errors.New("agent output is required")
	}

	return output, nil
}

// parseChatEvent 解析 SSE 事件数据
func parseChatEvent(event, data string) (string, string, error) {
	var textPayload string
	if err := json.Unmarshal([]byte(data), &textPayload); err == nil {
		switch event {
		case ChatEventDelta:
			return textPayload, "", nil
		case ChatEventComplete:
			return "", textPayload, nil
		case ChatEventFailed:
			return "", "", errors.New(utils.FirstNotEmpty(textPayload, "agent chat failed"))
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
		return "", "", errors.New(utils.FirstNotEmpty(payload.LastError.Msg, "agent chat failed"))
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
