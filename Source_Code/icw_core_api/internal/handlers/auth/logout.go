package auth

import (
	"github.com/gin-gonic/gin"

	"icw_common/gen/core/biz"
	"icw_core_api/internal/response"
	"icw_core_api/rpc/icw_core_biz/auth"
	"icw_core_api/utils"
)

// Logout 登出
// @router /auth/logout [POST]
func (h *Handler) Logout(c *gin.Context) {
	// Refresh Token 存在 HttpOnly Cookie 中，Access Token 存在 Authorization Header 中
	refreshToken, _ := c.Cookie(RefreshCookieName)

	rpcReq := &bizpb.LogoutRequest{
		AccessToken:  utils.BearerToken(c),
		RefreshToken: refreshToken,
	}
	rpcResp := &bizpb.LogoutResponse{}
	_ = auth.Logout(c.Request.Context(), h.CoreBizClient(), rpcReq, rpcResp)

	// 旧 Refresh Token 失效
	h.clearRefreshCookie(c)

	response.OK(c, utils.NewLogoutResponse(rpcResp))
}
