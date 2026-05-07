package project_detection

import (
	"database/sql"
	"encoding/json"

	"icw_common/enum"

	"icw_core_biz/repositories/mysql/model"
)

func jsonOrEmptyArray(raw json.RawMessage) string {
	if len(raw) == 0 {
		return "[]"
	}
	return string(raw)
}

func CheckRowsAffected(result sql.Result, err error) error {
	if err != nil {
		return err
	}
	affected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if affected == 0 {
		return sql.ErrNoRows
	}
	return nil
}

type sqlScanner interface {
	Scan(dest ...interface{}) error
}

func ScanProjectDetectionTask(scanner sqlScanner) (*model.ProjectDetectionTaskRecord, error) {
	task := &model.ProjectDetectionTaskRecord{}
	var status string
	if err := scanner.Scan(
		&task.Id,
		&task.Uuid,
		&task.UserId,
		&task.ProjectId,
		&task.ImageId,
		&task.ImageUuid,
		&status,
		&task.CorrosionShouldExecute,
		&task.CorrosionTaskId,
		&task.CrackShouldExecute,
		&task.CrackTaskId,
		&task.StainShouldExecute,
		&task.StainTaskId,
		&task.FlatnessShouldExecute,
		&task.FlatnessTaskId,
		&task.SpallingShouldExecute,
		&task.SpallingTaskId,
		&task.SummaryShouldExecute,
		&task.SummaryTaskId,
		&task.StartedAt,
		&task.FinishedAt,
		&task.CreatedAt,
		&task.UpdatedAt,
	); err != nil {
		return nil, err
	}
	task.Status = enum.ParseProjectDetectionTaskStatus(status)
	return task, nil
}

func scanProjectDetectionSubTask(scanner sqlScanner) (*model.ProjectDetectionSubTaskRecord, error) {
	record := &model.ProjectDetectionSubTaskRecord{}
	var status string
	if err := scanner.Scan(
		&record.Id,
		&record.Uuid,
		&record.MainTaskId,
		&record.UserId,
		&record.ProjectId,
		&record.ImageId,
		&status,
		&record.StartedAt,
		&record.FinishedAt,
		&record.CreatedAt,
		&record.UpdatedAt,
	); err != nil {
		return nil, err
	}
	record.Status = enum.ParseProjectDetectionSubTaskStatus(status)
	return record, nil
}
