package mysql

import (
	"context"
	"database/sql"
	"errors"

	"icw_core_biz/internal/services/project/consts"
)

// FindProjectByIdAndUserId 按用户 ID 和项目 ID 查询项目
func (r *Repository) FindProjectByIdAndUserId(ctx context.Context, userId, projectId uint64) (*ProjectRecord, error) {
	project := &ProjectRecord{}

	var (
		buildingDescription sql.NullString
		knownIssues         sql.NullString
		assessmentGoal      sql.NullString
	)

	err := r.mysql.QueryRowContext(ctx, `
		SELECT
			id, 
			user_id, 
			name, 
			building_name, 
			building_location, 
			built_year, 
			building_description,
			known_issues, 
			assessment_goal, 
			progress, 
			status, 
			created_at, 
			updated_at
		FROM projects
		WHERE id = ? AND user_id = ?
		LIMIT 1
	`, projectId, userId).Scan(
		&project.Id,
		&project.UserId,
		&project.Name,
		&project.BuildingName,
		&project.BuildingLocation,
		&project.BuiltYear,
		&buildingDescription,
		&knownIssues,
		&assessmentGoal,
		&project.Progress,
		&project.Status,
		&project.CreatedAt,
		&project.UpdatedAt,
	)

	if errors.Is(err, sql.ErrNoRows) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	project.BuildingDescription = buildingDescription.String
	project.KnownIssues = knownIssues.String
	project.AssessmentGoal = assessmentGoal.String

	return project, nil
}

// AdvanceProject 按用户 ID 和项目 ID 流转项目进度
func (r *Repository) AdvanceProject(ctx context.Context, userId, projectId uint64, fromProgress, toProgress uint8, status ProjectStatus) (bool, error) {
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

// UpdateProjectProfile 更新项目基础信息，并返回最新记录
func (r *Repository) UpdateProjectProfile(
	ctx context.Context,
	userId uint64,
	projectId uint64,
	name string,
	buildingName string,
	buildingLocation string,
	builtYear uint16,
	buildingDescription string,
	knownIssues string,
	assessmentGoal string,
) (*ProjectRecord, error) {
	var builtYearValue sql.NullInt64
	if builtYear > 0 {
		builtYearValue = sql.NullInt64{
			Int64: int64(builtYear),
			Valid: true,
		}
	}

	result, err := r.mysql.ExecContext(ctx, `
		UPDATE projects
		SET name = ?,
			building_name = ?,
			building_location = ?,
			built_year = ?,
			building_description = ?,
			known_issues = ?,
			assessment_goal = ?,
			updated_at = NOW(3)
		WHERE id = ? AND user_id = ? AND progress = ? AND status = ?
	`, name, buildingName, buildingLocation, builtYearValue, buildingDescription, knownIssues, assessmentGoal, projectId, userId, consts.ProgressInitializationFinished, ProjectStatusActive)
	if err != nil {
		return nil, err
	}

	affected, err := result.RowsAffected()
	if err != nil {
		return nil, err
	}

	if affected == 0 {
		project, err := r.FindProjectByIdAndUserId(ctx, userId, projectId)
		if err != nil || project == nil {
			return project, err
		}
		if project.Progress != consts.ProgressInitializationFinished || project.Status != string(ProjectStatusActive) {
			return nil, nil
		}
		return project, nil
	}

	return r.FindProjectByIdAndUserId(ctx, userId, projectId)
}
