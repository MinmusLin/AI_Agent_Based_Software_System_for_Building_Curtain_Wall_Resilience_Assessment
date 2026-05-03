package auth

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"icw_core_api/internal/handlers/common"
)

const (
	// RefreshCookieName Refresh Token 的 HttpOnly Cookie 名称
	RefreshCookieName = "icw_refresh_token"
)

// Handler 登录鉴权 Handler
type Handler struct {
	*common.BaseHandler
}

func NewHandler(deps *common.Deps) *Handler {
	if deps == nil {
		deps = &common.Deps{}
	}
	return &Handler{
		BaseHandler: common.NewBaseHandler(deps),
	}
}

// setRefreshCookie Server 签发 Refresh Token（登录成功或刷新 Token 成功）时写入 HttpOnly Cookie
// Refresh Token 使用 HttpOnly Cookie 保存，前端 JavaScript 无法读取，浏览器会在请求 /auth 路径下的接口时自动携带该 Cookie
func (h *Handler) setRefreshCookie(c *gin.Context, token string, maxAge int) {
	c.SetSameSite(http.SameSiteLaxMode)
	c.SetCookie(RefreshCookieName, token, maxAge, "/auth", "", false, true)
}

// clearRefreshCookie 登出或刷新 Token 失败时清理 HttpOnly Cookie 中的 Refresh Token
// maxAge = -1 表示让浏览器立即删除该 Cookie
func (h *Handler) clearRefreshCookie(c *gin.Context) {
	c.SetSameSite(http.SameSiteLaxMode)
	c.SetCookie(RefreshCookieName, "", -1, "/auth", "", false, true)
}
