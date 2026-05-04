package socket

import (
	"icw_core_api/internal/handlers/common"
	"icw_core_api/internal/socket"
)

// Handler WebSocket Handler
type Handler struct {
	*common.BaseHandler
	hub *socket.Hub
}

func NewHandler(deps *common.Deps, hub *socket.Hub) *Handler {
	return &Handler{
		BaseHandler: common.NewBaseHandler(deps),
		hub:         hub,
	}
}
