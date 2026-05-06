package core

import (
	"github.com/gin-gonic/gin"

	"icw_common/gen/core/biz"
	"icw_core_api/internal/response"
	"icw_core_api/rpc/icw_core_biz/project_core"
	"icw_core_api/utils"
)

// CreateProject 创建项目
// @router /project/core/create [POST]
func (h *Handler) CreateProject(c *gin.Context) {
	// 从 Gin Context 中获取当前登录用户
	user, err := utils.GetCurrentUser(c)
	if err != nil {
		response.WriteError(c, err)
		return
	}

	rpcReq := &bizpb.CreateProjectRequest{
		UserId: user.Id,
	}
	rpcResp := &bizpb.CreateProjectResponse{}
	if err := project_core.CreateProject(c.Request.Context(), h.CoreBizClient(), rpcReq, rpcResp); err != nil {
		response.WriteError(c, err)
		return
	}

	response.OK(c, utils.NewCreateProjectResponse(rpcResp))
}
