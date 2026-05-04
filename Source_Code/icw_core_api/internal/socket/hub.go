package socket

import (
	"sync"

	"github.com/gorilla/websocket"
)

// Hub 管理项目维度 WebSocket 连接
type Hub struct {
	mu      sync.RWMutex
	clients map[uint64]map[*Client]struct{}
}

// NewHub 创建 WebSocket Hub
func NewHub() *Hub {
	return &Hub{
		clients: make(map[uint64]map[*Client]struct{}),
	}
}

// Register 注册 WebSocket 连接
func (h *Hub) Register(projectId uint64, conn *websocket.Conn) *Client {
	client := newClient(h, projectId, conn)
	h.mu.Lock()
	defer h.mu.Unlock()
	if _, ok := h.clients[projectId]; !ok {
		h.clients[projectId] = make(map[*Client]struct{})
	}
	h.clients[projectId][client] = struct{}{}
	return client
}

// Unregister 注销 WebSocket 连接
func (h *Hub) Unregister(client *Client) {
	if h == nil || client == nil {
		return
	}
	h.mu.Lock()
	defer h.mu.Unlock()
	clients, ok := h.clients[client.projectId]
	if !ok {
		return
	}
	if _, ok := clients[client]; !ok {
		return
	}
	delete(clients, client)
	close(client.send)
	if len(clients) == 0 {
		delete(h.clients, client.projectId)
	}
}

// BroadcastProject 向指定项目的所有 WebSocket 连接广播消息
func (h *Hub) BroadcastProject(projectId uint64, projectCode string, payload []byte) {
	if h == nil || len(payload) == 0 {
		return
	}
	h.mu.RLock()
	clients := make([]*Client, 0, len(h.clients[projectId]))
	for client := range h.clients[projectId] {
		clients = append(clients, client)
	}
	h.mu.RUnlock()

	sent := 0
	dropped := 0
	for _, client := range clients {
		select {
		case client.send <- payload:
			sent++
		default:
			dropped++
			h.Unregister(client)
		}
	}

	if dropped > 0 {
		WSError("[%s] Broadcast websocket message failed, clients=%d sent=%d dropped=%d payload_bytes=%d",
			projectCode,
			len(clients),
			sent,
			dropped,
			len(payload),
		)
		return
	}
	WSInfo("[%s] Broadcast websocket message succeeded, clients=%d payload_bytes=%d",
		projectCode,
		len(clients),
		len(payload),
	)
}
