package core

import (
	"github.com/gin-gonic/gin"

	"icw_core_api/internal/dto/project"
	"icw_core_api/internal/response"
	"icw_core_api/utils"
	bizDto "icw_core_biz/pkg/dto/project"
)

// DeleteProject 删除项目
// @router /project/core/delete [POST]
func (h *Handler) DeleteProject(c *gin.Context) {
	var req project.DeleteProjectRequest
	if !response.BindJSON(c, &req) {
		return
	}

	// 从 Gin Context 中获取当前登录用户
	user, err := utils.GetCurrentUser(c)
	if err != nil {
		response.WriteError(c, err)
		return
	}

	// 将 Sqids 字符串解码为数字 ID
	projectId, err := utils.Decode(req.ProjectId)
	if err != nil {
		response.WriteError(c, err)
		return
	}

	rpcReq := &bizDto.DeleteProjectRequest{
		UserId:    user.Id,
		ProjectId: projectId,
	}
	rpcResp := &bizDto.DeleteProjectResponse{}
	if err := h.CallRPC(h.CoreBizRPCClient(), "ProjectCoreService.DeleteProject", rpcReq, rpcResp); err != nil {
		response.WriteError(c, err)
		return
	}

	response.OK(c, project.NewDeleteProjectResponse(rpcResp))
}
