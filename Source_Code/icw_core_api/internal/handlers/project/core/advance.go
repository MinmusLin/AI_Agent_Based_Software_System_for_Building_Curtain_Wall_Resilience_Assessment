package core

import (
	"github.com/gin-gonic/gin"

	"icw_core_api/internal/dto/project"
	"icw_core_api/internal/response"
	"icw_core_api/utils"
	bizDto "icw_core_biz/pkg/dto/project"
)

// AdvanceProject 项目进度流转
// @router /project/core/advance [POST]
func (h *Handler) AdvanceProject(c *gin.Context) {
	var req project.AdvanceProjectRequest
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

	rpcReq := &bizDto.AdvanceProjectRequest{
		UserId:       user.Id,
		ProjectId:    projectId,
		FromProgress: req.FromProgress,
		ToProgress:   req.ToProgress,
	}
	rpcResp := &bizDto.AdvanceProjectResponse{}
	if err := h.CallRPC(h.CoreBizRPCClient(), "ProjectCoreService.AdvanceProject", rpcReq, rpcResp); err != nil {
		response.WriteError(c, err)
		return
	}

	response.OK(c, project.NewAdvanceProjectResponse(rpcResp))
}
