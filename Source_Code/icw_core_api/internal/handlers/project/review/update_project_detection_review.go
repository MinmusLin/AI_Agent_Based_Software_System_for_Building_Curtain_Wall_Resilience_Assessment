package review

import (
	"github.com/gin-gonic/gin"

	"icw_common/gen/core/api"
	"icw_common/gen/core/biz"
	"icw_common/rpc/error"
	"icw_common/utils"

	"icw_core_api/internal/dto"
	"icw_core_api/internal/response"
	"icw_core_api/rpc/icw_core_biz/project_review"
	apiUtils "icw_core_api/utils"
)

// UpdateProjectDetectionReview 更新图像检测人工复核信息
// @router /project/review/update [POST]
func (h *Handler) UpdateProjectDetectionReview(c *gin.Context) {
	req := &apipb.UpdateProjectDetectionReviewRequest{}
	if !response.BindJSON(c, req) {
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

	rpcReq := &bizpb.UpdateProjectDetectionReviewRequest{
		UserId:    user.Id,
		ProjectId: projectId,
		TaskUuid:  req.TaskUuid,
		Verdict:   req.Verdict,
		Comment:   req.Comment,
	}
	rpcResp := &bizpb.UpdateProjectDetectionReviewResponse{}
	if err := project_review.UpdateProjectDetectionReview(c.Request.Context(), h.CoreBizClient(), rpcReq, rpcResp); err != nil {
		response.Error(c, err)
		return
	}

	response.OK(c, dto.NewUpdateProjectDetectionReviewResponse(rpcResp))
}
