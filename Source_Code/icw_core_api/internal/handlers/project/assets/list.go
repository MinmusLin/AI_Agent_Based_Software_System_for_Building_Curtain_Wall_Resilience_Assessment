package assets

import (
	"github.com/gin-gonic/gin"

	"icw_core_api/internal/dto/project"
	"icw_core_api/internal/response"
	"icw_core_api/utils"
	bizDto "icw_core_biz/pkg/dto/project"
)

// GetProjectAssets 获取项目图像列表
// @router /project/assets/list [GET]
func (h *Handler) GetProjectAssets(c *gin.Context) {
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

	rpcReq := &bizDto.GetProjectAssetsRequest{
		UserId:    user.Id,
		ProjectId: projectId,
	}
	rpcResp := &bizDto.GetProjectAssetsResponse{}
	if err := h.CallRPC(h.CoreBizClient(), "ProjectAssetsService.GetProjectAssets", rpcReq, rpcResp); err != nil {
		response.WriteError(c, err)
		return
	}

	response.OK(c, project.NewGetProjectAssetsResponse(rpcResp))
}
