package auth

import (
	"github.com/gin-gonic/gin"

	"icw_common/gen/core/api"
	"icw_common/gen/core/biz"

	"icw_core_api/internal/dto"
	"icw_core_api/internal/response"
	"icw_core_api/rpc/icw_core_biz/auth"
)

// Login 登录
// @router /auth/login [POST]
func (h *Handler) Login(c *gin.Context) {
	var req apipb.LoginRequest
	if !response.BindJSON(c, &req) {
		return
	}

	rpcReq := &bizpb.LoginRequest{
		Email: req.Email,
		Scene: req.Scene,
		Code:  req.Code,
	}
	rpcResp := &bizpb.LoginResponse{}
	if err := auth.Login(c.Request.Context(), h.CoreBizClient(), rpcReq, rpcResp); err != nil {
		response.WriteError(c, err)
		return
	}

	// 新 Refresh Token 写入 HttpOnly Cookie，旧 Refresh Token 失效
	h.setRefreshCookie(c, rpcResp.RefreshToken, int(rpcResp.RefreshTokenExpiresIn))

	response.OK(c, dto.NewLoginResponse(rpcResp))
}
