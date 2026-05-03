package profile

import (
	"github.com/gin-gonic/gin"

	"icw_core_api/internal/dto/project"
	"icw_core_api/internal/response"
	"icw_core_api/utils"
	bizDto "icw_core_biz/pkg/dto/project"
)

// UploadProjectThumbnail 上传项目缩略图
// @router /project/profile/thumbnail [POST]
func (h *Handler) UploadProjectThumbnail(c *gin.Context) {
	var req project.UploadProjectThumbnailRequest
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

	rpcReq := &bizDto.UploadProjectThumbnailRequest{
		UserId:      user.Id,
		ProjectId:   projectId,
		ContentType: "image/png",
	}
	rpcResp := &bizDto.UploadProjectThumbnailResponse{}
	if err := h.CoreBizCall(c.Request.Context(), "ProjectProfileService.UploadProjectThumbnail", rpcReq, rpcResp); err != nil {
		response.WriteError(c, err)
		return
	}

	response.OK(c, project.NewUploadProjectThumbnailResponse(rpcResp))
}
