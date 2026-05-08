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

// UploadProjectThumbnail 上传项目缩略图
// @router /project/profile/thumbnail [POST]
func (h *Handler) UploadProjectThumbnail(c *gin.Context) {
	var req apipb.UploadProjectThumbnailRequest
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

	rpcReq := &bizpb.UploadProjectThumbnailRequest{
		UserId:      user.Id,
		ProjectId:   projectId,
		ContentType: "image/png",
	}
	rpcResp := &bizpb.UploadProjectThumbnailResponse{}
	if err := project_profile.UploadProjectThumbnail(c.Request.Context(), h.CoreBizClient(), rpcReq, rpcResp); err != nil {
		response.Error(c, err)
		return
	}

	response.OK(c, dto.NewUploadProjectThumbnailResponse(rpcResp))
}
