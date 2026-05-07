package mysql

import (
	"context"
	"database/sql"
	"errors"
	"icw_core_biz/repositories/mysql/model"

	"icw_common/enum"
	"icw_common/gen/core/biz"
)

// FindProjectByIdAndUserId 按用户 ID 和项目 ID 查询项目
func (r *Repository) FindProjectByIdAndUserId(ctx context.Context, userId, projectId uint64) (*model.ProjectRecord, error) {
	project := &model.ProjectRecord{}

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

	project.Progress = enum.ParseProjectProgress(progress)
	project.Status = enum.ParseProjectStatus(status)

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
) (*model.ProjectRecord, error) {
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
	`, name, buildingName, buildingLocation, builtYearValue, buildingDescription, knownIssues, assessmentGoal, projectId, userId, enum.ProjectProgressUint8(bizpb.ProjectProgress_InitializationFinished), enum.ProjectStatusString(bizpb.ProjectStatus_Active))
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
		if project.Progress != bizpb.ProjectProgress_InitializationFinished || project.Status != bizpb.ProjectStatus_Active {
			return nil, nil
		}
		return project, nil
	}

	return r.FindProjectByIdAndUserId(ctx, userId, projectId)
}
