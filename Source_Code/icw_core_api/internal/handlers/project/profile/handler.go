package profile

import (
	"icw_core_api/internal/handlers/common"
)

// Handler 基础信息 Handler
type Handler struct {
	*common.BaseHandler
}

func NewHandler(deps *common.Deps) *Handler {
	return &Handler{
		BaseHandler: common.NewBaseHandler(deps),
	}
}
