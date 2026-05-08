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

// DeleteProjectImage 删除图像
// @router /project/assets/image/delete [POST]
func (h *Handler) DeleteProjectImage(c *gin.Context) {
	var req apipb.DeleteProjectImageRequest
	if !response.BindJSON(c, &req) {
		return
	}

	// 从 Gin Context 中获取当前登录用户
	user, err := apiUtils.GetCurrentUser(c)
	if err != nil {
		response.Error(c, err)
		return
	}

	// 将 Sqids 字符串解码为数字 ID
	projectId, err := utils.Decode(req.ProjectId)
	if err != nil {
		response.Error(c, rpc_error.BadRequestDefault(err.Error()))
		return
	}

	rpcReq := &bizpb.DeleteProjectImageRequest{
		UserId:     user.Id,
		ProjectId:  projectId,
		ImageUuids: req.ImageUuids,
	}
	rpcResp := &bizpb.DeleteProjectImageResponse{}
	if err := project_assets.DeleteProjectImage(c.Request.Context(), h.CoreBizClient(), rpcReq, rpcResp); err != nil {
		response.Error(c, err)
		return
	}

	response.OK(c, dto.NewDeleteProjectImageResponse(rpcResp))
}
