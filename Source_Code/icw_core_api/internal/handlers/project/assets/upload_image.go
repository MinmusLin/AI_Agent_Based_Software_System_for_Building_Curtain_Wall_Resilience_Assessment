package assets

import (
	"github.com/gin-gonic/gin"

	"icw_core_api/internal/dto/project"
	"icw_core_api/internal/response"
	"icw_core_api/utils"
	bizDto "icw_core_biz/pkg/dto/project"
)

// UploadProjectImage 上传图像
// @router /project/assets/image/upload [POST]
func (h *Handler) UploadProjectImage(c *gin.Context) {
	var req project.UploadProjectImageRequest
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
	groupId, err := utils.Decode(req.GroupId)
	if err != nil {
		response.WriteError(c, err)
		return
	}

	rpcReq := &bizDto.UploadProjectImageRequest{
		UserId:    user.Id,
		ProjectId: projectId,
		GroupId:   groupId,
		Images:    make([]*bizDto.UploadProjectImageItem, 0, len(req.Images)),
	}
	for _, image := range req.Images {
		if image == nil {
			continue
		}
		rpcReq.Images = append(rpcReq.Images, project.NewUploadProjectImageItem(image))
	}
	rpcResp := &bizDto.UploadProjectImageResponse{}
	if err := h.CoreBizCall(c.Request.Context(), "ProjectAssetsService.UploadProjectImage", rpcReq, rpcResp); err != nil {
		response.WriteError(c, err)
		return
	}

	response.OK(c, project.NewUploadProjectImageResponse(rpcResp))
}
