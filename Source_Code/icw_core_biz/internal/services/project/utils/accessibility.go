package utils

import (
	"context"

	"icw_core_biz/pkg/rpc_err"
	"icw_core_biz/repositories/mysql"
)

// ValidateProjectOwnership 校验用户是否拥有项目访问权限
func ValidateProjectOwnership(ctx context.Context, repo *mysql.Repository, userId, projectId uint64) (*mysql.ProjectRecord, error) {
	if userId == 0 || projectId == 0 {
		return nil, rpc_err.BadRequestDefault("invalid project ownership request")
	}

	// 按用户 ID 和项目 ID 查询项目
	project, err := repo.FindProjectByIdAndUserId(ctx, userId, projectId)
	if err != nil {
		return nil, err
	}

	// 如果项目不存在 / 项目不属于该用户 / 项目状态为已删除，则用户无项目访问权限
	if project == nil || project.Status == string(mysql.ProjectStatusDeleted) {
		return nil, rpc_err.BadRequest(rpc_err.DetailProjectNotAccessible, "project is not accessible")
	}

	return project, nil
}
