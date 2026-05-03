package mysql

import (
	"context"
	"database/sql"
	"errors"

	"icw_core_biz/pkg/dto"
)

// FindProjectByIdAndUserId 按用户 ID 和项目 ID 查询项目
func (r *Repository) FindProjectByIdAndUserId(ctx context.Context, userId, projectId uint64) (*ProjectRecord, error) {
	project := &ProjectRecord{}

	var (
		progress uint8
		status   string
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
		&project.BuildingDescription,
		&project.KnownIssues,
		&project.AssessmentGoal,
		&progress,
		&status,
		&project.CreatedAt,
		&project.UpdatedAt,
	)

	if errors.Is(err, sql.ErrNoRows) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	project.Progress = dto.ParseProjectProgress(progress)
	project.Status = dto.ParseProjectStatus(status)

	return project, nil
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
			assessment_goal = ?
		WHERE id = ? AND user_id = ? AND progress = ? AND status = ?
	`, name, buildingName, buildingLocation, builtYearValue, buildingDescription, knownIssues, assessmentGoal, projectId, userId, dto.ProjectProgressInitializationFinished.Uint8(), dto.ProjectStatusActive.String())
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
		if project.Progress != dto.ProjectProgressInitializationFinished || project.Status != dto.ProjectStatusActive {
			return nil, nil
		}
		return project, nil
	}

	return r.FindProjectByIdAndUserId(ctx, userId, projectId)
}
