package user

import (
	"github.com/gin-gonic/gin"

	"icw_core_api/internal/dto"
	"icw_core_api/internal/response"
	"icw_core_api/utils"
	bizDto "icw_core_biz/pkg/dto"
)

// UploadAvatar 上传用户自定义头像
// @router /user/avatar [POST]
func (h *Handler) UploadAvatar(c *gin.Context) {
	// GetCurrentUser 从 Gin Context 中获取当前登录用户
	user, err := utils.GetCurrentUser(c)
	if err != nil {
		response.WriteError(c, err)
		return
	}

	rpcReq := &bizDto.UploadAvatarRequest{
		UserId:      user.Id,
		Email:       user.Email,
		ContentType: "image/png",
	}
	rpcResp := &bizDto.UploadAvatarResponse{}
	if err := h.CallRPC(h.CoreBizRPCClient(), "UserService.UploadAvatar", rpcReq, rpcResp); err != nil {
		response.WriteError(c, err)
		return
	}

	response.OK(c, dto.NewUploadAvatarResponse(rpcResp))
}
