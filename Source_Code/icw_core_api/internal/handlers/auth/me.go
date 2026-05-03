package auth

import (
	"github.com/gin-gonic/gin"

	"icw_core_api/internal/dto"
	"icw_core_api/internal/response"
	"icw_core_api/utils"
	bizDto "icw_core_biz/pkg/dto"
)

// Me 获取用户信息
// @router /auth/me [GET]
func (h *Handler) Me(c *gin.Context) {
	// 正常路由会先经过 AuthRequired 登录鉴权中间件，预期已经调用过 icw.core.biz AuthService.Me 接口
	// 优先从 Gin Context 中获取用户信息，避免 RPC 重复调用
	if user, err := utils.GetCurrentUser(c); err == nil {
		response.OK(c, dto.NewMeResponse(&bizDto.MeResponse{
			User: user,
		}))
		return
	}

	// 兼容未经过 AuthRequired 登录鉴权中间件的调用场景
	rpcReq := &bizDto.MeRequest{
		AccessToken: utils.BearerToken(c),
	}
	rpcResp := &bizDto.MeResponse{}
	if err := h.CoreBizCall(c.Request.Context(), "AuthService.Me", rpcReq, rpcResp); err != nil {
		response.WriteError(c, err)
		return
	}

	response.OK(c, dto.NewMeResponse(rpcResp))
}
