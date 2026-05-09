package review

import (
	"strings"

	"github.com/gin-gonic/gin"

	"icw_common/gen/core/biz"
	"icw_common/rpc/error"
	"icw_common/utils"

	"icw_core_api/internal/dto"
	"icw_core_api/internal/response"
	"icw_core_api/rpc/icw_core_biz/project_review"
	apiUtils "icw_core_api/utils"
)

// GetProjectDetectionReview 获取图像检测人工复核信息
// @router /project/review/detail [GET]
func (h *Handler) GetProjectDetectionReview(c *gin.Context) {
	// 从 Gin Context 中获取当前登录用户
	user, err := apiUtils.GetCurrentUser(c)
	if err != nil {
		response.Error(c, err)
		return
	}

	// 将 Sqids 字符串解码为数字 ID
	projectId, err := utils.Decode(c.Query("project_id"))
	if err != nil {
		response.Error(c, rpc_error.BadRequestDefault(err.Error()))
		return
	}

	taskUuid := strings.TrimSpace(c.Query("task_uuid"))
	if taskUuid == "" {
		response.Error(c, rpc_error.BadRequestDefault("task uuid is required"))
		return
	}

	rpcReq := &bizpb.GetProjectDetectionReviewRequest{
		UserId:    user.Id,
		ProjectId: projectId,
		TaskUuid:  taskUuid,
	}
	rpcResp := &bizpb.GetProjectDetectionReviewResponse{}
	if err := project_review.GetProjectDetectionReview(c.Request.Context(), h.CoreBizClient(), rpcReq, rpcResp); err != nil {
		response.Error(c, err)
		return
	}

	response.OK(c, dto.NewGetProjectDetectionReviewResponse(rpcResp))
}
