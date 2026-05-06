package core

import (
	"github.com/gin-gonic/gin"

	"icw_common/gen/core/biz"
	"icw_core_api/internal/dto"
	"icw_core_api/internal/response"
	"icw_core_api/rpc/icw_core_biz/project_core"
	"icw_core_api/utils"
)

// ListProjects 获取项目列表
// @router /project/core/list [GET]
func (h *Handler) ListProjects(c *gin.Context) {
	// 从 Gin Context 中获取当前登录用户
	user, err := utils.GetCurrentUser(c)
	if err != nil {
		response.WriteError(c, err)
		return
	}

	rpcReq := &bizpb.ListProjectsRequest{
		UserId: user.Id,
	}
	rpcResp := &bizpb.ListProjectsResponse{}
	if err := project_core.ListProjects(c.Request.Context(), h.CoreBizClient(), rpcReq, rpcResp); err != nil {
		response.WriteError(c, err)
		return
	}

	response.OK(c, dto.NewListProjectsResponse(rpcResp))
}
