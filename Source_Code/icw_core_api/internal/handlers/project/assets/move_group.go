package assets

import (
	"github.com/gin-gonic/gin"

	"icw_common/gen/core/api"
	"icw_common/gen/core/biz"
	"icw_common/rpc/error"
	"icw_common/utils"
	"icw_core_api/internal/dto"
	"icw_core_api/internal/response"
	"icw_core_api/rpc/icw_core_biz/project_assets"
	apiUtils "icw_core_api/utils"
)

// MoveProjectGroup 移动图像组
// @router /project/assets/group/move [POST]
func (h *Handler) MoveProjectGroup(c *gin.Context) {
	var req apipb.MoveProjectGroupRequest
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
	groupId, err := utils.Decode(req.GroupId)
	if err != nil {
		response.WriteError(c, rpc_error.BadRequestDefault(err.Error()))
		return
	}
	var previousGroupId uint64
	if req.PreviousGroupId != "" {
		previousGroupId, err = utils.Decode(req.PreviousGroupId)
		if err != nil {
			response.WriteError(c, rpc_error.BadRequestDefault(err.Error()))
			return
		}
	}
	var nextGroupId uint64
	if req.NextGroupId != "" {
		nextGroupId, err = utils.Decode(req.NextGroupId)
		if err != nil {
			response.WriteError(c, rpc_error.BadRequestDefault(err.Error()))
			return
		}
	}

	rpcReq := &bizpb.MoveProjectGroupRequest{
		UserId:          user.Id,
		ProjectId:       projectId,
		GroupId:         groupId,
		PreviousGroupId: previousGroupId,
		NextGroupId:     nextGroupId,
		MoveToFirst:     req.MoveToFirst,
		MoveToLast:      req.MoveToLast,
	}
	rpcResp := &bizpb.MoveProjectGroupResponse{}
	if err := project_assets.MoveProjectGroup(c.Request.Context(), h.CoreBizClient(), rpcReq, rpcResp); err != nil {
		response.WriteError(c, err)
		return
	}

	response.OK(c, dto.NewMoveProjectGroupResponse(rpcResp))
}
