package mysql

import (
	"context"

	"github.com/google/uuid"

	"icw_common/enum"
	"icw_common/gen/core/biz"
	"icw_common/utils"

	"icw_core_biz/repositories/mysql/model"
	mysqlUtils "icw_core_biz/repositories/mysql/utils"
)

// CreateProjectReport 按用户 ID 和项目 ID 创建或复用项目评估报告记录
func (r *Repository) CreateProjectReport(ctx context.Context, userId, projectId uint64) (*model.ProjectReportRecord, error) {
	_, err := r.mysql.ExecContext(ctx, `
		INSERT INTO project_reports(uuid, user_id, project_id, status)
		VALUES (?, ?, ?, ?)
		ON DUPLICATE KEY UPDATE
			status = VALUES(status),
			result_json = NULL
	`, uuid.NewString(), userId, projectId, enum.ProjectReportStatusString(bizpb.ProjectReportStatus_Pending))
	if err != nil {
		return nil, err
	}

	return r.GetProjectReport(ctx, userId, projectId)
}

// GetProjectReport 按用户 ID 和项目 ID 查询项目评估报告记录
func (r *Repository) GetProjectReport(ctx context.Context, userId, projectId uint64) (*model.ProjectReportRecord, error) {
	return model.ScanProjectReport(r.mysql.QueryRowContext(ctx, `
		SELECT id, uuid, user_id, project_id, status, CAST(COALESCE(result_json, JSON_OBJECT()) AS CHAR), created_at, updated_at
		FROM project_reports
		WHERE user_id = ? AND project_id = ?
	`, userId, projectId))
}

// UpdateProjectReportResult 按项目 ID 更新项目评估报告结果
func (r *Repository) UpdateProjectReportResult(ctx context.Context, projectId uint64, status bizpb.ProjectReportStatus_Value, resultJSON string) (*model.ProjectReportRecord, error) {
	statusText := enum.ProjectReportStatusString(status)
	if statusText == "" {
		return nil, model.ErrProjectReportStatusInvalid
	}
	compactedJSON := ""
	if status == bizpb.ProjectReportStatus_Succeeded {
		var err error
		compactedJSON, err = utils.CompactJSONObjectString(resultJSON)
		if err != nil {
			return nil, err
		}
	}

	tx, err := r.mysql.BeginTx(ctx, nil)
	if err != nil {
		return nil, err
	}

	defer func() {
		_ = tx.Rollback()
	}()

	result, err := tx.ExecContext(ctx, `
		UPDATE project_reports
		SET status = ?,
			result_json = CASE WHEN ? = ? THEN ? ELSE result_json END
		WHERE project_id = ?
	`, statusText, statusText, enum.ProjectReportStatusString(bizpb.ProjectReportStatus_Succeeded), compactedJSON, projectId)
	if err := mysqlUtils.CheckRowsAffected(result, err); err != nil {
		return nil, err
	}

	if status == bizpb.ProjectReportStatus_Succeeded {
		result, err = tx.ExecContext(ctx, `
			UPDATE projects
			SET progress = ?, status = ?
			WHERE id = ?
		`, enum.ProjectProgressUint8(bizpb.ProjectProgress_ReportFinished), enum.ProjectStatusString(bizpb.ProjectStatus_Completed), projectId)
		if err := mysqlUtils.CheckRowsAffected(result, err); err != nil {
			return nil, err
		}
	}

	if err := tx.Commit(); err != nil {
		return nil, err
	}

	var reportUserId uint64
	if err := r.mysql.QueryRowContext(ctx, `
		SELECT user_id
		FROM project_reports
		WHERE project_id = ?
	`, projectId).Scan(&reportUserId); err != nil {
		return nil, err
	}

	return r.GetProjectReport(ctx, reportUserId, projectId)
}
