package profile

import (
	"log"

	"github.com/gin-gonic/gin"

	"icw_core_api/internal/dto/project"
	"icw_core_api/internal/response"
	"icw_core_api/utils"
	bizDto "icw_core_biz/pkg/dto/project"
)

// UpdateProjectProfile 更新项目基础信息
// @router /project/profile/update [POST]
func (h *Handler) UpdateProjectProfile(c *gin.Context) {
	var req project.UpdateProjectProfileRequest
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

	rpcReq := &bizDto.UpdateProjectProfileRequest{
		UserId:              user.Id,
		ProjectId:           projectId,
		Name:                req.Name,
		BuildingName:        req.BuildingName,
		BuildingLocation:    req.BuildingLocation,
		BuiltYear:           req.BuiltYear,
		BuildingDescription: req.BuildingDescription,
		KnownIssues:         req.KnownIssues,
		AssessmentGoal:      req.AssessmentGoal,
	}
	rpcResp := &bizDto.UpdateProjectProfileResponse{}
	if err := h.CoreBizClient().Call("ProjectProfileService.UpdateProjectProfile", rpcReq, rpcResp); err != nil || rpcResp == nil {
		log.Printf("[ERROR] Call icw.core.biz ProjectProfileService.UpdateProjectProfile failed, req: %s, resp: %s, err: %v", utils.JSONF(rpcReq), utils.JSONF(rpcResp), err)
		response.WriteError(c, err)
		return
	}

	response.OK(c, project.NewUpdateProjectProfileResponse(rpcResp))
}
