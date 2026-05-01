package core

import (
	"icw_core_api/internal/handlers/common"
)

// Handler 项目核心 Handler
type Handler struct {
	*common.BaseHandler
}

func NewHandler(deps *common.Deps) *Handler {
	return &Handler{
		BaseHandler: common.NewBaseHandler(deps),
	}
}
