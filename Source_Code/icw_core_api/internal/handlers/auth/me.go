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

// Me .
// @router /auth/me [GET]
func (h *Handler) Me(c *gin.Context) {
	req := dto.MeRequest{
		AccessToken: bearerToken(c),
	}
	if req.AccessToken == "" {
		response.Error(c, http.StatusUnauthorized, "UNAUTHORIZED", "access token is empty")
		return
	}

	rpcReq := &bizDto.MeRequest{
		AccessToken: req.AccessToken,
	}
	rpcResp := &bizDto.MeResponse{}
	if err := h.rpc.Call("AuthService.Me", rpcReq, rpcResp); err != nil || rpcResp == nil {
		log.Printf("Call icw.core.biz AuthService.Me failed, req: %s, resp: %s, err: %v", utils.JSONF(rpcReq), utils.JSONF(rpcResp), err)
		writeRPCError(c, err)
		return
	}

	response.OK(c, dto.NewMeResponse(rpcResp))
}
