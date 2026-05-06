package user

import (
	"github.com/gin-gonic/gin"

	"icw_common/gen/core/biz"
	"icw_core_api/internal/response"
	"icw_core_api/rpc/icw_core_biz/user"
	"icw_core_api/utils"
)

// DeleteAvatar 删除用户自定义头像
// @router /user/avatar [DELETE]
func (h *Handler) DeleteAvatar(c *gin.Context) {
	// 从 Gin Context 中获取当前登录用户
	currentUser, err := utils.GetCurrentUser(c)
	if err != nil {
		response.WriteError(c, err)
		return
	}

	rpcReq := &bizpb.DeleteAvatarRequest{
		UserId: currentUser.Id,
		Email:  currentUser.Email,
	}
	rpcResp := &bizpb.DeleteAvatarResponse{}
	if err := user.DeleteAvatar(c.Request.Context(), h.CoreBizClient(), rpcReq, rpcResp); err != nil {
		response.WriteError(c, err)
		return
	}

	response.OK(c, utils.NewDeleteAvatarResponse(rpcResp))
}
