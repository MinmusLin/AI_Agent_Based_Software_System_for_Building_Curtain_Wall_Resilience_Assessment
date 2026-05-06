package socket

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"

	"icw_common/consts"
	"icw_common/gen/core/biz"
	"icw_common/rpc_err"
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

// SetupAssetsWebSocket 建立图像资产 WebSocket 连接
// @router /socket/setup/assets [GET]
func (h *Handler) SetupAssetsWebSocket(c *gin.Context) {
	h.setupProjectWebSocket(c, consts.SocketScopeProjectAssets)
}

// SetupDetectionWebSocket 建立智能检测 WebSocket 连接
// @router /socket/setup/detection [GET]
func (h *Handler) SetupDetectionWebSocket(c *gin.Context) {
	h.setupProjectWebSocket(c, consts.SocketScopeProjectDetection)
}

func (h *Handler) setupProjectWebSocket(c *gin.Context, socketScope string) {
	projectCode := c.Query("project_id")
	ticket := c.Query("ticket")
	projectId, err := utils.Decode(projectCode)
	if err != nil {
		response.WriteError(c, rpc_err.BadRequestDefault(err.Error()))
		return
	}

	rpcReq := &bizpb.ValidateSocketTicketRequest{
		ProjectCode: projectCode,
		SocketScope: socketScope,
		Ticket:      ticket,
	}
	rpcResp := &bizpb.ValidateSocketTicketResponse{}
	if err := socket.ValidateSocketTicket(c.Request.Context(), h.CoreBizClient(), rpcReq, rpcResp); err != nil {
		response.WriteError(c, err)
		return
	}

	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		return
	}

	client := h.hub.Register(projectId, socketScope, conn)
	go client.WritePump()
	client.ReadPump()
}
