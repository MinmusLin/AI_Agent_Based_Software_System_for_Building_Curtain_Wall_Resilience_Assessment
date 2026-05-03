package user

import (
	"github.com/gin-gonic/gin"

	"icw_core_api/internal/dto"
	"icw_core_api/internal/response"
	"icw_core_api/utils"
	bizDto "icw_core_biz/pkg/dto"
)

// GetAvatar 获取用户头像
// @router /user/avatar [GET]
func (h *Handler) GetAvatar(c *gin.Context) {
	// GetCurrentUser 从 Gin Context 中获取当前登录用户
	user, err := utils.GetCurrentUser(c)
	if err != nil {
		response.WriteError(c, err)
		return
	}

	rpcReq := &bizDto.GetAvatarRequest{
		UserId: user.Id,
		Email:  user.Email,
	}
	rpcResp := &bizDto.GetAvatarResponse{}
	if err := h.CallRPC(h.CoreBizRPCClient(), "UserService.GetAvatar", rpcReq, rpcResp); err != nil {
		response.WriteError(c, err)
		return
	}

	response.OK(c, dto.NewGetAvatarResponse(rpcResp))
}
