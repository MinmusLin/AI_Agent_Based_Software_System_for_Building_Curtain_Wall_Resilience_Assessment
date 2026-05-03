package assets

import (
	"github.com/gin-gonic/gin"

	"icw_core_api/internal/dto/project"
	"icw_core_api/internal/response"
	"icw_core_api/utils"
	bizDto "icw_core_biz/pkg/dto/project"
)

// ReportProjectImage 上报图像
// @router /project/assets/image/report [POST]
func (h *Handler) ReportProjectImage(c *gin.Context) {
	var req project.ReportProjectImageRequest
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

	rpcReq := &bizDto.ReportProjectImageRequest{
		UserId:    user.Id,
		ProjectId: projectId,
		ImageUuid: req.ImageUuid,
		Status:    req.Status,
	}
	rpcResp := &bizDto.ReportProjectImageResponse{}
	if err := h.CallRPC(h.CoreBizRPCClient(), "ProjectAssetsService.ReportProjectImage", rpcReq, rpcResp); err != nil {
		response.WriteError(c, err)
		return
	}

	response.OK(c, project.NewReportProjectImageResponse(rpcResp))
}
