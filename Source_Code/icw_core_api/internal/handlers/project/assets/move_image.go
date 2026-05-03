package assets

import (
	"github.com/gin-gonic/gin"

	"icw_core_api/internal/dto/project"
	"icw_core_api/internal/response"
	"icw_core_api/utils"
	bizDto "icw_core_biz/pkg/dto/project"
)

// MoveProjectImage 移动图像
// @router /project/assets/image/move [POST]
func (h *Handler) MoveProjectImage(c *gin.Context) {
	var req project.MoveProjectImageRequest
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
	targetGroupId, err := utils.Decode(req.TargetGroupId)
	if err != nil {
		response.WriteError(c, err)
		return
	}

	rpcReq := &bizDto.MoveProjectImageRequest{
		UserId:        user.Id,
		ProjectId:     projectId,
		ImageUuids:    req.ImageUuids,
		TargetGroupId: targetGroupId,
	}
	rpcResp := &bizDto.MoveProjectImageResponse{}
	if err := h.CallRPC(h.CoreBizRPCClient(), "ProjectAssetsService.MoveProjectImage", rpcReq, rpcResp); err != nil {
		response.WriteError(c, err)
		return
	}

	response.OK(c, project.NewMoveProjectImageResponse(rpcResp))
}
