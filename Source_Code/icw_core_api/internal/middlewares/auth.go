package middlewares

import (
	"log"
	"net/http"
	"net/rpc"

	"github.com/gin-gonic/gin"

	"icw_core_api/internal/response"
	"icw_core_api/utils"
	bizDto "icw_core_biz/pkg/dto"
)

const ContextUser = "current_user"

// AuthRequired 登录鉴权中间件
// 只校验 Access Token，不处理 Refresh Token。Access Token 过期后，前端会收到 401 状态码
// 前端通过 /auth/refresh 使用 HttpOnly Cookie 里的 Refresh Token 换取新的 Access Token
func AuthRequired(coreBizClient *rpc.Client) gin.HandlerFunc {
	return func(c *gin.Context) {
		if coreBizClient == nil {
			response.Error(c, http.StatusInternalServerError, "INTERNAL_ERROR", "icw.core.biz client is nil")
			c.Abort()
			return
		}

		// 调用 icw.core.biz AuthService.Me 接口完成 Token 校验并获取用户信息
		rpcReq := &bizDto.MeRequest{
			AccessToken: utils.BearerToken(c),
		}
		rpcResp := &bizDto.MeResponse{}
		if err := coreBizClient.Call("AuthService.Me", rpcReq, rpcResp); err != nil || rpcResp == nil {
			log.Printf("Call icw.core.biz AuthService.Me failed, req: %s, resp: %s, err: %v", utils.JSONF(rpcReq), utils.JSONF(rpcResp), err)
			response.WriteRPCError(c, err)
			c.Abort()
			return
		}

		// 将用户信息写入请求上下文，交给后续 Handler 使用
		c.Set(ContextUser, rpcResp.User)

		c.Next()
	}
}
