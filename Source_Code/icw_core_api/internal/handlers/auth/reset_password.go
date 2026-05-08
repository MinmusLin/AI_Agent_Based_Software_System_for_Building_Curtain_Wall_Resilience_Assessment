package auth

import (
	"github.com/gin-gonic/gin"

	"icw_common/gen/core/api"
	"icw_common/gen/core/biz"

	"icw_core_api/internal/dto"
	"icw_core_api/internal/response"
	"icw_core_api/rpc/icw_core_biz/auth"
)

// ResetPassword 重置密码
// @router /auth/reset-password [POST]
func (h *Handler) ResetPassword(c *gin.Context) {
	var req apipb.ResetPasswordRequest
	if !response.BindJSON(c, &req) {
		return
	}

	rpcReq := &bizpb.ResetPasswordRequest{
		Email:       req.Email,
		EmailCode:   req.EmailCode,
		NewPassword: req.NewPassword,
	}
	rpcResp := &bizpb.ResetPasswordResponse{}
	if err := auth.ResetPassword(c.Request.Context(), h.CoreBizClient(), rpcReq, rpcResp); err != nil {
		response.Error(c, err)
		return
	}

	response.OK(c, dto.NewResetPasswordResponse(rpcResp))
}
