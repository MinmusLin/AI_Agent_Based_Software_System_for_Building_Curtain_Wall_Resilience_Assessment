package socket

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"

	"icw_common/consts"
	"icw_common/gen/core/api"
	"icw_common/gen/core/biz"
	"icw_common/rpc/error"
	"icw_common/utils"

	"icw_core_api/internal/response"
	"icw_core_api/rpc/icw_core_biz/socket"
)

// upgrader 从 HTTP 协议升级至 WebSocket 协议
var upgrader = websocket.Upgrader{
	CheckOrigin: func(_ *http.Request) bool {
		return true
	},
}

// SetupWebSocket 建立 WebSocket 连接
// @router /socket/setup [GET]
func (h *Handler) SetupWebSocket(c *gin.Context) {
	req := &apipb.SetupWebSocketRequest{
		ProjectId: c.Query("project_id"),
		Scope:     strings.TrimSpace(c.Query("scope")),
		Ticket:    strings.TrimSpace(c.Query("ticket")),
	}
	projectId, err := utils.Decode(req.ProjectId)
	if err != nil {
		response.Error(c, rpc_error.BadRequestDefault(err.Error()))
		return
	}

	// 校验是否是有效的 WebSocket 连接范围
	if !isValidSocketScope(req.Scope) {
		response.Error(c, rpc_error.BadRequestDefault("socket scope is invalid"))
		return
	}

	rpcReq := &bizpb.ValidateSocketTicketRequest{
		ProjectCode: req.ProjectId,
		Scope:       req.Scope,
		Ticket:      req.Ticket,
	}
	rpcResp := &bizpb.ValidateSocketTicketResponse{}
	if err := socket.ValidateSocketTicket(c.Request.Context(), h.CoreBizClient(), rpcReq, rpcResp); err != nil {
		response.Error(c, err)
		return
	}

	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		return
	}

	client := h.hub.Register(projectId, req.Scope, conn)
	go client.WritePump()
	client.ReadPump()
}

// isValidSocketScope 校验是否是有效的 WebSocket 连接范围
func isValidSocketScope(scope string) bool {
	return scope == consts.SocketScopeProjectAssets || scope == consts.SocketScopeProjectDetection || scope == consts.SocketScopeProjectReport
}
