package auth

import (
	"net/http"
	"net/rpc"

	"github.com/gin-gonic/gin"

	"icw_core_api/configs"
)

const RefreshCookieName = "icw_refresh_token"

type Handler struct {
	cfg configs.Config
	rpc *rpc.Client
}

func NewHandler(cfg configs.Config, rpcClient *rpc.Client) *Handler {
	return &Handler{
		cfg: cfg,
		rpc: rpcClient,
	}
}

func (h *Handler) setRefreshCookie(c *gin.Context, token string, maxAge int) {
	c.SetSameSite(http.SameSiteLaxMode)
	c.SetCookie(RefreshCookieName, token, maxAge, "/auth", "", false, true)
}

func (h *Handler) clearRefreshCookie(c *gin.Context) {
	c.SetSameSite(http.SameSiteLaxMode)
	c.SetCookie(RefreshCookieName, "", -1, "/auth", "", false, true)
}
