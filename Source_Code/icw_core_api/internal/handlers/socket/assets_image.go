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
	projectCode := c.Query("project_id")
	ticket := c.Query("ticket")
	projectId, err := utils.Decode(projectCode)
	if err != nil {
		response.WriteError(c, rpc_err.BadRequestDefault("id is invalid"))
		return
	}

	rpcReq := &bizpb.ValidateSocketTicketRequest{
		ProjectCode: projectCode,
		SocketScope: consts.SocketScopeProjectAssets,
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

	client := h.hub.Register(projectId, conn)
	go client.WritePump()
	client.ReadPump()
}
