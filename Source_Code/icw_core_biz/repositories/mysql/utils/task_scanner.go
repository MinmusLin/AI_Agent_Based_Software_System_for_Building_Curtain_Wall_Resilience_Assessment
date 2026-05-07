package utils

import (
	"icw_common/enum"

	"icw_core_biz/repositories/mysql/model"
)

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

func ScanProjectDetectionSubTask(scanner sqlScanner) (*model.ProjectDetectionSubTaskRecord, error) {
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
