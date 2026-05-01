package user

import (
	"log"

	"github.com/gin-gonic/gin"

	"icw_core_api/internal/dto"
	"icw_core_api/internal/response"
	"icw_core_api/utils"
	bizDto "icw_core_biz/pkg/dto"
)

// DeleteAvatar 删除用户自定义头像
// @router /user/avatar [DELETE]
func (h *Handler) DeleteAvatar(c *gin.Context) {
	user, err := utils.GetCurrentUser(c)
	if err != nil {
		response.WriteRPCError(c, err)
		return
	}

	rpcReq := &bizDto.DeleteAvatarRequest{
		UserId: user.Id,
		Email:  user.Email,
	}
	rpcResp := &bizDto.DeleteAvatarResponse{}
	if err := h.CoreBizClient().Call("UserService.DeleteAvatar", rpcReq, rpcResp); err != nil || rpcResp == nil {
		log.Printf("[ERROR] Call icw.core.biz UserService.DeleteAvatar failed, req: %s, resp: %s, err: %v", utils.JSONF(rpcReq), utils.JSONF(rpcResp), err)
		response.WriteRPCError(c, err)
		return
	}

	response.OK(c, dto.NewDeleteAvatarResponse(rpcResp))
}
