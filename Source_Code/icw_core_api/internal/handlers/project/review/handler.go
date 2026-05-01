package review

import (
	"icw_core_api/internal/handlers/common"
)

// Handler 人工复核 Handler
type Handler struct {
	*common.BaseHandler
}

func NewHandler(deps *common.Deps) *Handler {
	return &Handler{
		BaseHandler: common.NewBaseHandler(deps),
	}
}
