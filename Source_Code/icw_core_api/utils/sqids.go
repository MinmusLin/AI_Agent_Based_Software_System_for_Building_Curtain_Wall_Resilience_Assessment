package utils

import (
	"errors"
	"log"
	"strings"

	"github.com/sqids/sqids-go"
)

// MinProjectIdLength 生成 Sqids ID 的最小长度
const MinProjectIdLength uint8 = 10

var projectIdCodec = mustNewProjectIdCodec()

// mustNewProjectIdCodec 创建 Sqids 编解码器实例
func mustNewProjectIdCodec() *sqids.Sqids {
	codec, err := sqids.New(sqids.Options{
		MinLength: MinProjectIdLength,
	})
	if err != nil {
		log.Fatalf("Failed to initialize sqids codec: %v", err)
	}
	return codec
}

// Encode 将数字 ID 编码为 Sqids 字符串
func Encode(id uint64) string {
	if id == 0 {
		return ""
	}
	encoded, err := projectIdCodec.Encode([]uint64{id})
	if err != nil {
		return ""
	}
	return encoded
}

// Decode 将 Sqids 字符串解码为数字 ID
func Decode(id string) (uint64, error) {
	id = strings.TrimSpace(id)
	if id == "" {
		return 0, errors.New("id is empty")
	}
	decoded := projectIdCodec.Decode(id)
	if len(decoded) != 1 || decoded[0] == 0 {
		return 0, errors.New("id is invalid")
	}
	canonicalId, err := projectIdCodec.Encode(decoded)
	if err != nil || canonicalId != id {
		return 0, errors.New("id is invalid")
	}
	return decoded[0], nil
}
