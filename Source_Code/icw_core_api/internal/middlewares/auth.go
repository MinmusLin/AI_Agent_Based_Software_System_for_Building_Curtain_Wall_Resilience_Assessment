package middlewares

import (
	"github.com/gin-gonic/gin"

	"icw_core_api/consts"
	"icw_core_api/internal/handlers/common"
	"icw_core_api/internal/response"
	"icw_core_api/utils"
	"icw_core_biz/pkg/dto"
)

// AuthRequired 登录鉴权中间件
// 只校验 Access Token，不处理 Refresh Token。Access Token 过期后，前端会收到 401 状态码
// 前端通过 /auth/refresh 使用 HttpOnly Cookie 里的 Refresh Token 换取新的 Access Token
func AuthRequired(coreBizClient *common.RPCClient) gin.HandlerFunc {
	return func(c *gin.Context) {
		// 调用 icw.core.biz AuthService.Me 接口完成 Token 校验并获取用户信息
		rpcReq := &dto.MeRequest{
			AccessToken: utils.BearerToken(c),
		}
		rpcResp := &dto.MeResponse{}
		if err := common.CallRPC(c.Request.Context(), coreBizClient, "AuthService.Me", rpcReq, rpcResp); err != nil {
			response.WriteError(c, err)
			c.Abort()
			return
		}

		// 将用户信息写入请求上下文，交给后续 Handler 使用
		c.Set(consts.ContextUser, rpcResp.User)

		c.Next()
	}
}
