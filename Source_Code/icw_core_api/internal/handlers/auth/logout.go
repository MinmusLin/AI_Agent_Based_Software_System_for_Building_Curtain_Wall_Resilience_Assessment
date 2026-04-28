package auth

import (
	"log"

	"github.com/gin-gonic/gin"

	"icw_core_api/internal/dto"
	"icw_core_api/internal/response"
	"icw_core_api/utils"
	bizDto "icw_core_biz/pkg/dto"
)

// Logout .
// @router /auth/logout [POST]
func (h *Handler) Logout(c *gin.Context) {
	refreshToken, _ := c.Cookie(RefreshCookieName)
	req := dto.LogoutRequest{
		AccessToken:  bearerToken(c),
		RefreshToken: refreshToken,
	}

	rpcReq := &bizDto.LogoutRequest{
		AccessToken:  req.AccessToken,
		RefreshToken: req.RefreshToken,
	}
	rpcResp := &bizDto.LogoutResponse{}
	err := h.rpc.Call("AuthService.Logout", rpcReq, rpcResp)
	if err != nil {
		log.Printf("Call icw.core.biz AuthService.Logout failed, req: %s, resp: %s, err: %v", utils.JSONF(rpcReq), utils.JSONF(rpcResp), err)
	}

	h.clearRefreshCookie(c)

	response.OK(c, dto.NewLogoutResponse(rpcResp))
}
