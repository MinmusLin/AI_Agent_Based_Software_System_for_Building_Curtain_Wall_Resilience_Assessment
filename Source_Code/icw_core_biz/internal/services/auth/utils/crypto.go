package utils

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"strings"

	"github.com/google/uuid"
)

// NewEmailCode 生成 6 位数字邮箱验证码
func NewEmailCode() (string, error) {
	randomBytes := make([]byte, 4)
	if _, err := rand.Read(randomBytes); err != nil {
		return "", errors.New("failed to generate random bytes")
	}
	value := uint32(randomBytes[0])<<24 | uint32(randomBytes[1])<<16 | uint32(randomBytes[2])<<8 | uint32(randomBytes[3])
	return fmt.Sprintf("%06d", value%1000000), nil
}

// HashEmailCode 对邮箱验证码和 EMAIL_CODE_SECRET 做 SHA-256 哈希
func HashEmailCode(code, secret string) string {
	sum := sha256.Sum256([]byte(code + "." + secret))
	return hex.EncodeToString(sum[:])
}

// NewRefreshToken 生成 Refresh Token 明文和对应 Token Id
// 明文格式为 <token_id>.<random_secret>，Token Id 便于数据库定位记录，Random Secret 用于证明客户端持有真实 Token
func NewRefreshToken() (string, string, error) {
	tokenId := uuid.NewString()
	randomBytes := make([]byte, 32)
	if _, err := rand.Read(randomBytes); err != nil {
		return "", "", errors.New("failed to generate random bytes")
	}
	return tokenId + "." + hex.EncodeToString(randomBytes), tokenId, nil
}

// HashRefreshToken 对 Refresh Token 明文做 SHA-256 哈希
func HashRefreshToken(token string) string {
	sum := sha256.Sum256([]byte(token))
	return hex.EncodeToString(sum[:])
}

// ParseRefreshTokenId 从 Refresh Token 明文中解析 Token Id
func ParseRefreshTokenId(token string) string {
	tokenId, _, ok := strings.Cut(token, ".")
	if !ok || tokenId == "" {
		return ""
	}
	return tokenId
}
