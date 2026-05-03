package assets

import (
	"github.com/gin-gonic/gin"

	"icw_core_api/internal/dto/project"
	"icw_core_api/internal/response"
	"icw_core_api/utils"
	bizDto "icw_core_biz/pkg/dto/project"
)

// MoveProjectGroup 移动图像组
// @router /project/assets/group/move [POST]
func (h *Handler) MoveProjectGroup(c *gin.Context) {
	var req project.MoveProjectGroupRequest
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
	groupId, err := utils.Decode(req.GroupId)
	if err != nil {
		response.WriteError(c, err)
		return
	}
	var previousGroupId uint64
	if req.PreviousGroupId != "" {
		previousGroupId, err = utils.Decode(req.PreviousGroupId)
		if err != nil {
			response.WriteError(c, err)
			return
		}
	}
	var nextGroupId uint64
	if req.NextGroupId != "" {
		nextGroupId, err = utils.Decode(req.NextGroupId)
		if err != nil {
			response.WriteError(c, err)
			return
		}
	}

	rpcReq := &bizDto.MoveProjectGroupRequest{
		UserId:          user.Id,
		ProjectId:       projectId,
		GroupId:         groupId,
		PreviousGroupId: previousGroupId,
		NextGroupId:     nextGroupId,
		MoveToFirst:     req.MoveToFirst,
		MoveToLast:      req.MoveToLast,
	}
	rpcResp := &bizDto.MoveProjectGroupResponse{}
	if err := h.CoreBizCall(c.Request.Context(), "ProjectAssetsService.MoveProjectGroup", rpcReq, rpcResp); err != nil {
		response.WriteError(c, err)
		return
	}

	response.OK(c, project.NewMoveProjectGroupResponse(rpcResp))
}
