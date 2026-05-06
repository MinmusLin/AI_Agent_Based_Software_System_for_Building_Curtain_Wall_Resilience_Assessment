package socket

import (
	"sync"

	"github.com/gorilla/websocket"
)

// Hub 管理项目维度 WebSocket 连接
type Hub struct {
	mu      sync.RWMutex
	clients map[uint64]map[string]map[*Client]struct{}
}

// NewHub 创建 WebSocket Hub
func NewHub() *Hub {
	return &Hub{
		clients: make(map[uint64]map[string]map[*Client]struct{}),
	}
}

// Register 注册 WebSocket 连接
func (h *Hub) Register(projectId uint64, scope string, conn *websocket.Conn) *Client {
	client := newClient(h, projectId, scope, conn)
	h.mu.Lock()
	defer h.mu.Unlock()
	if _, ok := h.clients[projectId]; !ok {
		h.clients[projectId] = make(map[string]map[*Client]struct{})
	}
	if _, ok := h.clients[projectId][scope]; !ok {
		h.clients[projectId][scope] = make(map[*Client]struct{})
	}
	h.clients[projectId][scope][client] = struct{}{}
	return client
}

// Unregister 注销 WebSocket 连接
func (h *Hub) Unregister(client *Client) {
	if h == nil || client == nil {
		return
	}
	h.mu.Lock()
	defer h.mu.Unlock()
	scopes, ok := h.clients[client.projectId]
	if !ok {
		return
	}
	clients, ok := scopes[client.scope]
	if !ok {
		return
	}
	if _, ok := clients[client]; !ok {
		return
	}
	delete(clients, client)
	close(client.send)
	if len(clients) == 0 {
		delete(scopes, client.scope)
	}
	if len(scopes) == 0 {
		delete(h.clients, client.projectId)
	}
}

// BroadcastProject 向指定项目的所有 WebSocket 连接广播消息
func (h *Hub) BroadcastProject(projectId uint64, projectCode, scope string, payload []byte) {
	if h == nil || len(payload) == 0 {
		return
	}
	h.mu.RLock()
	scopes := h.clients[projectId]
	clients := make([]*Client, 0, len(scopes[scope]))
	for client := range scopes[scope] {
		if client == nil {
			continue
		}
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
		WSError("[%s|%s] Broadcast failed, clients=%d sent=%d dropped=%d payload_bytes=%d", projectCode, scope, len(clients), sent, dropped, len(payload))
		return
	}
	WSInfo("[%s|%s] Broadcast succeeded, clients=%d payload_bytes=%d", projectCode, scope, len(clients), len(payload))
}
