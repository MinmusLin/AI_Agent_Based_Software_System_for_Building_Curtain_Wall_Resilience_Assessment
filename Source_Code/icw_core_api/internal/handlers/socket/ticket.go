package socket

import (
	"strings"

	"github.com/gin-gonic/gin"

	"icw_common/consts"
	"icw_common/gen/core/api"
	"icw_common/gen/core/biz"
	"icw_common/rpc/error"
	"icw_common/utils"
	"icw_core_api/internal/dto"
	"icw_core_api/internal/response"
	"icw_core_api/rpc/icw_core_biz/socket"
	apiUtils "icw_core_api/utils"
)

// CreateSocketTicket 创建 WebSocket 连接票据
// @router /socket/ticket [POST]
func (h *Handler) CreateSocketTicket(c *gin.Context) {
	var req apipb.CreateSocketTicketRequest
	if !response.BindJSON(c, &req) {
		return
	}

	// 从 Gin Context 中获取当前登录用户
	user, err := apiUtils.GetCurrentUser(c)
	if err != nil {
		response.WriteError(c, err)
		return
	}

	// 将 Sqids 字符串解码为数字 ID
	projectId, err := utils.Decode(req.ProjectId)
	if err != nil {
		response.WriteError(c, rpc_error.BadRequestDefault(err.Error()))
		return
	}

	socketScope := strings.TrimSpace(req.SocketScope)
	if socketScope == "" {
		response.WriteError(c, rpc_error.BadRequestDefault("socket scope is required"))
		return
	}
	if socketScope != consts.SocketScopeProjectAssets && socketScope != consts.SocketScopeProjectDetection {
		response.WriteError(c, rpc_error.BadRequestDefault("socket scope is invalid"))
		return
	}

	rpcReq := &bizpb.CreateSocketTicketRequest{
		UserId:      user.Id,
		ProjectId:   projectId,
		ProjectCode: req.ProjectId,
		SocketScope: socketScope,
	}
	rpcResp := &bizpb.CreateSocketTicketResponse{}
	if err := socket.CreateSocketTicket(c.Request.Context(), h.CoreBizClient(), rpcReq, rpcResp); err != nil {
		response.WriteError(c, err)
		return
	}

	response.OK(c, dto.NewCreateSocketTicketResponse(rpcResp))
}
