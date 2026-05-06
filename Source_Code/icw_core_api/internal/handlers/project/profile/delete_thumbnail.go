package profile

import (
	"github.com/gin-gonic/gin"

	"icw_common/gen/core/biz"
	"icw_common/rpc_err"
	"icw_common/utils"
	"icw_core_api/internal/response"
	"icw_core_api/rpc/icw_core_biz/project_profile"
	apiUtils "icw_core_api/utils"
)

// DeleteProjectThumbnail 删除项目缩略图
// @router /project/profile/thumbnail [DELETE]
func (h *Handler) DeleteProjectThumbnail(c *gin.Context) {
	// 从 Gin Context 中获取当前登录用户
	user, err := apiUtils.GetCurrentUser(c)
	if err != nil {
		response.WriteError(c, err)
		return
	}

	// 将 Sqids 字符串解码为数字 ID
	projectId, err := utils.Decode(c.Query("project_id"))
	if err != nil {
		response.WriteError(c, rpc_err.BadRequestDefault(err.Error()))
		return
	}

	rpcReq := &bizpb.DeleteProjectThumbnailRequest{
		UserId:    user.Id,
		ProjectId: projectId,
	}
	rpcResp := &bizpb.DeleteProjectThumbnailResponse{}
	if err := project_profile.DeleteProjectThumbnail(c.Request.Context(), h.CoreBizClient(), rpcReq, rpcResp); err != nil {
		response.WriteError(c, err)
		return
	}

	response.OK(c, apiUtils.NewDeleteProjectThumbnailResponse(rpcResp))
}
