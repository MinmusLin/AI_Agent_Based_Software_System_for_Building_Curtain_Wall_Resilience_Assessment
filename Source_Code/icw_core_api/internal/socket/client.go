package socket

import (
	"time"

	"github.com/gorilla/websocket"
)

const (
	// WriteWait 写消息超时时间
	WriteWait = 10 * time.Second
	// PingPeriod 服务端发送 Ping 时间间隔
	PingPeriod = 55 * time.Second
	// PongWait 客户端响应 Pong 最长等待时间
	PongWait = 60 * time.Second
	// MaxMessageSize 允许读取最大消息大小
	MaxMessageSize = 512
	// SendBufferSize 发送缓冲队列大小
	SendBufferSize = 16
)

// webSocketClient 单 WebSocket 连接
type webSocketClient struct {
	hub       *Hub
	conn      *websocket.Conn
	send      chan []byte
	projectId uint64
	scope     string
}

// newWebSocketClient 创建单 WebSocket 连接
func newWebSocketClient(hub *Hub, projectId uint64, scope string, conn *websocket.Conn) *webSocketClient {
	return &webSocketClient{
		hub:       hub,
		conn:      conn,
		send:      make(chan []byte, SendBufferSize),
		projectId: projectId,
		scope:     scope,
	}
}

// ReadPump 消费客户端消息，并维护连接生命周期
func (c *webSocketClient) ReadPump() {
	defer func() {
		c.hub.Unregister(c)
		_ = c.conn.Close()
	}()
	c.conn.SetReadLimit(MaxMessageSize)
	_ = c.conn.SetReadDeadline(time.Now().Add(PongWait))
	c.conn.SetPongHandler(func(string) error {
		_ = c.conn.SetReadDeadline(time.Now().Add(PongWait))
		return nil
	})
	for {
		if _, _, err := c.conn.ReadMessage(); err != nil {
			return
		}
	}
}

// WritePump 生产客户端消息，并维护连接心跳检测
func (c *webSocketClient) WritePump() {
	ticker := time.NewTicker(PingPeriod)
	defer func() {
		ticker.Stop()
		_ = c.conn.Close()
	}()
	for {
		select {
		case message, ok := <-c.send:
			_ = c.conn.SetWriteDeadline(time.Now().Add(WriteWait))
			if !ok {
				_ = c.conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}
			if err := c.conn.WriteMessage(websocket.TextMessage, message); err != nil {
				return
			}
		case <-ticker.C:
			_ = c.conn.SetWriteDeadline(time.Now().Add(WriteWait))
			if err := c.conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
		}
	}
}
