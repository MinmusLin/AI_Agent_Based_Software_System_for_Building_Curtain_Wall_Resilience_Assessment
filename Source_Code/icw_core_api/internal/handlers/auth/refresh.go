package auth

import (
	"github.com/gin-gonic/gin"

	"icw_core_api/internal/dto"
	"icw_core_api/internal/response"
	bizDto "icw_core_biz/pkg/dto"
	"icw_core_biz/pkg/rpc_err"
)

// Refresh 刷新 Token
// @router /auth/refresh [POST]
// 该接口在 Access Token 过期后被前端调用的，因此不要求 Access Token 有效
func (h *Handler) Refresh(c *gin.Context) {
	// Refresh Token 存在 HttpOnly Cookie 中
	refreshToken, _ := c.Cookie(RefreshCookieName)

	rpcReq := &bizDto.RefreshRequest{
		RefreshToken: refreshToken,
	}
	rpcResp := &bizDto.RefreshResponse{}
	if err := h.CallRPC(h.CoreBizRPCClient(), "AuthService.Refresh", rpcReq, rpcResp); err != nil {
		response.WriteError(c, err)
		code, _, _ := rpc_err.Parse(err)
		if code == rpc_err.CodeUnauthorized {
			// 旧 Refresh Token 失效
			h.clearRefreshCookie(c)
		}
		return
	}

	// 新 Refresh Token 写入 HttpOnly Cookie，旧 Refresh Token 失效
	h.setRefreshCookie(c, rpcResp.RefreshToken, rpcResp.RefreshTokenExpiresIn)

	response.OK(c, dto.NewRefreshResponse(rpcResp))
}
