package project_detection

import (
	"context"
	"database/sql"
	"encoding/json"

	"icw_common/gen/core/biz"
	"icw_common/gen/core/common"

	"icw_core_biz/repositories/mysql/model"
	"icw_core_biz/repositories/mysql/utils"
)

// updateProjectDetectionCorrosionTaskResultTx 更新项目图像金属锈蚀检测子任务报告与状态
func updateProjectDetectionCorrosionTaskResultTx(ctx context.Context, tx *sql.Tx, taskUuid string, status bizpb.ProjectDetectionSubTaskStatus_Value, resultJSON, artifactSha256Map string) error {
	if status != bizpb.ProjectDetectionSubTaskStatus_Succeeded {
		return updateProjectDetectionSubTaskStatusByTableTx(ctx, tx, "project_detection_corrosion_tasks", taskUuid, status, false, true, model.ErrProjectDetectionSubTaskStatusInvalid)
	}

	report := &commonpb.ProjectDetectionCorrosionResult{}
	if err := json.Unmarshal([]byte(resultJSON), report); err != nil {
		return err
	}
	regionsJSON, err := utils.JsonOrEmptyArray(report.Regions)
	if err != nil {
		return err
	}

	result, err := tx.ExecContext(ctx, `
		UPDATE project_detection_corrosion_tasks
		SET has_corrosion = ?,
			corrosion_count = ?,
			max_confidence = ?,
			average_confidence = ?,
			corrosion_pixels = ?,
			corrosion_ratio = ?,
			regions = ?,
			artifact_sha256_map = ?,
			runtime_seconds = ?
		WHERE uuid = ?
	`, report.HasCorrosion, report.CorrosionCount, report.MaxConfidence, report.AverageConfidence, report.CorrosionPixels, report.CorrosionRatio, regionsJSON, utils.JsonStringOrEmptyObject(artifactSha256Map), report.RuntimeSeconds, taskUuid)
	if err := utils.CheckRowsAffected(result, err); err != nil {
		return err
	}

	return updateProjectDetectionSubTaskStatusByTableTx(ctx, tx, "project_detection_corrosion_tasks", taskUuid, status, false, true, model.ErrProjectDetectionSubTaskStatusInvalid)
}

// updateProjectDetectionCrackTaskResultTx 更新项目图像石材裂缝检测子任务报告与状态
func updateProjectDetectionCrackTaskResultTx(ctx context.Context, tx *sql.Tx, taskUuid string, status bizpb.ProjectDetectionSubTaskStatus_Value, resultJSON, artifactSha256Map string) error {
	if status != bizpb.ProjectDetectionSubTaskStatus_Succeeded {
		return updateProjectDetectionSubTaskStatusByTableTx(ctx, tx, "project_detection_crack_tasks", taskUuid, status, false, true, model.ErrProjectDetectionSubTaskStatusInvalid)
	}

	report := &commonpb.ProjectDetectionCrackResult{}
	if err := json.Unmarshal([]byte(resultJSON), report); err != nil {
		return err
	}
	regionsJSON, err := utils.JsonOrEmptyArray(report.Regions)
	if err != nil {
		return err
	}

	result, err := tx.ExecContext(ctx, `
		UPDATE project_detection_crack_tasks
		SET has_crack = ?,
			crack_count = ?,
			crack_pixels = ?,
			crack_ratio = ?,
			regions = ?,
			artifact_sha256_map = ?,
			runtime_seconds = ?
		WHERE uuid = ?
	`, report.HasCrack, report.CrackCount, report.CrackPixels, report.CrackRatio, regionsJSON, utils.JsonStringOrEmptyObject(artifactSha256Map), report.RuntimeSeconds, taskUuid)
	if err := utils.CheckRowsAffected(result, err); err != nil {
		return err
	}

	return updateProjectDetectionSubTaskStatusByTableTx(ctx, tx, "project_detection_crack_tasks", taskUuid, status, false, true, model.ErrProjectDetectionSubTaskStatusInvalid)
}

// updateProjectDetectionStainTaskResultTx 更新项目图像石材污渍检测子任务报告与状态
func updateProjectDetectionStainTaskResultTx(ctx context.Context, tx *sql.Tx, taskUuid string, status bizpb.ProjectDetectionSubTaskStatus_Value, resultJSON, artifactSha256Map string) error {
	if status != bizpb.ProjectDetectionSubTaskStatus_Succeeded {
		return updateProjectDetectionSubTaskStatusByTableTx(ctx, tx, "project_detection_stain_tasks", taskUuid, status, false, true, model.ErrProjectDetectionSubTaskStatusInvalid)
	}

	report := &commonpb.ProjectDetectionStainResult{}
	if err := json.Unmarshal([]byte(resultJSON), report); err != nil {
		return err
	}
	regionsJSON, err := utils.JsonOrEmptyArray(report.Regions)
	if err != nil {
		return err
	}

	result, err := tx.ExecContext(ctx, `
		UPDATE project_detection_stain_tasks
		SET has_stain = ?,
			stain_count = ?,
			average_stain_ratio = ?,
			max_stain_ratio = ?,
			regions = ?,
			artifact_sha256_map = ?,
			runtime_seconds = ?
		WHERE uuid = ?
	`, report.HasStain, report.StainCount, report.AverageStainRatio, report.MaxStainRatio, regionsJSON, utils.JsonStringOrEmptyObject(artifactSha256Map), report.RuntimeSeconds, taskUuid)
	if err := utils.CheckRowsAffected(result, err); err != nil {
		return err
	}

	return updateProjectDetectionSubTaskStatusByTableTx(ctx, tx, "project_detection_stain_tasks", taskUuid, status, false, true, model.ErrProjectDetectionSubTaskStatusInvalid)
}

// updateProjectDetectionFlatnessTaskResultTx 更新项目图像玻璃平整度检测子任务报告与状态
func updateProjectDetectionFlatnessTaskResultTx(ctx context.Context, tx *sql.Tx, taskUuid string, status bizpb.ProjectDetectionSubTaskStatus_Value, resultJSON, artifactSha256Map string) error {
	if status != bizpb.ProjectDetectionSubTaskStatus_Succeeded {
		return updateProjectDetectionSubTaskStatusByTableTx(ctx, tx, "project_detection_flatness_tasks", taskUuid, status, false, true, model.ErrProjectDetectionSubTaskStatusInvalid)
	}

	report := &commonpb.ProjectDetectionFlatnessResult{}
	if err := json.Unmarshal([]byte(resultJSON), report); err != nil {
		return err
	}
	regionsJSON, err := utils.JsonOrEmptyArray(report.Regions)
	if err != nil {
		return err
	}

	result, err := tx.ExecContext(ctx, `
		UPDATE project_detection_flatness_tasks
		SET result = ?,
			uneven_count = ?,
			regions = ?,
			artifact_sha256_map = ?,
			runtime_seconds = ?
		WHERE uuid = ?
	`, report.Result, report.UnevenCount, regionsJSON, utils.JsonStringOrEmptyObject(artifactSha256Map), report.RuntimeSeconds, taskUuid)
	if err := utils.CheckRowsAffected(result, err); err != nil {
		return err
	}

	return updateProjectDetectionSubTaskStatusByTableTx(ctx, tx, "project_detection_flatness_tasks", taskUuid, status, false, true, model.ErrProjectDetectionSubTaskStatusInvalid)
}

// updateProjectDetectionSpallingTaskResultTx 更新项目图像玻璃爆裂检测子任务报告与状态
func updateProjectDetectionSpallingTaskResultTx(ctx context.Context, tx *sql.Tx, taskUuid string, status bizpb.ProjectDetectionSubTaskStatus_Value, resultJSON, artifactSha256Map string) error {
	if status != bizpb.ProjectDetectionSubTaskStatus_Succeeded {
		return updateProjectDetectionSubTaskStatusByTableTx(ctx, tx, "project_detection_spalling_tasks", taskUuid, status, false, true, model.ErrProjectDetectionSubTaskStatusInvalid)
	}

	report := &commonpb.ProjectDetectionSpallingResult{}
	if err := json.Unmarshal([]byte(resultJSON), report); err != nil {
		return err
	}

	result, err := tx.ExecContext(ctx, `
		UPDATE project_detection_spalling_tasks
		SET has_spalling = ?,
			confidence = ?,
			artifact_sha256_map = ?,
			runtime_seconds = ?
		WHERE uuid = ?
	`, report.HasSpalling, report.Confidence, utils.JsonStringOrEmptyObject(artifactSha256Map), report.RuntimeSeconds, taskUuid)
	if err := utils.CheckRowsAffected(result, err); err != nil {
		return err
	}

	return updateProjectDetectionSubTaskStatusByTableTx(ctx, tx, "project_detection_spalling_tasks", taskUuid, status, false, true, model.ErrProjectDetectionSubTaskStatusInvalid)
}
