package detection

import (
	"strings"

	"github.com/gin-gonic/gin"

	"icw_common/gen/core/biz"
	"icw_common/rpc/error"
	"icw_common/utils"

	"icw_core_api/internal/dto"
	"icw_core_api/internal/response"
	"icw_core_api/rpc/icw_core_biz/project_detection"
	apiUtils "icw_core_api/utils"
)

// GetImageDetectionResult 获取图像检测结果
func (h *Handler) GetImageDetectionResult(c *gin.Context) {
	// 从 Gin Context 中获取当前登录用户
	user, err := apiUtils.GetCurrentUser(c)
	if err != nil {
		response.WriteError(c, err)
		return
	}

	// 将 Sqids 字符串解码为数字 ID
	projectId, err := utils.Decode(c.Query("project_id"))
	if err != nil {
		response.WriteError(c, rpc_error.BadRequestDefault(err.Error()))
		return
	}

	imageUuid := strings.TrimSpace(c.Query("image_uuid"))
	if imageUuid == "" {
		response.WriteError(c, rpc_error.BadRequestDefault("image uuid is required"))
		return
	}

	rpcReq := &bizpb.GetImageDetectionResultRequest{
		UserId:    user.Id,
		ProjectId: projectId,
		ImageUuid: imageUuid,
	}
	rpcResp := &bizpb.GetImageDetectionResultResponse{}
	if err := project_detection.GetImageDetectionResult(c.Request.Context(), h.CoreBizClient(), rpcReq, rpcResp); err != nil {
		response.WriteError(c, err)
		return
	}

	response.OK(c, dto.NewGetImageDetectionResultResponse(rpcResp))
}
