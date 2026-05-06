package auth

import (
	"github.com/gin-gonic/gin"

	"icw_common/gen/core/api"
	"icw_common/gen/core/biz"
	"icw_core_api/internal/dto"
	"icw_core_api/internal/response"
	"icw_core_api/rpc/icw_core_biz/auth"
)

// SendEmailCode 发送邮箱验证码
// @router /auth/send-email-code [POST]
func (h *Handler) SendEmailCode(c *gin.Context) {
	var req apipb.SendEmailCodeRequest
	if !response.BindJSON(c, &req) {
		return
	}

	rpcReq := &bizpb.SendEmailCodeRequest{
		Email: req.Email,
		Scene: req.Scene,
	}
	rpcResp := &bizpb.SendEmailCodeResponse{}
	if err := auth.SendEmailCode(c.Request.Context(), h.CoreBizClient(), rpcReq, rpcResp); err != nil {
		response.WriteError(c, err)
		return
	}

	response.OK(c, dto.NewSendEmailCodeResponse(rpcResp))
}
