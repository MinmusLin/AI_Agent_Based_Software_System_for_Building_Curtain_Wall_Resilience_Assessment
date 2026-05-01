package middlewares

import (
	"errors"
	"log"
	"net/rpc"

	"github.com/gin-gonic/gin"

	"icw_core_api/consts"
	"icw_core_api/internal/response"
	"icw_core_api/utils"
	"icw_core_biz/pkg/dto"
)

// AuthRequired 登录鉴权中间件
// 只校验 Access Token，不处理 Refresh Token。Access Token 过期后，前端会收到 401 状态码
// 前端通过 /auth/refresh 使用 HttpOnly Cookie 里的 Refresh Token 换取新的 Access Token
func AuthRequired(coreBizClient *rpc.Client) gin.HandlerFunc {
	return func(c *gin.Context) {
		if coreBizClient == nil {
			response.WriteRPCError(c, errors.New("icw.core.biz client is nil"))
			c.Abort()
			return
		}

		// 调用 icw.core.biz AuthService.Me 接口完成 Token 校验并获取用户信息
		rpcReq := &dto.MeRequest{
			AccessToken: utils.BearerToken(c),
		}
		rpcResp := &dto.MeResponse{}
		if err := coreBizClient.Call("AuthService.Me", rpcReq, rpcResp); err != nil || rpcResp == nil {
			log.Printf("[ERROR] Call icw.core.biz AuthService.Me failed, req: %s, resp: %s, err: %v", utils.JSONF(rpcReq), utils.JSONF(rpcResp), err)
			response.WriteRPCError(c, err)
			c.Abort()
			return
		}

		// 将用户信息写入请求上下文，交给后续 Handler 使用
		c.Set(consts.ContextUser, rpcResp.User)

		c.Next()
	}
}
