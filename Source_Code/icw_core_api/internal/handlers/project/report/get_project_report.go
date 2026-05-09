package report

import (
	"github.com/gin-gonic/gin"

	"icw_common/gen/core/biz"
	"icw_common/rpc/error"
	"icw_common/utils"

	"icw_core_api/internal/dto"
	"icw_core_api/internal/response"
	"icw_core_api/rpc/icw_core_biz/project_report"
	apiUtils "icw_core_api/utils"
)

// GetProjectReport 获取项目评估报告
// @router /project/report/detail [GET]
func (h *Handler) GetProjectReport(c *gin.Context) {
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

	rpcReq := &bizpb.GetProjectReportRequest{
		UserId:    user.Id,
		ProjectId: projectId,
	}
	rpcResp := &bizpb.GetProjectReportResponse{}
	if err := project_report.GetProjectReport(c.Request.Context(), h.CoreBizClient(), rpcReq, rpcResp); err != nil {
		response.Error(c, err)
		return
	}

	response.OK(c, dto.NewGetProjectReportResponse(rpcResp))
}
