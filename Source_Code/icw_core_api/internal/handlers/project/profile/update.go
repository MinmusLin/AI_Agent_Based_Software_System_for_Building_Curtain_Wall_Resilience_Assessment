package profile

import (
	"github.com/gin-gonic/gin"

	"icw_common/gen/core/api"
	"icw_common/gen/core/biz"
	"icw_common/rpc/error"
	"icw_common/utils"
	"icw_core_api/internal/dto"
	"icw_core_api/internal/response"
	"icw_core_api/rpc/icw_core_biz/project_profile"
	apiUtils "icw_core_api/utils"
)

// UpdateProjectProfile 更新项目基础信息
// @router /project/profile/update [POST]
func (h *Handler) UpdateProjectProfile(c *gin.Context) {
	var req apipb.UpdateProjectProfileRequest
	if !response.BindJSON(c, &req) {
		return
	}

	// 从 Gin Context 中获取当前登录用户
	user, err := apiUtils.GetCurrentUser(c)
	if err != nil {
		response.WriteError(c, err)
		return
	}

	// 将 Sqids 字符串解码为数字 ID
	projectId, err := utils.Decode(req.ProjectId)
	if err != nil {
		response.WriteError(c, rpc_error.BadRequestDefault(err.Error()))
		return
	}

	rpcReq := &bizpb.UpdateProjectProfileRequest{
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
	rpcResp := &bizpb.UpdateProjectProfileResponse{}
	if err := project_profile.UpdateProjectProfile(c.Request.Context(), h.CoreBizClient(), rpcReq, rpcResp); err != nil {
		response.WriteError(c, err)
		return
	}

	response.OK(c, dto.NewUpdateProjectProfileResponse(rpcResp))
}
