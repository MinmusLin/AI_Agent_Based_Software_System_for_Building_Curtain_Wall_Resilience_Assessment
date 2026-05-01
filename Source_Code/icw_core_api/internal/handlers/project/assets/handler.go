package assets

import (
	"icw_core_api/internal/handlers/common"
)

// Handler 图像资产 Handler
type Handler struct {
	*common.BaseHandler
}

func NewHandler(deps *common.Deps) *Handler {
	return &Handler{
		BaseHandler: common.NewBaseHandler(deps),
	}
}
