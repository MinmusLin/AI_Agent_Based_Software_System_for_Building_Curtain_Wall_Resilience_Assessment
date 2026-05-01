package detection

import (
	"icw_core_api/internal/handlers/common"
)

// Handler 智能检测 Handler
type Handler struct {
	*common.BaseHandler
}

func NewHandler(deps *common.Deps) *Handler {
	return &Handler{
		BaseHandler: common.NewBaseHandler(deps),
	}
}
