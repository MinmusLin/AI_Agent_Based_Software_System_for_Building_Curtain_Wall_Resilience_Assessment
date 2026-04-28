package auth

import (
	"log"

	"github.com/gin-gonic/gin"

	"icw_core_api/internal/dto"
	"icw_core_api/internal/response"
	"icw_core_api/utils"
	bizDto "icw_core_biz/pkg/dto"
)

// ResetPassword .
// @router /auth/reset-password [POST]
func (h *Handler) ResetPassword(c *gin.Context) {
	var req dto.ResetPasswordRequest
	if !bindJSON(c, &req) {
		return
	}

	rpcReq := &bizDto.ResetPasswordRequest{
		Email:       req.Email,
		EmailCode:   req.EmailCode,
		NewPassword: req.NewPassword,
	}
	rpcResp := &bizDto.ResetPasswordResponse{}
	if err := h.rpc.Call("AuthService.ResetPassword", rpcReq, rpcResp); err != nil || rpcResp == nil {
		log.Printf("Call icw.core.biz AuthService.ResetPassword failed, req: %s, resp: %s, err: %v", utils.JSONF(rpcReq), utils.JSONF(rpcResp), err)
		writeRPCError(c, err)
		return
	}

	response.OK(c, dto.NewResetPasswordResponse(rpcResp))
}
