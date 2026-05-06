package auth

import (
	"github.com/gin-gonic/gin"

	"icw_common/gen/core/biz"
	"icw_common/rpc_err"
	"icw_core_api/internal/dto"
	"icw_core_api/internal/response"
	"icw_core_api/rpc/icw_core_biz/auth"
)

// Refresh 刷新 Token
// @router /auth/refresh [POST]
// 该接口在 Access Token 过期后被前端调用的，因此不要求 Access Token 有效
func (h *Handler) Refresh(c *gin.Context) {
	// Refresh Token 存在 HttpOnly Cookie 中
	refreshToken, _ := c.Cookie(RefreshCookieName)

	rpcReq := &bizpb.RefreshRequest{
		RefreshToken: refreshToken,
	}
	rpcResp := &bizpb.RefreshResponse{}
	if err := auth.Refresh(c.Request.Context(), h.CoreBizClient(), rpcReq, rpcResp); err != nil {
		response.WriteError(c, err)
		code, _, _ := rpc_err.Parse(err)
		if code == rpc_err.CodeUnauthorized {
			// 旧 Refresh Token 失效
			h.clearRefreshCookie(c)
		}
		return
	}

	// 新 Refresh Token 写入 HttpOnly Cookie，旧 Refresh Token 失效
	h.setRefreshCookie(c, rpcResp.RefreshToken, int(rpcResp.RefreshTokenExpiresIn))

	response.OK(c, dto.NewRefreshResponse(rpcResp))
}
