package user

import (
	"github.com/gin-gonic/gin"

	"icw_core_api/internal/dto"
	"icw_core_api/internal/response"
	"icw_core_api/utils"
	bizDto "icw_core_biz/pkg/dto"
)

// DeleteAvatar 删除用户自定义头像
// @router /user/avatar [DELETE]
func (h *Handler) DeleteAvatar(c *gin.Context) {
	// 从 Gin Context 中获取当前登录用户
	user, err := utils.GetCurrentUser(c)
	if err != nil {
		response.WriteError(c, err)
		return
	}

	rpcReq := &bizDto.DeleteAvatarRequest{
		UserId: user.Id,
		Email:  user.Email,
	}
	rpcResp := &bizDto.DeleteAvatarResponse{}
	if err := h.CoreBizCall(c.Request.Context(), "UserService.DeleteAvatar", rpcReq, rpcResp); err != nil {
		response.WriteError(c, err)
		return
	}

	response.OK(c, dto.NewDeleteAvatarResponse(rpcResp))
}
