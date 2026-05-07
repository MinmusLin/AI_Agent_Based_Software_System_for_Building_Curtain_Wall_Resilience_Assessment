package user

import (
	"github.com/gin-gonic/gin"

	"icw_common/gen/core/biz"

	"icw_core_api/internal/dto"
	"icw_core_api/internal/response"
	"icw_core_api/rpc/icw_core_biz/user"
	"icw_core_api/utils"
)

// UploadAvatar 上传用户自定义头像
// @router /user/avatar [POST]
func (h *Handler) UploadAvatar(c *gin.Context) {
	// 从 Gin Context 中获取当前登录用户
	currentUser, err := utils.GetCurrentUser(c)
	if err != nil {
		response.WriteError(c, err)
		return
	}

	rpcReq := &bizpb.UploadAvatarRequest{
		UserId:      currentUser.Id,
		Email:       currentUser.Email,
		ContentType: "image/png",
	}
	rpcResp := &bizpb.UploadAvatarResponse{}
	if err := user.UploadAvatar(c.Request.Context(), h.CoreBizClient(), rpcReq, rpcResp); err != nil {
		response.WriteError(c, err)
		return
	}

	response.OK(c, dto.NewUploadAvatarResponse(rpcResp))
}
