package socket

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"

	"icw_core_api/internal/dto"
	"icw_core_api/internal/response"
	bizDto "icw_core_biz/pkg/dto"
	"icw_core_biz/utils"
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
		response.WriteError(c, err)
		return
	}

	rpcReq := &bizDto.ValidateSocketTicketRequest{
		ProjectCode: projectCode,
		SocketScope: dto.SocketScopeProjectAssets,
		Ticket:      ticket,
	}
	rpcResp := &bizDto.ValidateSocketTicketResponse{}
	if err := h.CoreBizCall(c.Request.Context(), "SocketService.ValidateSocketTicket", rpcReq, rpcResp); err != nil {
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
