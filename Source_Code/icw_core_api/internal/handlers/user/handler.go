package user

import (
	"icw_core_api/internal/handlers/common"
)

// Handler 用户业务 Handler
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
