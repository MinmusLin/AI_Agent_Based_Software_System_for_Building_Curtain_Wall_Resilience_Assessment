package auth

import (
	"log"

	"github.com/gin-gonic/gin"

	"icw_core_api/internal/dto"
	"icw_core_api/internal/response"
	"icw_core_api/utils"
	bizDto "icw_core_biz/pkg/dto"
)

// Login 登录
// @router /auth/login [POST]
func (h *Handler) Login(c *gin.Context) {
	var req dto.LoginRequest
	if !response.BindJSON(c, &req) {
		return
	}

	rpcReq := &bizDto.LoginRequest{
		Email: req.Email,
		Scene: req.Scene,
		Code:  req.Code,
	}
	rpcResp := &bizDto.LoginResponse{}
	if err := h.CoreBizClient().Call("AuthService.Login", rpcReq, rpcResp); err != nil || rpcResp == nil {
		log.Printf("[ERROR] Call icw.core.biz AuthService.Login failed, req: %s, resp: %s, err: %v", utils.JSONF(rpcReq), utils.JSONF(rpcResp), err)
		response.WriteRPCError(c, err)
		return
	}

	// 新 Refresh Token 写入 HttpOnly Cookie，旧 Refresh Token 失效
	h.setRefreshCookie(c, rpcResp.RefreshToken, rpcResp.RefreshTokenExpiresIn)

	response.OK(c, dto.NewLoginResponse(rpcResp))
}
