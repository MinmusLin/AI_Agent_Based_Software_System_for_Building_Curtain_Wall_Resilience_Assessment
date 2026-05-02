package core

import (
	"log"

	"github.com/gin-gonic/gin"

	"icw_core_api/internal/dto/project"
	"icw_core_api/internal/response"
	"icw_core_api/utils"
	bizDto "icw_core_biz/pkg/dto/project"
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

	rpcReq := &bizDto.ListProjectsRequest{
		UserId: user.Id,
	}
	rpcResp := &bizDto.ListProjectsResponse{}
	if err := h.CoreBizClient().Call("ProjectCoreService.ListProjects", rpcReq, rpcResp); err != nil || rpcResp == nil {
		log.Printf("[ERROR] Call icw.core.biz ProjectCoreService.ListProjects failed, req: %s, resp: %s, err: %v", utils.JSONF(rpcReq), utils.JSONF(rpcResp), err)
		response.WriteError(c, err)
		return
	}

	response.OK(c, project.NewListProjectsResponse(rpcResp))
}
