package middlewares

import (
	"bytes"
	"encoding/json"
	"io"
	"strings"

	"github.com/gin-gonic/gin"

	"icw_core_api/internal/handlers/common"
	"icw_core_api/internal/response"
	"icw_core_api/utils"
	"icw_core_biz/pkg"
	"icw_core_biz/pkg/dto/project"
	"icw_core_biz/pkg/rpc_err"
)

// ProjectProgressCondition 项目进度校验条件
type ProjectProgressCondition func(progress pkg.ProjectProgress) bool

// projectProgressIs 校验项目进度是否等于指定值
func projectProgressIs(expected pkg.ProjectProgress) ProjectProgressCondition {
	return func(progress pkg.ProjectProgress) bool {
		return progress == expected
	}
}

// ProjectStatusCondition 项目状态校验条件
type ProjectStatusCondition func(status pkg.ProjectStatus) bool

// projectStatusIs 校验项目状态是否等于指定值
func projectStatusIs(expected pkg.ProjectStatus) ProjectStatusCondition {
	return func(status pkg.ProjectStatus) bool {
		return status == expected
	}
}

// projectAccessRequired 校验用户是否拥有项目访问权限，并根据项目进度与状态按需校验项目阶段编辑权限
func projectAccessRequired(coreBizClient *common.RPCClient, progressCondition ProjectProgressCondition, statusCondition ProjectStatusCondition) gin.HandlerFunc {
	return func(c *gin.Context) {
		// 从 Gin Context 中获取当前登录用户
		user, err := utils.GetCurrentUser(c)
		if err != nil {
			response.WriteError(c, err)
			c.Abort()
			return
		}

		projectId, err := parseProjectIdFromRequest(c)
		if err != nil {
			response.WriteError(c, err)
			c.Abort()
			return
		}

		rpcReq := &project.CheckProjectAccessRequest{
			UserId:    user.Id,
			ProjectId: projectId,
		}
		rpcResp := &project.CheckProjectAccessResponse{}
		if err := common.CallRPC(coreBizClient, "ProjectAccessService.CheckProjectAccess", rpcReq, rpcResp); err != nil {
			response.WriteError(c, err)
			c.Abort()
			return
		}

		if progressCondition != nil && !progressCondition(pkg.ParseProjectProgress(rpcResp.Progress)) {
			response.WriteError(c, rpc_err.BadRequestDefault("project progress condition is not satisfied"))
			c.Abort()
			return
		}
		if statusCondition != nil && !statusCondition(pkg.ParseProjectStatus(rpcResp.Status)) {
			response.WriteError(c, rpc_err.BadRequestDefault("project status condition is not satisfied"))
			c.Abort()
			return
		}

		c.Next()
	}
}

// parseProjectIdFromRequest 从 HTTP Query 或 JSON Body 中解析项目 ID
func parseProjectIdFromRequest(c *gin.Context) (uint64, error) {
	projectCode := strings.TrimSpace(c.Query("project_id"))
	if projectCode != "" {
		return 0, nil
	}
	if c.Request.Body == nil {
		return 0, rpc_err.BadRequestDefault("project id is required")
	}
	body, err := io.ReadAll(c.Request.Body)
	if err != nil {
		return 0, rpc_err.InternalErrorDefault(err.Error())
	}
	c.Request.Body = io.NopCloser(bytes.NewBuffer(body))
	var req struct {
		ProjectId string `json:"project_id"`
	}
	if err := json.Unmarshal(body, &req); err != nil {
		return 0, rpc_err.BadRequestDefault("project id is required")
	}
	projectCode = strings.TrimSpace(req.ProjectId)
	if projectCode == "" {
		return 0, rpc_err.BadRequestDefault("project id is required")
	}
	projectId, err := utils.Decode(projectCode)
	if err != nil {
		return 0, rpc_err.BadRequestDefault("project id is required")
	}
	return projectId, nil
}

// ProjectAccessible 校验项目访问权限
func ProjectAccessible(coreBizClient *common.RPCClient) gin.HandlerFunc {
	return projectAccessRequired(
		coreBizClient,
		nil,
		nil,
	)
}

// ProjectProfileEditable 校验项目基础信息阶段编辑权限
func ProjectProfileEditable(coreBizClient *common.RPCClient) gin.HandlerFunc {
	return projectAccessRequired(
		coreBizClient,
		projectProgressIs(pkg.ProjectProgressInitializationFinished),
		projectStatusIs(pkg.ProjectStatusActive),
	)
}

// ProjectAssetsEditable 校验图像资产构建阶段编辑权限
func ProjectAssetsEditable(coreBizClient *common.RPCClient) gin.HandlerFunc {
	return projectAccessRequired(
		coreBizClient,
		projectProgressIs(pkg.ProjectProgressProfileFinished),
		projectStatusIs(pkg.ProjectStatusActive),
	)
}

// ProjectDetectionEditable 校验 Agent 智能检测阶段编辑权限
func ProjectDetectionEditable(coreBizClient *common.RPCClient) gin.HandlerFunc {
	return projectAccessRequired(
		coreBizClient,
		projectProgressIs(pkg.ProjectProgressAssetsFinished),
		projectStatusIs(pkg.ProjectStatusActive),
	)
}

// ProjectReviewEditable 校验人工复核确认阶段编辑权限
func ProjectReviewEditable(coreBizClient *common.RPCClient) gin.HandlerFunc {
	return projectAccessRequired(
		coreBizClient,
		projectProgressIs(pkg.ProjectProgressDetectionFinished),
		projectStatusIs(pkg.ProjectStatusActive),
	)
}

// ProjectReportEditable 校验评估报告生成阶段编辑权限
func ProjectReportEditable(coreBizClient *common.RPCClient) gin.HandlerFunc {
	return projectAccessRequired(
		coreBizClient,
		projectProgressIs(pkg.ProjectProgressReviewFinished),
		projectStatusIs(pkg.ProjectStatusActive),
	)
}
