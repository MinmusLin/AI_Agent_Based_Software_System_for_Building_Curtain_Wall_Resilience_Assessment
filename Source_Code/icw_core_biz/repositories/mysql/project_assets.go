package mysql

import (
	"context"
	"database/sql"
	"errors"

	"icw_core_biz/internal/services/project/consts"
)

// ListProjectGroups 按用户 ID 和项目 ID 查询图像组列表
func (r *Repository) ListProjectGroups(ctx context.Context, userId, projectId uint64) ([]*ProjectGroupRecord, error) {
	rows, err := r.mysql.QueryContext(ctx, `
		SELECT id, project_id, user_id, name, CAST(sort_order AS CHAR), created_at, updated_at
		FROM project_groups
		WHERE user_id = ? AND project_id = ?
		ORDER BY sort_order ASC, created_at ASC, id ASC
	`, userId, projectId)
	if err != nil {
		return nil, err
	}
	defer func() {
		_ = rows.Close()
	}()

	groups := make([]*ProjectGroupRecord, 0)
	for rows.Next() {
		group := &ProjectGroupRecord{}
		if err := rows.Scan(
			&group.Id,
			&group.ProjectId,
			&group.UserId,
			&group.Name,
			&group.SortOrder,
			&group.CreatedAt,
			&group.UpdatedAt,
		); err != nil {
			return nil, err
		}
		groups = append(groups, group)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return groups, nil
}

// ListProjectImages 按用户 ID 和项目 ID 查询图像列表
func (r *Repository) ListProjectImages(ctx context.Context, userId, projectId uint64) ([]*ProjectImageRecord, error) {
	rows, err := r.mysql.QueryContext(ctx, `
		SELECT
			id,
			group_id,
			project_id,
			user_id,
			uuid,
			file_name,
			content_type,
			size_bytes,
			width,
			height,
			CAST(metadata AS CHAR),
			status,
			uploaded_at,
			created_at,
			updated_at
		FROM project_group_images
		WHERE user_id = ? AND project_id = ?
		ORDER BY created_at ASC, id ASC
	`, userId, projectId)
	if err != nil {
		return nil, err
	}
	defer func() {
		_ = rows.Close()
	}()

	images := make([]*ProjectImageRecord, 0)
	for rows.Next() {
		image := &ProjectImageRecord{}
		var status string
		if err := rows.Scan(
			&image.Id,
			&image.GroupId,
			&image.ProjectId,
			&image.UserId,
			&image.Uuid,
			&image.FileName,
			&image.ContentType,
			&image.SizeBytes,
			&image.Width,
			&image.Height,
			&image.Metadata,
			&status,
			&image.UploadedAt,
			&image.CreatedAt,
			&image.UpdatedAt,
		); err != nil {
			return nil, err
		}
		image.Status = consts.ParseProjectImageStatus(status)
		images = append(images, image)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return images, nil
}

// ListProjectImagesByGroupId 按用户 ID、项目 ID 和图像组 ID 查询图像列表
func (r *Repository) ListProjectImagesByGroupId(ctx context.Context, userId, projectId, groupId uint64) ([]*ProjectImageRecord, error) {
	rows, err := r.mysql.QueryContext(ctx, `
		SELECT
			id,
			group_id,
			project_id,
			user_id,
			uuid,
			file_name,
			content_type,
			size_bytes,
			width,
			height,
			CAST(metadata AS CHAR),
			status,
			uploaded_at,
			created_at,
			updated_at
		FROM project_group_images
		WHERE user_id = ? AND project_id = ? AND group_id = ?
		ORDER BY created_at ASC, id ASC
	`, userId, projectId, groupId)
	if err != nil {
		return nil, err
	}
	defer func() {
		_ = rows.Close()
	}()

	images := make([]*ProjectImageRecord, 0)
	for rows.Next() {
		image := &ProjectImageRecord{}
		var status string
		if err := rows.Scan(
			&image.Id,
			&image.GroupId,
			&image.ProjectId,
			&image.UserId,
			&image.Uuid,
			&image.FileName,
			&image.ContentType,
			&image.SizeBytes,
			&image.Width,
			&image.Height,
			&image.Metadata,
			&status,
			&image.UploadedAt,
			&image.CreatedAt,
			&image.UpdatedAt,
		); err != nil {
			return nil, err
		}
		image.Status = consts.ParseProjectImageStatus(status)
		images = append(images, image)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return images, nil
}

// FindProjectGroupById 按用户 ID、项目 ID 和图像组 ID 查询图像组
func (r *Repository) FindProjectGroupById(ctx context.Context, userId, projectId, groupId uint64) (*ProjectGroupRecord, error) {
	group := &ProjectGroupRecord{}

	err := r.mysql.QueryRowContext(ctx, `
		SELECT id, project_id, user_id, name, CAST(sort_order AS CHAR), created_at, updated_at
		FROM project_groups
		WHERE id = ? AND user_id = ? AND project_id = ?
		LIMIT 1
	`, groupId, userId, projectId).Scan(
		&group.Id,
		&group.ProjectId,
		&group.UserId,
		&group.Name,
		&group.SortOrder,
		&group.CreatedAt,
		&group.UpdatedAt,
	)

	if errors.Is(err, sql.ErrNoRows) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	return group, nil
}

// FindProjectImageById 按用户 ID、项目 ID 和图像 ID 查询图像
func (r *Repository) FindProjectImageById(ctx context.Context, userId, projectId, imageId uint64) (*ProjectImageRecord, error) {
	image := &ProjectImageRecord{}
	var status string

	err := r.mysql.QueryRowContext(ctx, `
		SELECT
			id,
			group_id,
			project_id,
			user_id,
			uuid,
			file_name,
			content_type,
			size_bytes,
			width,
			height,
			CAST(metadata AS CHAR),
			status,
			uploaded_at,
			created_at,
			updated_at
		FROM project_group_images
		WHERE id = ? AND user_id = ? AND project_id = ?
		LIMIT 1
	`, imageId, userId, projectId).Scan(
		&image.Id,
		&image.GroupId,
		&image.ProjectId,
		&image.UserId,
		&image.Uuid,
		&image.FileName,
		&image.ContentType,
		&image.SizeBytes,
		&image.Width,
		&image.Height,
		&image.Metadata,
		&status,
		&image.UploadedAt,
		&image.CreatedAt,
		&image.UpdatedAt,
	)

	if errors.Is(err, sql.ErrNoRows) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	image.Status = consts.ParseProjectImageStatus(status)

	return image, nil
}

// FindProjectImageByUuid 按用户 ID、项目 ID 和图像 UUID 查询图像
func (r *Repository) FindProjectImageByUuid(ctx context.Context, userId, projectId uint64, imageUuid string) (*ProjectImageRecord, error) {
	image := &ProjectImageRecord{}
	var status string

	err := r.mysql.QueryRowContext(ctx, `
		SELECT
			id,
			group_id,
			project_id,
			user_id,
			uuid,
			file_name,
			content_type,
			size_bytes,
			width,
			height,
			CAST(metadata AS CHAR),
			status,
			uploaded_at,
			created_at,
			updated_at
		FROM project_group_images
		WHERE uuid = ? AND user_id = ? AND project_id = ?
		LIMIT 1
	`, imageUuid, userId, projectId).Scan(
		&image.Id,
		&image.GroupId,
		&image.ProjectId,
		&image.UserId,
		&image.Uuid,
		&image.FileName,
		&image.ContentType,
		&image.SizeBytes,
		&image.Width,
		&image.Height,
		&image.Metadata,
		&status,
		&image.UploadedAt,
		&image.CreatedAt,
		&image.UpdatedAt,
	)

	if errors.Is(err, sql.ErrNoRows) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	image.Status = consts.ParseProjectImageStatus(status)

	return image, nil
}

// CreateProjectGroup 按用户 ID 和项目 ID 创建图像组
func (r *Repository) CreateProjectGroup(ctx context.Context, userId, projectId uint64, name string) (*ProjectGroupRecord, error) {
	result, err := r.mysql.ExecContext(ctx, `
		INSERT INTO project_groups(project_id, user_id, name, sort_order)
		SELECT ?, ?, ?, COALESCE(MAX(sort_order), -1) + 1
		FROM project_groups
		WHERE project_id = ? AND user_id = ?
	`, projectId, userId, name, projectId, userId)
	if err != nil {
		return nil, err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return nil, err
	}

	return r.FindProjectGroupById(ctx, userId, projectId, uint64(id))
}

// CreateProjectImage 按用户 ID、项目 ID 和图像组 ID 创建图像
func (r *Repository) CreateProjectImage(ctx context.Context, userId, projectId, groupId uint64, imageUuid, fileName, contentType string, sizeBytes uint64, width, height uint32, metadata string) (*ProjectImageRecord, error) {
	result, err := r.mysql.ExecContext(ctx, `
		INSERT INTO project_group_images(
			group_id,
			project_id,
			user_id,
			uuid,
			file_name,
			content_type,
			size_bytes,
			width,
			height,
			metadata,
			status
		)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
	`, groupId, projectId, userId, imageUuid, fileName, contentType, sizeBytes, width, height, metadata, consts.ProjectImageStatusPending.String())
	if err != nil {
		return nil, err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return nil, err
	}

	return r.FindProjectImageById(ctx, userId, projectId, uint64(id))
}

// DeleteProjectGroup 按用户 ID、项目 ID 和图像组 ID 删除图像组
func (r *Repository) DeleteProjectGroup(ctx context.Context, userId, projectId, groupId uint64) (bool, error) {
	tx, err := r.mysql.BeginTx(ctx, nil)
	if err != nil {
		return false, err
	}
	defer func() {
		_ = tx.Rollback()
	}()

	var groupCount uint64
	if err := tx.QueryRowContext(ctx, `
		SELECT COUNT(*)
		FROM project_groups
		WHERE user_id = ? AND project_id = ?
		FOR UPDATE
	`, userId, projectId).Scan(&groupCount); err != nil {
		return false, err
	}
	if groupCount <= 1 {
		return false, ErrProjectGroupCannotDeleteLast
	}

	if _, err := tx.ExecContext(ctx, `
		DELETE FROM project_group_images
		WHERE user_id = ? AND project_id = ? AND group_id = ?
	`, userId, projectId, groupId); err != nil {
		return false, err
	}

	result, err := tx.ExecContext(ctx, `
		DELETE FROM project_groups
		WHERE id = ? AND user_id = ? AND project_id = ?
	`, groupId, userId, projectId)
	if err != nil {
		return false, err
	}

	affected, err := result.RowsAffected()
	if err != nil {
		return false, err
	}
	if affected == 0 {
		return false, nil
	}

	if err := tx.Commit(); err != nil {
		return false, err
	}

	return true, nil
}

// DeleteProjectImage 按用户 ID、项目 ID 和图像 UUID 删除项目图像记录
func (r *Repository) DeleteProjectImage(ctx context.Context, userId, projectId uint64, imageUuid string) (bool, error) {
	result, err := r.mysql.ExecContext(ctx, `
		DELETE FROM project_group_images
		WHERE uuid = ? AND user_id = ? AND project_id = ?
	`, imageUuid, userId, projectId)
	if err != nil {
		return false, err
	}

	affected, err := result.RowsAffected()
	if err != nil {
		return false, err
	}

	return affected > 0, nil
}

// MoveProjectGroup 按用户 ID、项目 ID 和图像组 ID 移动图像组
func (r *Repository) MoveProjectGroup(ctx context.Context, userId, projectId, groupId, previousGroupId, nextGroupId uint64, moveToFirst, moveToLast bool) (*ProjectGroupRecord, error) {
	var sortOrder string

	// 移动到置顶位置
	if moveToFirst && !moveToLast {
		err := r.mysql.QueryRowContext(ctx, `
			SELECT CAST(COALESCE(MIN(sort_order), 0) - 1 AS CHAR)
			FROM project_groups
			WHERE user_id = ? AND project_id = ?
		`, userId, projectId).Scan(&sortOrder)

		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		if err != nil {
			return nil, err
		}
	}

	// 移动到置底位置
	if !moveToFirst && moveToLast {
		err := r.mysql.QueryRowContext(ctx, `
			SELECT CAST(COALESCE(MAX(sort_order), 0) + 1 AS CHAR)
			FROM project_groups
			WHERE user_id = ? AND project_id = ?
		`, userId, projectId).Scan(&sortOrder)

		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		if err != nil {
			return nil, err
		}
	}

	// 移动到居中位置
	if !moveToFirst && !moveToLast {
		err := r.mysql.QueryRowContext(ctx, `
			SELECT CAST((previous_group.sort_order + next_group.sort_order) / 2 AS CHAR)
			FROM project_groups previous_group
			JOIN project_groups next_group
				ON next_group.id = ?
				AND next_group.user_id = ?
				AND next_group.project_id = ?
			WHERE previous_group.id = ?
				AND previous_group.user_id = ?
				AND previous_group.project_id = ?
		`, nextGroupId, userId, projectId, previousGroupId, userId, projectId).Scan(&sortOrder)

		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		if err != nil {
			return nil, err
		}
	}

	// 非法参数
	if moveToFirst && moveToLast {
		return nil, errors.New("move to first and move to last cannot both be true")
	}

	if sortOrder == "" {
		return nil, nil
	}

	result, err := r.mysql.ExecContext(ctx, `
		UPDATE project_groups
		SET sort_order = ?
		WHERE id = ? AND user_id = ? AND project_id = ?
	`, sortOrder, groupId, userId, projectId)
	if err != nil {
		return nil, err
	}

	if _, err = result.RowsAffected(); err != nil {
		return nil, err
	}

	return r.FindProjectGroupById(ctx, userId, projectId, groupId)
}

// MoveProjectImage 按用户 ID、项目 ID 和图像 UUID 移动图像
func (r *Repository) MoveProjectImage(ctx context.Context, userId, projectId uint64, imageUuid string, targetGroupId uint64) (*ProjectImageRecord, error) {
	result, err := r.mysql.ExecContext(ctx, `
		UPDATE project_group_images
		SET group_id = ?
		WHERE uuid = ? AND user_id = ? AND project_id = ?
	`, targetGroupId, imageUuid, userId, projectId)
	if err != nil {
		return nil, err
	}

	if _, err := result.RowsAffected(); err != nil {
		return nil, err
	}

	return r.FindProjectImageByUuid(ctx, userId, projectId, imageUuid)
}

// UpdateProjectGroupName 按用户 ID、项目 ID 和图像组 ID 更新图像组名称
func (r *Repository) UpdateProjectGroupName(ctx context.Context, userId, projectId, groupId uint64, name string) (*ProjectGroupRecord, error) {
	result, err := r.mysql.ExecContext(ctx, `
		UPDATE project_groups
		SET name = ?
		WHERE id = ? AND user_id = ? AND project_id = ?
	`, name, groupId, userId, projectId)
	if err != nil {
		return nil, err
	}

	if _, err := result.RowsAffected(); err != nil {
		return nil, err
	}

	return r.FindProjectGroupById(ctx, userId, projectId, groupId)
}

// UpdateProjectImageStatus 按用户 ID、项目 ID 和图像 UUID 更新图像状态
func (r *Repository) UpdateProjectImageStatus(ctx context.Context, userId, projectId uint64, imageUuid string, status consts.ProjectImageStatus) (*ProjectImageRecord, error) {
	result, err := r.mysql.ExecContext(ctx, `
		UPDATE project_group_images
		SET status = ?,
			uploaded_at = CASE WHEN ? = ? THEN NOW(3) ELSE NULL END
		WHERE uuid = ? AND user_id = ? AND project_id = ?
	`, status.String(), status.String(), consts.ProjectImageStatusUploaded.String(), imageUuid, userId, projectId)
	if err != nil {
		return nil, err
	}

	if _, err := result.RowsAffected(); err != nil {
		return nil, err
	}

	return r.FindProjectImageByUuid(ctx, userId, projectId, imageUuid)
}

// CountProjectGroups 按用户 ID 和项目 ID 统计图像组数量
func (r *Repository) CountProjectGroups(ctx context.Context, userId, projectId uint64) (uint64, error) {
	var count uint64

	err := r.mysql.QueryRowContext(ctx, `
		SELECT COUNT(*)
		FROM project_groups
		WHERE user_id = ? AND project_id = ?
	`, userId, projectId).Scan(&count)

	return count, err
}

// GetProjectAssetsReadyStats 按用户 ID 和项目 ID 统计项目图像状态校验
func (r *Repository) GetProjectAssetsReadyStats(ctx context.Context, userId, projectId uint64) (*ProjectAssetsReadyStats, error) {
	stats := &ProjectAssetsReadyStats{}

	err := r.mysql.QueryRowContext(ctx, `
		SELECT
			COUNT(CASE WHEN i.status = ? THEN 1 END) AS pending_images,
			COUNT(CASE WHEN i.status = ? THEN 1 END) AS uploaded_images,
			COUNT(CASE WHEN i.status = ? THEN 1 END) AS failed_images,
			COUNT(DISTINCT CASE WHEN i.id IS NULL THEN g.id END) AS empty_groups
		FROM project_groups g
		LEFT JOIN project_group_images i
			ON i.group_id = g.id
			AND i.project_id = g.project_id
			AND i.user_id = g.user_id
		WHERE g.user_id = ? AND g.project_id = ?
	`, consts.ProjectImageStatusPending.String(), consts.ProjectImageStatusUploaded.String(), consts.ProjectImageStatusFailed.String(), userId, projectId).Scan(
		&stats.PendingImageCount,
		&stats.UploadedImageCount,
		&stats.FailedImageCount,
		&stats.EmptyGroupCount,
	)
	if err != nil {
		return nil, err
	}

	return stats, nil
}
