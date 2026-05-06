package user

import (
	"github.com/gin-gonic/gin"

	"icw_common/gen/core/biz"
	"icw_core_api/internal/response"
	"icw_core_api/rpc/icw_core_biz/user"
	"icw_core_api/utils"
)

// GetAvatar 获取用户头像
// @router /user/avatar [GET]
func (h *Handler) GetAvatar(c *gin.Context) {
	// 从 Gin Context 中获取当前登录用户
	currentUser, err := utils.GetCurrentUser(c)
	if err != nil {
		response.WriteError(c, err)
		return
	}

	rpcReq := &bizpb.GetAvatarRequest{
		UserId: currentUser.Id,
		Email:  currentUser.Email,
	}
	rpcResp := &bizpb.GetAvatarResponse{}
	if err := user.GetAvatar(c.Request.Context(), h.CoreBizClient(), rpcReq, rpcResp); err != nil {
		response.WriteError(c, err)
		return
	}

	response.OK(c, utils.NewGetAvatarResponse(rpcResp))
}
