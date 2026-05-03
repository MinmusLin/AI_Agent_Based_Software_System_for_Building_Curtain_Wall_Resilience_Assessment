package auth

import (
	"github.com/gin-gonic/gin"

	"icw_core_api/internal/dto"
	"icw_core_api/internal/response"
	bizDto "icw_core_biz/pkg/dto"
)

// ResetPassword 重置密码
// @router /auth/reset-password [POST]
func (h *Handler) ResetPassword(c *gin.Context) {
	var req dto.ResetPasswordRequest
	if !response.BindJSON(c, &req) {
		return
	}

	rpcReq := &bizDto.ResetPasswordRequest{
		Email:       req.Email,
		EmailCode:   req.EmailCode,
		NewPassword: req.NewPassword,
	}
	rpcResp := &bizDto.ResetPasswordResponse{}
	if err := h.CallRPC(h.CoreBizClient(), "AuthService.ResetPassword", rpcReq, rpcResp); err != nil {
		response.WriteError(c, err)
		return
	}

	response.OK(c, dto.NewResetPasswordResponse(rpcResp))
}
