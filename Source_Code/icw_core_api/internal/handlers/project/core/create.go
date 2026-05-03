package core

import (
	"github.com/gin-gonic/gin"

	"icw_core_api/internal/dto/project"
	"icw_core_api/internal/response"
	"icw_core_api/utils"
	bizDto "icw_core_biz/pkg/dto/project"
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

	rpcReq := &bizDto.CreateProjectRequest{
		UserId: user.Id,
	}
	rpcResp := &bizDto.CreateProjectResponse{}
	if err := h.CallRPC(h.CoreBizRPCClient(), "ProjectCoreService.CreateProject", rpcReq, rpcResp); err != nil {
		response.WriteError(c, err)
		return
	}

	response.OK(c, project.NewCreateProjectResponse(rpcResp))
}
