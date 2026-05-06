package mysql

import (
	"context"
	"database/sql"

	"icw_common/enum"
	"icw_common/gen/core/biz"
	projectConsts "icw_core_biz/internal/services/project/consts"
)

type advanceProjectTxFunc func(ctx context.Context, tx *sql.Tx, userId, projectId uint64) error

// AdvanceProject 按用户 ID 和项目 ID 流转项目进度
func (r *Repository) AdvanceProject(ctx context.Context, userId, projectId uint64, fromProgress, toProgress bizpb.ProjectProgress_Value, status bizpb.ProjectStatus_Value, fn advanceProjectTxFunc) (bool, error) {
	tx, err := r.mysql.BeginTx(ctx, nil)
	if err != nil {
		return false, err
	}
	defer func() {
		_ = tx.Rollback()
	}()

	result, err := tx.ExecContext(ctx, `
		UPDATE projects
		SET progress = ?, status = ?
		WHERE id = ? AND user_id = ? AND progress = ? AND status = ?
	`, enum.ProjectProgressUint8(toProgress), enum.ProjectStatusString(status), projectId, userId, enum.ProjectProgressUint8(fromProgress), enum.ProjectStatusString(bizpb.ProjectStatus_Active))
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

	if fn != nil {
		// 通过 MySQL 事务执行项目进度流转后置扩展点
		if err := fn(ctx, tx, userId, projectId); err != nil {
			return false, err
		}
	}

	if err := tx.Commit(); err != nil {
		return false, err
	}

	return true, nil
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
	`, enum.ProjectStatusString(bizpb.ProjectStatus_Deleted), projectId, userId, enum.ProjectStatusString(bizpb.ProjectStatus_Deleted))
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
		WHERE user_id = ? AND status IN (?, ?)
		ORDER BY created_at DESC
	`, userId, enum.ProjectStatusString(bizpb.ProjectStatus_Active), enum.ProjectStatusString(bizpb.ProjectStatus_Completed))
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
		var (
			progress uint8
			status   string
		)
		if err := rows.Scan(
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
		); err != nil {
			return nil, nil, err
		}
		project.Progress = enum.ParseProjectProgress(progress)
		project.Status = enum.ParseProjectStatus(status)
		if project.Status == bizpb.ProjectStatus_Active {
			activeProjects = append(activeProjects, project)
		} else if project.Status == bizpb.ProjectStatus_Completed {
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

// PostAdvanceProjectProfileToAssets 项目进度流转后置扩展点：项目基础信息阶段 -> 图像资产构建阶段
func PostAdvanceProjectProfileToAssets(ctx context.Context, tx *sql.Tx, userId, projectId uint64) error {
	_, err := tx.ExecContext(ctx, `
		INSERT INTO project_groups(project_id, user_id, name, sort_order)
		SELECT ?, ?, ?, COALESCE(MAX(sort_order), -1) + 1
		FROM project_groups
		WHERE user_id = ? AND project_id = ?
	`, projectId, userId, projectConsts.DefaultProjectGroupName, userId, projectId)

	if IsDuplicateEntryError(err) {
		return nil
	}

	return err
}
