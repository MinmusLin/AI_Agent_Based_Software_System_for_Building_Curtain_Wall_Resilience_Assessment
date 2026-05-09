package mysql

import (
	"context"
	"database/sql"
	"strings"

	"icw_common/enum"
	"icw_common/gen/core/biz"

	"icw_core_biz/repositories/mysql/model"
	"icw_core_biz/repositories/mysql/utils"
)

// GetProjectDetectionReview 按用户 ID、项目 ID 和主任务 UUID 查询图像检测人工复核信息
func (r *Repository) GetProjectDetectionReview(ctx context.Context, userId, projectId uint64, taskUuid string) (*bizpb.ProjectDetectionReview, error) {
	review := &bizpb.ProjectDetectionReview{}

	var (
		verdict   sql.NullString
		comment   sql.NullString
		updatedAt sql.NullTime
	)

	err := r.mysql.QueryRowContext(ctx, `
		SELECT uuid, image_uuid, review_verdict, review_comment, updated_at
		FROM project_detection_tasks
		WHERE user_id = ? AND project_id = ? AND uuid = ?
	`, userId, projectId, taskUuid).Scan(&review.TaskUuid, &review.ImageUuid, &verdict, &comment, &updatedAt)
	if err != nil {
		return nil, err
	}

	review.Verdict = enum.ParseProjectDetectionReviewVerdict(utils.NullString(verdict))
	review.Comment = utils.NullString(comment)
	review.UpdatedAt = utils.NullTimeString(updatedAt)

	return review, nil
}

// UpdateProjectDetectionReview 按用户 ID、项目 ID 和主任务 UUID 更新图像检测人工复核信息
func (r *Repository) UpdateProjectDetectionReview(ctx context.Context, userId, projectId uint64, taskUuid string, verdict bizpb.ProjectDetectionReviewVerdict_Value, comment string) (*bizpb.ProjectDetectionReview, error) {
	verdictText := enum.ProjectDetectionReviewVerdictString(verdict)
	if verdictText == "" {
		return nil, model.ErrProjectDetectionReviewVerdictInvalid
	}

	var verdictValue interface{}
	if verdict == bizpb.ProjectDetectionReviewVerdict_Unknown {
		verdictValue = nil
	} else {
		verdictValue = verdictText
	}

	result, err := r.mysql.ExecContext(ctx, `
		UPDATE project_detection_tasks
		SET review_verdict = ?, review_comment = ?
		WHERE user_id = ? AND project_id = ? AND uuid = ?
	`, verdictValue, strings.TrimSpace(comment), userId, projectId, taskUuid)
	if err := utils.CheckRowsAffected(result, err); err != nil {
		return nil, err
	}

	return r.GetProjectDetectionReview(ctx, userId, projectId, taskUuid)
}
