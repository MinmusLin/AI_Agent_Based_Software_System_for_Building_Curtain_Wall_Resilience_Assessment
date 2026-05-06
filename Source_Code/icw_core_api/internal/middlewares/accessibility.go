package middlewares

import (
	"bytes"
	"encoding/json"
	"io"

	"github.com/gin-gonic/gin"

	"icw_common/enum"
	"icw_common/gen/core/biz"
	"icw_common/rpc_err"
	"icw_common/utils"
	"icw_core_api/internal/response"
	"icw_core_api/rpc/icw_core_biz"
	"icw_core_api/rpc/icw_core_biz/project_core"
	apiUtils "icw_core_api/utils"
)

// ProjectProgressCondition 项目进度校验条件
type ProjectProgressCondition func(progress bizpb.ProjectProgress_Value) bool

// projectProgressIs 校验项目进度是否等于指定值
func projectProgressIs(expected bizpb.ProjectProgress_Value) ProjectProgressCondition {
	return func(progress bizpb.ProjectProgress_Value) bool {
		return progress == expected
	}
}

// ProjectStatusCondition 项目状态校验条件
type ProjectStatusCondition func(status bizpb.ProjectStatus_Value) bool

// projectStatusIs 校验项目状态是否等于指定值
func projectStatusIs(expected bizpb.ProjectStatus_Value) ProjectStatusCondition {
	return func(status bizpb.ProjectStatus_Value) bool {
		return status == expected
	}
}

// projectAccessRequired 校验用户是否拥有项目访问权限，并根据项目进度与状态按需校验项目阶段编辑权限
func projectAccessRequired(coreBizClient *icw_core_biz.Client, progressCondition ProjectProgressCondition, statusCondition ProjectStatusCondition) gin.HandlerFunc {
	return func(c *gin.Context) {
		// 从 Gin Context 中获取当前登录用户
		user, err := apiUtils.GetCurrentUser(c)
		if err != nil {
			response.WriteError(c, err)
			c.Abort()
			return
		}

		// 从 HTTP Query 或 JSON Body 中解析项目 ID
		projectId, err := parseProjectIdFromRequest(c)
		if err != nil {
			response.WriteError(c, err)
			c.Abort()
			return
		}

		rpcReq := &bizpb.CheckProjectAccessRequest{
			UserId:    user.Id,
			ProjectId: projectId,
		}
		rpcResp := &bizpb.CheckProjectAccessResponse{}
		if err := project_core.CheckProjectAccess(c.Request.Context(), coreBizClient, rpcReq, rpcResp); err != nil {
			response.WriteError(c, err)
			c.Abort()
			return
		}

		// 根据项目进度按需校验项目阶段编辑权限
		if progressCondition != nil && !progressCondition(enum.ParseProjectProgress(uint8(rpcResp.Progress))) {
			response.WriteError(c, rpc_err.BadRequestDefault("project progress condition is not satisfied"))
			c.Abort()
			return
		}

		// 根据项目状态按需校验项目阶段编辑权限
		if statusCondition != nil && !statusCondition(enum.ParseProjectStatus(rpcResp.Status)) {
			response.WriteError(c, rpc_err.BadRequestDefault("project status condition is not satisfied"))
			c.Abort()
			return
		}

		c.Next()
	}
}

// parseProjectIdFromRequest 从 HTTP Query 或 JSON Body 中解析项目 ID
func parseProjectIdFromRequest(c *gin.Context) (uint64, error) {
	// 从 HTTP Query 中解析项目 ID
	projectCode := c.Query("project_id")
	if projectCode != "" {
		projectId, err := utils.Decode(projectCode)
		if err != nil {
			return 0, rpc_err.BadRequestDefault(err.Error())
		}
		return projectId, nil
	}

	// 从 JSON Body 中解析项目 ID
	if c.Request.Body == nil {
		return 0, rpc_err.BadRequestDefault("project id is required")
	}
	body, err := io.ReadAll(c.Request.Body)
	if err != nil {
		return 0, err
	}
	c.Request.Body = io.NopCloser(bytes.NewBuffer(body))
	var req struct {
		ProjectId string `json:"project_id"`
	}
	if err := json.Unmarshal(body, &req); err != nil {
		return 0, rpc_err.BadRequestDefault(err.Error())
	}
	projectCode = req.ProjectId
	if projectCode == "" {
		return 0, rpc_err.BadRequestDefault("project id is required")
	}
	projectId, err := utils.Decode(projectCode)
	if err != nil {
		return 0, rpc_err.BadRequestDefault(err.Error())
	}

	return projectId, nil
}

// ProjectAccessible 校验项目访问权限
func ProjectAccessible(coreBizClient *icw_core_biz.Client) gin.HandlerFunc {
	return projectAccessRequired(
		coreBizClient,
		nil,
		nil,
	)
}

// ProjectProfileEditable 校验项目基础信息阶段编辑权限
func ProjectProfileEditable(coreBizClient *icw_core_biz.Client) gin.HandlerFunc {
	return projectAccessRequired(
		coreBizClient,
		projectProgressIs(bizpb.ProjectProgress_InitializationFinished),
		projectStatusIs(bizpb.ProjectStatus_Active),
	)
}

// ProjectAssetsEditable 校验图像资产构建阶段编辑权限
func ProjectAssetsEditable(coreBizClient *icw_core_biz.Client) gin.HandlerFunc {
	return projectAccessRequired(
		coreBizClient,
		projectProgressIs(bizpb.ProjectProgress_ProfileFinished),
		projectStatusIs(bizpb.ProjectStatus_Active),
	)
}

// ProjectDetectionEditable 校验 Agent 智能检测阶段编辑权限
func ProjectDetectionEditable(coreBizClient *icw_core_biz.Client) gin.HandlerFunc {
	return projectAccessRequired(
		coreBizClient,
		projectProgressIs(bizpb.ProjectProgress_AssetsFinished),
		projectStatusIs(bizpb.ProjectStatus_Active),
	)
}

// ProjectReviewEditable 校验人工复核确认阶段编辑权限
func ProjectReviewEditable(coreBizClient *icw_core_biz.Client) gin.HandlerFunc {
	return projectAccessRequired(
		coreBizClient,
		projectProgressIs(bizpb.ProjectProgress_DetectionFinished),
		projectStatusIs(bizpb.ProjectStatus_Active),
	)
}

// ProjectReportEditable 校验评估报告生成阶段编辑权限
func ProjectReportEditable(coreBizClient *icw_core_biz.Client) gin.HandlerFunc {
	return projectAccessRequired(
		coreBizClient,
		projectProgressIs(bizpb.ProjectProgress_ReviewFinished),
		projectStatusIs(bizpb.ProjectStatus_Active),
	)
}
