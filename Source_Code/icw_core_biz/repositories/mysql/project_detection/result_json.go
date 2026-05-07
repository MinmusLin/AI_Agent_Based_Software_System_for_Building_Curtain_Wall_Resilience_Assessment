package project_detection

import (
	"context"
	"database/sql"
	"encoding/json"

	"icw_common/gen/core/biz"

	"icw_core_biz/repositories/mysql/model"
	"icw_core_biz/repositories/mysql/utils"
)

// updateProjectDetectionCorrosionTaskResultTx 更新项目图像金属锈蚀检测子任务报告与状态
func updateProjectDetectionCorrosionTaskResultTx(ctx context.Context, tx *sql.Tx, taskUuid string, status bizpb.ProjectDetectionSubTaskStatus_Value, resultJSON string) error {
	if status != bizpb.ProjectDetectionSubTaskStatus_Succeeded {
		return updateProjectDetectionSubTaskStatusByTableTx(ctx, tx, "project_detection_corrosion_tasks", taskUuid, status, false, true, model.ErrProjectDetectionSubTaskStatusInvalid)
	}

	report := struct {
		HasCorrosion      bool            `json:"has_corrosion"`
		CorrosionCount    uint64          `json:"corrosion_count"`
		MaxConfidence     float64         `json:"max_confidence"`
		AverageConfidence float64         `json:"average_confidence"`
		CorrosionPixels   uint64          `json:"corrosion_pixels"`
		CorrosionRatio    float64         `json:"corrosion_ratio"`
		Regions           json.RawMessage `json:"regions"`
		RuntimeSeconds    float64         `json:"runtime_seconds"`
	}{}
	if err := json.Unmarshal([]byte(resultJSON), &report); err != nil {
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
			runtime_seconds = ?
		WHERE uuid = ?
	`, report.HasCorrosion, report.CorrosionCount, report.MaxConfidence, report.AverageConfidence, report.CorrosionPixels, report.CorrosionRatio, utils.JsonOrEmptyArray(report.Regions), report.RuntimeSeconds, taskUuid)
	if err := utils.CheckRowsAffected(result, err); err != nil {
		return err
	}

	return updateProjectDetectionSubTaskStatusByTableTx(ctx, tx, "project_detection_corrosion_tasks", taskUuid, status, false, true, model.ErrProjectDetectionSubTaskStatusInvalid)
}

// updateProjectDetectionCrackTaskResultTx 更新项目图像石材裂缝检测子任务报告与状态
func updateProjectDetectionCrackTaskResultTx(ctx context.Context, tx *sql.Tx, taskUuid string, status bizpb.ProjectDetectionSubTaskStatus_Value, resultJSON string) error {
	if status != bizpb.ProjectDetectionSubTaskStatus_Succeeded {
		return updateProjectDetectionSubTaskStatusByTableTx(ctx, tx, "project_detection_crack_tasks", taskUuid, status, false, true, model.ErrProjectDetectionSubTaskStatusInvalid)
	}

	report := struct {
		HasCrack       bool            `json:"has_crack"`
		CrackCount     uint64          `json:"crack_count"`
		CrackPixels    uint64          `json:"crack_pixels"`
		CrackRatio     float64         `json:"crack_ratio"`
		Regions        json.RawMessage `json:"regions"`
		RuntimeSeconds float64         `json:"runtime_seconds"`
	}{}
	if err := json.Unmarshal([]byte(resultJSON), &report); err != nil {
		return err
	}

	result, err := tx.ExecContext(ctx, `
		UPDATE project_detection_crack_tasks
		SET has_crack = ?,
			crack_count = ?,
			crack_pixels = ?,
			crack_ratio = ?,
			regions = ?,
			runtime_seconds = ?
		WHERE uuid = ?
	`, report.HasCrack, report.CrackCount, report.CrackPixels, report.CrackRatio, utils.JsonOrEmptyArray(report.Regions), report.RuntimeSeconds, taskUuid)
	if err := utils.CheckRowsAffected(result, err); err != nil {
		return err
	}

	return updateProjectDetectionSubTaskStatusByTableTx(ctx, tx, "project_detection_crack_tasks", taskUuid, status, false, true, model.ErrProjectDetectionSubTaskStatusInvalid)
}

// updateProjectDetectionStainTaskResultTx 更新项目图像石材污渍检测子任务报告与状态
func updateProjectDetectionStainTaskResultTx(ctx context.Context, tx *sql.Tx, taskUuid string, status bizpb.ProjectDetectionSubTaskStatus_Value, resultJSON string) error {
	if status != bizpb.ProjectDetectionSubTaskStatus_Succeeded {
		return updateProjectDetectionSubTaskStatusByTableTx(ctx, tx, "project_detection_stain_tasks", taskUuid, status, false, true, model.ErrProjectDetectionSubTaskStatusInvalid)
	}

	report := struct {
		HasStain          bool            `json:"has_stain"`
		StainCount        uint64          `json:"stain_count"`
		AverageStainRatio float64         `json:"average_stain_ratio"`
		MaxStainRatio     float64         `json:"max_stain_ratio"`
		Regions           json.RawMessage `json:"regions"`
		RuntimeSeconds    float64         `json:"runtime_seconds"`
	}{}
	if err := json.Unmarshal([]byte(resultJSON), &report); err != nil {
		return err
	}

	result, err := tx.ExecContext(ctx, `
		UPDATE project_detection_stain_tasks
		SET has_stain = ?,
			stain_count = ?,
			average_stain_ratio = ?,
			max_stain_ratio = ?,
			regions = ?,
			runtime_seconds = ?
		WHERE uuid = ?
	`, report.HasStain, report.StainCount, report.AverageStainRatio, report.MaxStainRatio, utils.JsonOrEmptyArray(report.Regions), report.RuntimeSeconds, taskUuid)
	if err := utils.CheckRowsAffected(result, err); err != nil {
		return err
	}

	return updateProjectDetectionSubTaskStatusByTableTx(ctx, tx, "project_detection_stain_tasks", taskUuid, status, false, true, model.ErrProjectDetectionSubTaskStatusInvalid)
}

// updateProjectDetectionFlatnessTaskResultTx 更新项目图像玻璃平整度检测子任务报告与状态
func updateProjectDetectionFlatnessTaskResultTx(ctx context.Context, tx *sql.Tx, taskUuid string, status bizpb.ProjectDetectionSubTaskStatus_Value, resultJSON string) error {
	if status != bizpb.ProjectDetectionSubTaskStatus_Succeeded {
		return updateProjectDetectionSubTaskStatusByTableTx(ctx, tx, "project_detection_flatness_tasks", taskUuid, status, false, true, model.ErrProjectDetectionSubTaskStatusInvalid)
	}

	report := struct {
		Result         string          `json:"result"`
		UnevenCount    uint64          `json:"uneven_count"`
		Regions        json.RawMessage `json:"regions"`
		RuntimeSeconds float64         `json:"runtime_seconds"`
	}{}
	if err := json.Unmarshal([]byte(resultJSON), &report); err != nil {
		return err
	}

	result, err := tx.ExecContext(ctx, `
		UPDATE project_detection_flatness_tasks
		SET result = ?,
			uneven_count = ?,
			regions = ?,
			runtime_seconds = ?
		WHERE uuid = ?
	`, report.Result, report.UnevenCount, utils.JsonOrEmptyArray(report.Regions), report.RuntimeSeconds, taskUuid)
	if err := utils.CheckRowsAffected(result, err); err != nil {
		return err
	}

	return updateProjectDetectionSubTaskStatusByTableTx(ctx, tx, "project_detection_flatness_tasks", taskUuid, status, false, true, model.ErrProjectDetectionSubTaskStatusInvalid)
}

// updateProjectDetectionSpallingTaskResultTx 更新项目图像玻璃爆裂检测子任务报告与状态
func updateProjectDetectionSpallingTaskResultTx(ctx context.Context, tx *sql.Tx, taskUuid string, status bizpb.ProjectDetectionSubTaskStatus_Value, resultJSON string) error {
	if status != bizpb.ProjectDetectionSubTaskStatus_Succeeded {
		return updateProjectDetectionSubTaskStatusByTableTx(ctx, tx, "project_detection_spalling_tasks", taskUuid, status, false, true, model.ErrProjectDetectionSubTaskStatusInvalid)
	}

	report := struct {
		HasSpalling    bool    `json:"has_spalling"`
		Confidence     float64 `json:"confidence"`
		RuntimeSeconds float64 `json:"runtime_seconds"`
	}{}
	if err := json.Unmarshal([]byte(resultJSON), &report); err != nil {
		return err
	}

	result, err := tx.ExecContext(ctx, `
		UPDATE project_detection_spalling_tasks
		SET has_spalling = ?,
			confidence = ?,
			runtime_seconds = ?
		WHERE uuid = ?
	`, report.HasSpalling, report.Confidence, report.RuntimeSeconds, taskUuid)
	if err := utils.CheckRowsAffected(result, err); err != nil {
		return err
	}

	return updateProjectDetectionSubTaskStatusByTableTx(ctx, tx, "project_detection_spalling_tasks", taskUuid, status, false, true, model.ErrProjectDetectionSubTaskStatusInvalid)
}
