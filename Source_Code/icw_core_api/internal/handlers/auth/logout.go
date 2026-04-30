package auth

import (
	"log"

	"github.com/gin-gonic/gin"

	"icw_core_api/internal/dto"
	"icw_core_api/internal/response"
	"icw_core_api/utils"
	bizDto "icw_core_biz/pkg/dto"
)

// Logout 登出
// @router /auth/logout [POST]
func (h *Handler) Logout(c *gin.Context) {
	// Refresh Token 存在 HttpOnly Cookie 中，Access Token 存在 Authorization Header 中
	refreshToken, _ := c.Cookie(RefreshCookieName)

	rpcReq := &bizDto.LogoutRequest{
		AccessToken:  utils.BearerToken(c),
		RefreshToken: refreshToken,
	}
	rpcResp := &bizDto.LogoutResponse{}
	err := h.rpc.Call("AuthService.Logout", rpcReq, rpcResp)
	if err != nil {
		log.Printf("[ERROR] Call icw.core.biz AuthService.Logout failed, req: %s, resp: %s, err: %v", utils.JSONF(rpcReq), utils.JSONF(rpcResp), err)
	}

	// 旧 Refresh Token 失效
	h.clearRefreshCookie(c)

	response.OK(c, dto.NewLogoutResponse(rpcResp))
}
