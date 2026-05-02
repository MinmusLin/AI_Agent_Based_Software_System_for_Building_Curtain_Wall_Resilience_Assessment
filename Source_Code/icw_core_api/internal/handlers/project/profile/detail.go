package profile

import (
	"log"

	"github.com/gin-gonic/gin"

	"icw_core_api/internal/dto/project"
	"icw_core_api/internal/response"
	"icw_core_api/utils"
	bizDto "icw_core_biz/pkg/dto/project"
)

// GetProjectProfile 获取项目基础信息
// @router /project/profile/detail [GET]
func (h *Handler) GetProjectProfile(c *gin.Context) {
	// 从 Gin Context 中获取当前登录用户
	user, err := utils.GetCurrentUser(c)
	if err != nil {
		response.WriteRPCError(c, err)
		return
	}

	// 将 Sqids 字符串解码为数字 ID
	projectId, err := utils.Decode(c.Query("project_id"))
	if err != nil {
		response.WriteRPCError(c, err)
		return
	}

	rpcReq := &bizDto.GetProjectProfileRequest{
		UserId:    user.Id,
		ProjectId: projectId,
	}
	rpcResp := &bizDto.GetProjectProfileResponse{}
	if err := h.CoreBizClient().Call("ProjectProfileService.GetProjectProfile", rpcReq, rpcResp); err != nil || rpcResp == nil {
		log.Printf("[ERROR] Call icw.core.biz ProjectProfileService.GetProjectProfile failed, req: %s, resp: %s, err: %v", utils.JSONF(rpcReq), utils.JSONF(rpcResp), err)
		response.WriteRPCError(c, err)
		return
	}

	response.OK(c, project.NewGetProjectProfileResponse(rpcResp))
}
