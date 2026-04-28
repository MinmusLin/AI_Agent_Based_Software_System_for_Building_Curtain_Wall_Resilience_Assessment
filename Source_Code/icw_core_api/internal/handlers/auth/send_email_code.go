package auth

import (
	"log"

	"github.com/gin-gonic/gin"

	"icw_core_api/internal/dto"
	"icw_core_api/internal/response"
	"icw_core_api/utils"
	bizDto "icw_core_biz/pkg/dto"
)

// SendEmailCode .
// @router /auth/send-email-code [POST]
func (h *Handler) SendEmailCode(c *gin.Context) {
	var req dto.SendEmailCodeRequest
	if !bindJSON(c, &req) {
		return
	}

	rpcReq := &bizDto.SendEmailCodeRequest{
		Email: req.Email,
		Scene: req.Scene,
	}
	rpcResp := &bizDto.SendEmailCodeResponse{}
	if err := h.rpc.Call("AuthService.SendEmailCode", rpcReq, rpcResp); err != nil || rpcResp == nil {
		log.Printf("Call icw.core.biz AuthService.SendEmailCode failed, req: %s, resp: %s, err: %v", utils.JSONF(rpcReq), utils.JSONF(rpcResp), err)
		writeRPCError(c, err)
		return
	}

	response.OK(c, dto.NewSendEmailCodeResponse(rpcResp))
}
