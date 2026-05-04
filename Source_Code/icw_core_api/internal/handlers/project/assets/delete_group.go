package assets

import (
	"github.com/gin-gonic/gin"

	"icw_core_api/internal/dto/project"
	"icw_core_api/internal/response"
	"icw_core_api/utils"
	bizDto "icw_core_biz/pkg/dto/project"
	bizUtils "icw_core_biz/utils"
)

// DeleteProjectGroup 删除图像组
// @router /project/assets/group/delete [POST]
func (h *Handler) DeleteProjectGroup(c *gin.Context) {
	var req project.DeleteProjectGroupRequest
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
	projectId, err := bizUtils.Decode(req.ProjectId)
	if err != nil {
		response.WriteError(c, err)
		return
	}
	groupId, err := bizUtils.Decode(req.GroupId)
	if err != nil {
		response.WriteError(c, err)
		return
	}

	rpcReq := &bizDto.DeleteProjectGroupRequest{
		UserId:    user.Id,
		ProjectId: projectId,
		GroupId:   groupId,
	}
	rpcResp := &bizDto.DeleteProjectGroupResponse{}
	if err := h.CoreBizCall(c.Request.Context(), "ProjectAssetsService.DeleteProjectGroup", rpcReq, rpcResp); err != nil {
		response.WriteError(c, err)
		return
	}

	response.OK(c, project.NewDeleteProjectGroupResponse(rpcResp))
}
