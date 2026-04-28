package auth

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"

	"icw_core_api/internal/dto"
	"icw_core_api/internal/response"
	"icw_core_api/utils"
	bizDto "icw_core_biz/pkg/dto"
)

// Refresh .
// @router /auth/refresh [POST]
func (h *Handler) Refresh(c *gin.Context) {
	refreshToken, err := c.Cookie(RefreshCookieName)
	if err != nil || refreshToken == "" {
		response.Error(c, http.StatusUnauthorized, "UNAUTHORIZED", "refresh token is empty")
		return
	}
	req := dto.RefreshRequest{
		RefreshToken: refreshToken,
	}

	rpcReq := &bizDto.RefreshRequest{
		RefreshToken: req.RefreshToken,
	}
	rpcResp := &bizDto.RefreshResponse{}
	if err := h.rpc.Call("AuthService.Refresh", rpcReq, rpcResp); err != nil || rpcResp == nil {
		log.Printf("Call icw.core.biz AuthService.Refresh failed, req: %s, resp: %s, err: %v", utils.JSONF(rpcReq), utils.JSONF(rpcResp), err)
		response.WriteRPCError(c, err)
		h.clearRefreshCookie(c)
		return
	}

	h.setRefreshCookie(c, rpcResp.RefreshToken, rpcResp.RefreshTokenExpiresIn)

	response.OK(c, dto.NewRefreshResponse(rpcResp))
}
