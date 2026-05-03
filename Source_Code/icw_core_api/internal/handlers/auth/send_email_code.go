package auth

import (
	"github.com/gin-gonic/gin"

	"icw_core_api/internal/dto"
	"icw_core_api/internal/response"
	bizDto "icw_core_biz/pkg/dto"
)

// SendEmailCode 发送邮箱验证码
// @router /auth/send-email-code [POST]
func (h *Handler) SendEmailCode(c *gin.Context) {
	var req dto.SendEmailCodeRequest
	if !response.BindJSON(c, &req) {
		return
	}

	rpcReq := &bizDto.SendEmailCodeRequest{
		Email: req.Email,
		Scene: req.Scene,
	}
	rpcResp := &bizDto.SendEmailCodeResponse{}
	if err := h.CallRPC(h.CoreBizRPCClient(), "AuthService.SendEmailCode", rpcReq, rpcResp); err != nil {
		response.WriteError(c, err)
		return
	}

	response.OK(c, dto.NewSendEmailCodeResponse(rpcResp))
}
