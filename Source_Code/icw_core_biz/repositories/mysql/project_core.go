package mysql

import (
	"context"
)

// AdvanceProject 按用户 ID 和项目 ID 流转项目进度
func (r *Repository) AdvanceProject(ctx context.Context, userId, projectId uint64, fromProgress, toProgress ProjectProgress, status ProjectStatus) (bool, error) {
	result, err := r.mysql.ExecContext(ctx, `
		UPDATE projects
		SET progress = ?, status = ?
		WHERE id = ? AND user_id = ? AND progress = ? AND status = ?
	`, toProgress, status, projectId, userId, fromProgress, ProjectStatusActive)
	if err != nil {
		return false, err
	}

	affected, err := result.RowsAffected()
	if err != nil {
		return false, err
	}

	return affected > 0, nil
}

// CreateProject 按用户 ID 和项目名称创建项目，并返回最新记录
func (r *Repository) CreateProject(ctx context.Context, userId uint64, name string) (*ProjectRecord, error) {
	result, err := r.mysql.ExecContext(ctx, `
		INSERT INTO projects(user_id, name, building_name, building_location, building_description, known_issues, assessment_goal)
		VALUES (?, ?, '', '', '', '', '')
	`, userId, name)
	if err != nil {
		return nil, err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return nil, err
	}

	return r.FindProjectByIdAndUserId(ctx, userId, uint64(id))
}

// DeleteProject 按用户 ID 和项目 ID 软删除项目（非物理删除）
func (r *Repository) DeleteProject(ctx context.Context, userId, projectId uint64) (bool, error) {
	result, err := r.mysql.ExecContext(ctx, `
		UPDATE projects
		SET status = ?
		WHERE id = ? AND user_id = ? AND status != ?
	`, ProjectStatusDeleted, projectId, userId, ProjectStatusDeleted)
	if err != nil {
		return false, err
	}

	affected, err := result.RowsAffected()
	if err != nil {
		return false, err
	}

	return affected > 0, nil
}

// ListProjects 按用户 ID 查询项目列表
func (r *Repository) ListProjects(ctx context.Context, userId uint64) ([]*ProjectRecord, []*ProjectRecord, error) {
	rows, err := r.mysql.QueryContext(ctx, `
		SELECT id, user_id, name, building_name, building_location, progress, status, created_at, updated_at
		FROM projects
		WHERE user_id = ? AND status IN (?, ?)
		ORDER BY created_at DESC
	`, userId, ProjectStatusActive, ProjectStatusCompleted)
	if err != nil {
		return nil, nil, err
	}

	defer func() {
		_ = rows.Close()
	}()

	activeProjects := make([]*ProjectRecord, 0)
	completedProjects := make([]*ProjectRecord, 0)

	for rows.Next() {
		project := &ProjectRecord{}
		if err := rows.Scan(
			&project.Id,
			&project.UserId,
			&project.Name,
			&project.BuildingName,
			&project.BuildingLocation,
			&project.Progress,
			&project.Status,
			&project.CreatedAt,
			&project.UpdatedAt,
		); err != nil {
			return nil, nil, err
		}
		if project.Status == string(ProjectStatusActive) {
			activeProjects = append(activeProjects, project)
		} else if project.Status == string(ProjectStatusCompleted) {
			completedProjects = append(completedProjects, project)
		} else {
			continue
		}
	}
	if err := rows.Err(); err != nil {
		return nil, nil, err
	}

	return activeProjects, completedProjects, nil
}
