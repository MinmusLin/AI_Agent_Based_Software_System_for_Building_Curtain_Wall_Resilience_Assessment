package utils

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
)

const (
	// TicketBytes WebSocket 连接票据大小
	TicketBytes = 32
)

// SocketTicketContext WebSocket 连接票据上下文
type SocketTicketContext struct {
	UserId      uint64
	ProjectId   uint64
	ProjectCode string
	SocketScope string
	RequestId   string
	CreateAt    string
}

// NewTicket 生成 WebSocket 连接票据
func NewTicket() (string, error) {
	bytes := make([]byte, TicketBytes)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	return base64.RawURLEncoding.EncodeToString(bytes), nil
}

// TicketHash 生成 WebSocket 连接票据哈希
func TicketHash(ticket string) string {
	hash := sha256.Sum256([]byte(ticket))
	return hex.EncodeToString(hash[:])
}
