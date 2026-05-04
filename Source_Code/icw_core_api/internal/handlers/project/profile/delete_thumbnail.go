package profile

import (
	"github.com/gin-gonic/gin"

	"icw_core_api/internal/dto/project"
	"icw_core_api/internal/response"
	"icw_core_api/utils"
	bizDto "icw_core_biz/pkg/dto/project"
	bizUtils "icw_core_biz/utils"
)

// DeleteProjectThumbnail 删除项目缩略图
// @router /project/profile/thumbnail [DELETE]
func (h *Handler) DeleteProjectThumbnail(c *gin.Context) {
	// 从 Gin Context 中获取当前登录用户
	user, err := utils.GetCurrentUser(c)
	if err != nil {
		response.WriteError(c, err)
		return
	}

	// 将 Sqids 字符串解码为数字 ID
	projectId, err := bizUtils.Decode(c.Query("project_id"))
	if err != nil {
		response.WriteError(c, err)
		return
	}

	rpcReq := &bizDto.DeleteProjectThumbnailRequest{
		UserId:    user.Id,
		ProjectId: projectId,
	}
	rpcResp := &bizDto.DeleteProjectThumbnailResponse{}
	if err := h.CoreBizCall(c.Request.Context(), "ProjectProfileService.DeleteProjectThumbnail", rpcReq, rpcResp); err != nil {
		response.WriteError(c, err)
		return
	}

	response.OK(c, project.NewDeleteProjectThumbnailResponse(rpcResp))
}
