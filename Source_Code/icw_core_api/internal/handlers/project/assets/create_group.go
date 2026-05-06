package assets

import (
	"github.com/gin-gonic/gin"

	"icw_common/gen/core/api"
	"icw_common/gen/core/biz"
	"icw_common/rpc_err"
	"icw_common/utils"
	"icw_core_api/internal/dto"
	"icw_core_api/internal/response"
	"icw_core_api/rpc/icw_core_biz/project_assets"
	apiUtils "icw_core_api/utils"
)

// CreateProjectGroup 创建图像组
// @router /project/assets/group/create [POST]
func (h *Handler) CreateProjectGroup(c *gin.Context) {
	var req apipb.CreateProjectGroupRequest
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
		response.WriteError(c, rpc_err.BadRequestDefault(err.Error()))
		return
	}

	rpcReq := &bizpb.CreateProjectGroupRequest{
		UserId:    user.Id,
		ProjectId: projectId,
	}
	rpcResp := &bizpb.CreateProjectGroupResponse{}
	if err := project_assets.CreateProjectGroup(c.Request.Context(), h.CoreBizClient(), rpcReq, rpcResp); err != nil {
		response.WriteError(c, err)
		return
	}

	response.OK(c, dto.NewCreateProjectGroupResponse(rpcResp))
}
