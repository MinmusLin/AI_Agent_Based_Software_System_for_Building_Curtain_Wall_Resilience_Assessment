package socket

import (
	"github.com/gin-gonic/gin"

	"icw_core_api/internal/dto"
	"icw_core_api/internal/response"
	"icw_core_api/utils"
	bizDto "icw_core_biz/pkg/dto"
	bizUtils "icw_core_biz/utils"
)

// CreateSocketTicket 创建 WebSocket 连接票据
// @router /socket/ticket [POST]
func (h *Handler) CreateSocketTicket(c *gin.Context) {
	var req dto.CreateSocketTicketRequest
	if !response.BindJSON(c, &req) {
		return
	}

	// 从 Gin Context 中获取当前登录用户
	user, err := utils.GetCurrentUser(c)
	if err != nil {
		response.WriteError(c, err)
		return
	}

	// 将 Sqids 字符串解码为数字 ID
	projectId, err := bizUtils.Decode(req.ProjectId)
	if err != nil {
		response.WriteError(c, err)
		return
	}

	rpcReq := &bizDto.CreateSocketTicketRequest{
		UserId:      user.Id,
		ProjectId:   projectId,
		ProjectCode: req.ProjectId,
		SocketScope: dto.SocketScopeProjectAssets,
		RequestId:   utils.GetRequestId(c.Request.Context()),
	}
	rpcResp := &bizDto.CreateSocketTicketResponse{}
	if err := h.CoreBizCall(c.Request.Context(), "SocketService.CreateSocketTicket", rpcReq, rpcResp); err != nil {
		response.WriteError(c, err)
		return
	}

	response.OK(c, dto.NewCreateSocketTicketResponse(rpcResp))
}
