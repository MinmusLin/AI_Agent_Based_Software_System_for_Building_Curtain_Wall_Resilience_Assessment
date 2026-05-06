package utils

import (
	"errors"
	"strings"

	"github.com/sqids/sqids-go"
)

const (
	// MinSqidsIdLength 生成 Sqids ID 的最小长度
	MinSqidsIdLength uint8 = 10
)

var (
	// sqidsCodec 全局 Sqids 编解码器实例
	sqidsCodec = newSqidsCodec()
)

// newSqidsCodec 创建全局 Sqids 编解码器实例
func newSqidsCodec() *sqids.Sqids {
	codec, err := sqids.New(sqids.Options{
		MinLength: MinSqidsIdLength,
	})
	if err != nil {
		return nil
	}
	return codec
}

// Encode 将数字 ID 编码为 Sqids 字符串
func Encode(id uint64) string {
	if sqidsCodec == nil {
		return ""
	}
	if id == 0 {
		return ""
	}
	encoded, err := sqidsCodec.Encode([]uint64{id})
	if err != nil {
		return ""
	}
	return encoded
}

// Decode 将 Sqids 字符串解码为数字 ID
func Decode(id string) (uint64, error) {
	if sqidsCodec == nil {
		return 0, errors.New("sqids codec is nil")
	}
	id = strings.TrimSpace(id)
	if id == "" {
		return 0, errors.New("id is empty")
	}
	decoded := sqidsCodec.Decode(id)
	if len(decoded) != 1 || decoded[0] == 0 {
		return 0, errors.New("id is invalid")
	}
	canonicalId, err := sqidsCodec.Encode(decoded)
	if err != nil || canonicalId != id {
		return 0, errors.New("id is invalid")
	}
	return decoded[0], nil
}
