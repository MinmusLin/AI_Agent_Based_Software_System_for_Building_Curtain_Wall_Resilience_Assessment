package mysql

import (
	"context"
	"database/sql"

	"icw_common/enum"
	"icw_common/gen/core/biz"

	"icw_core_biz/repositories/mysql/model"
)

// ListProjectImageObjectReferences 查询所有项目图像对象引用
func (r *Repository) ListProjectImageObjectReferences(ctx context.Context) ([]model.ProjectImageObjectReference, error) {
	rows, err := r.mysql.QueryContext(ctx, `
		SELECT project_id, uuid
		FROM project_group_images
	`)
	return scanProjectImageObjectReferences(rows, err)
}

// ListFailedProjectImageObjectReferences 查询上传失败的项目图像对象引用
func (r *Repository) ListFailedProjectImageObjectReferences(ctx context.Context) ([]model.ProjectImageObjectReference, error) {
	rows, err := r.mysql.QueryContext(ctx, `
		SELECT project_id, uuid
		FROM project_group_images
		WHERE status = ?
	`, enum.ProjectImageStatusString(bizpb.ProjectImageStatus_Failed))
	return scanProjectImageObjectReferences(rows, err)
}

// ListProjectDetectionObjectReferences 查询所有项目检测产物对象引用
func (r *Repository) ListProjectDetectionObjectReferences(ctx context.Context) ([]model.ProjectDetectionObjectReference, error) {
	rows, err := r.mysql.QueryContext(ctx, projectDetectionObjectReferenceSQL(""))
	return scanProjectDetectionObjectReferences(rows, err)
}

// ListFailedProjectDetectionObjectReferences 查询检测失败的项目检测产物对象引用
func (r *Repository) ListFailedProjectDetectionObjectReferences(ctx context.Context) ([]model.ProjectDetectionObjectReference, error) {
	rows, err := r.mysql.QueryContext(
		ctx,
		projectDetectionObjectReferenceSQL("WHERE t.status = ?"),
		enum.ProjectDetectionSubTaskStatusString(bizpb.ProjectDetectionSubTaskStatus_Failed),
		enum.ProjectDetectionSubTaskStatusString(bizpb.ProjectDetectionSubTaskStatus_Failed),
		enum.ProjectDetectionSubTaskStatusString(bizpb.ProjectDetectionSubTaskStatus_Failed),
		enum.ProjectDetectionSubTaskStatusString(bizpb.ProjectDetectionSubTaskStatus_Failed),
		enum.ProjectDetectionSubTaskStatusString(bizpb.ProjectDetectionSubTaskStatus_Failed),
	)
	return scanProjectDetectionObjectReferences(rows, err)
}

// ListProjectReportObjectReferences 查询所有项目报告对象引用
func (r *Repository) ListProjectReportObjectReferences(ctx context.Context) ([]model.ProjectReportObjectReference, error) {
	rows, err := r.mysql.QueryContext(ctx, `
		SELECT project_id
		FROM project_reports
	`)
	return scanProjectReportObjectReferences(rows, err)
}

// ListFailedProjectReportObjectReferences 查询生成失败的项目报告对象引用
func (r *Repository) ListFailedProjectReportObjectReferences(ctx context.Context) ([]model.ProjectReportObjectReference, error) {
	rows, err := r.mysql.QueryContext(ctx, `
		SELECT project_id
		FROM project_reports
		WHERE status = ?
	`, enum.ProjectReportStatusString(bizpb.ProjectReportStatus_Failed))
	return scanProjectReportObjectReferences(rows, err)
}

// projectDetectionObjectReferenceSQL 构造图像检测子任务对象引用查询 SQL
func projectDetectionObjectReferenceSQL(whereClause string) string {
	return `
		SELECT t.project_id, i.uuid, 'corrosion' AS task_code
		FROM project_detection_corrosion_tasks t
		JOIN project_group_images i ON i.id = t.image_id
		` + whereClause + `
		UNION ALL
		SELECT t.project_id, i.uuid, 'crack' AS task_code
		FROM project_detection_crack_tasks t
		JOIN project_group_images i ON i.id = t.image_id
		` + whereClause + `
		UNION ALL
		SELECT t.project_id, i.uuid, 'stain' AS task_code
		FROM project_detection_stain_tasks t
		JOIN project_group_images i ON i.id = t.image_id
		` + whereClause + `
		UNION ALL
		SELECT t.project_id, i.uuid, 'flatness' AS task_code
		FROM project_detection_flatness_tasks t
		JOIN project_group_images i ON i.id = t.image_id
		` + whereClause + `
		UNION ALL
		SELECT t.project_id, i.uuid, 'spalling' AS task_code
		FROM project_detection_spalling_tasks t
		JOIN project_group_images i ON i.id = t.image_id
		` + whereClause + `
	`
}

// scanProjectImageObjectReferences 扫描项目图像对象引用查询结果
func scanProjectImageObjectReferences(rows *sql.Rows, err error) ([]model.ProjectImageObjectReference, error) {
	if err != nil {
		return nil, err
	}

	defer func() {
		_ = rows.Close()
	}()

	items := make([]model.ProjectImageObjectReference, 0)
	for rows.Next() {
		var item model.ProjectImageObjectReference
		if err := rows.Scan(&item.ProjectId, &item.ImageUuid); err != nil {
			return nil, err
		}
		items = append(items, item)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return items, nil
}

// scanProjectDetectionObjectReferences 扫描项目检测产物对象引用查询结果
func scanProjectDetectionObjectReferences(rows *sql.Rows, err error) ([]model.ProjectDetectionObjectReference, error) {
	if err != nil {
		return nil, err
	}

	defer func() {
		_ = rows.Close()
	}()

	items := make([]model.ProjectDetectionObjectReference, 0)
	for rows.Next() {
		var item model.ProjectDetectionObjectReference
		if err := rows.Scan(&item.ProjectId, &item.ImageUuid, &item.TaskCode); err != nil {
			return nil, err
		}
		items = append(items, item)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return items, nil
}

// scanProjectReportObjectReferences 扫描项目报告对象引用查询结果
func scanProjectReportObjectReferences(rows *sql.Rows, err error) ([]model.ProjectReportObjectReference, error) {
	if err != nil {
		return nil, err
	}

	defer func() {
		_ = rows.Close()
	}()

	items := make([]model.ProjectReportObjectReference, 0)
	for rows.Next() {
		var item model.ProjectReportObjectReference
		if err := rows.Scan(&item.ProjectId); err != nil {
			return nil, err
		}
		items = append(items, item)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return items, nil
}
