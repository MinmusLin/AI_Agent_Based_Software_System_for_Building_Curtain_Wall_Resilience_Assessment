package assets

import (
	"github.com/gin-gonic/gin"

	"icw_core_api/internal/dto/project"
	"icw_core_api/internal/response"
	"icw_core_api/utils"
	bizDto "icw_core_biz/pkg/dto/project"
	bizUtils "icw_core_biz/utils"
)

// DeleteProjectImage 删除图像
// @router /project/assets/image/delete [POST]
func (h *Handler) DeleteProjectImage(c *gin.Context) {
	var req project.DeleteProjectImageRequest
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

	rpcReq := &bizDto.DeleteProjectImageRequest{
		UserId:     user.Id,
		ProjectId:  projectId,
		ImageUuids: req.ImageUuids,
	}
	rpcResp := &bizDto.DeleteProjectImageResponse{}
	if err := h.CoreBizCall(c.Request.Context(), "ProjectAssetsService.DeleteProjectImage", rpcReq, rpcResp); err != nil {
		response.WriteError(c, err)
		return
	}

	response.OK(c, project.NewDeleteProjectImageResponse(rpcResp))
}
