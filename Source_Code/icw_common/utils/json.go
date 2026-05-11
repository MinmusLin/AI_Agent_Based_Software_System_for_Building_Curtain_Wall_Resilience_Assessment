package utils

import (
	"bytes"
	"encoding/json"
	"errors"
	"strings"
)

// compactJSONString 将 JSON 字符串压缩为单行 JSON
func compactJSONString(value string) (string, error) {
	value = strings.TrimSpace(value)
	if value == "" {
		return "", errors.New("json is required")
	}
	buffer := &bytes.Buffer{}
	if err := json.Compact(buffer, []byte(value)); err != nil {
		return "", err
	}
	return buffer.String(), nil
}

// CompactJSONObjectString 将 JSON 对象字符串压缩为单行 JSON
func CompactJSONObjectString(value string) (string, error) {
	compacted, err := compactJSONString(value)
	if err != nil {
		return "", err
	}
	object := map[string]json.RawMessage{}
	if err := json.Unmarshal([]byte(compacted), &object); err != nil {
		return "", err
	}
	return compacted, nil
}
