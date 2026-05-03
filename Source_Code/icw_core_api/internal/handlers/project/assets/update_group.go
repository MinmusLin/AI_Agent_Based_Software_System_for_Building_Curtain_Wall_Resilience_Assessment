package assets

import (
	"github.com/gin-gonic/gin"

	"icw_core_api/internal/dto/project"
	"icw_core_api/internal/response"
	"icw_core_api/utils"
	bizDto "icw_core_biz/pkg/dto/project"
)

// UpdateProjectGroup 更新图像组
// @router /project/assets/group/update [POST]
func (h *Handler) UpdateProjectGroup(c *gin.Context) {
	var req project.UpdateProjectGroupRequest
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

	rpcReq := &bizDto.UpdateProjectGroupRequest{
		UserId:    user.Id,
		ProjectId: projectId,
		GroupId:   groupId,
		Name:      req.Name,
	}
	rpcResp := &bizDto.UpdateProjectGroupResponse{}
	if err := h.CallRPC(h.CoreBizRPCClient(), "ProjectAssetsService.UpdateProjectGroup", rpcReq, rpcResp); err != nil {
		response.WriteError(c, err)
		return
	}

	response.OK(c, project.NewUpdateProjectGroupResponse(rpcResp))
}
