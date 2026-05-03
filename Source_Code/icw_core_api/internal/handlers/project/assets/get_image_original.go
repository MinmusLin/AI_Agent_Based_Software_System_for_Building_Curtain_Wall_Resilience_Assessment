package assets

import (
	"github.com/gin-gonic/gin"

	"icw_core_api/internal/dto/project"
	"icw_core_api/internal/response"
	"icw_core_api/utils"
	bizDto "icw_core_biz/pkg/dto/project"
)

// GetProjectImageOriginal 获取原图
// @router /project/assets/image/original [GET]
func (h *Handler) GetProjectImageOriginal(c *gin.Context) {
	// 从 Gin Context 中获取当前登录用户
	user, err := utils.GetCurrentUser(c)
	if err != nil {
		response.WriteError(c, err)
		return
	}

	// 将 Sqids 字符串解码为数字 ID
	projectId, err := utils.Decode(c.Query("project_id"))
	if err != nil {
		response.WriteError(c, err)
		return
	}

	rpcReq := &bizDto.GetProjectImageOriginalRequest{
		UserId:    user.Id,
		ProjectId: projectId,
		ImageUuid: c.Query("image_uuid"),
	}
	rpcResp := &bizDto.GetProjectImageOriginalResponse{}
	if err := h.CoreBizCall("ProjectAssetsService.GetProjectImageOriginal", rpcReq, rpcResp); err != nil {
		response.WriteError(c, err)
		return
	}

	response.OK(c, project.NewGetProjectImageOriginalResponse(rpcResp))
}
