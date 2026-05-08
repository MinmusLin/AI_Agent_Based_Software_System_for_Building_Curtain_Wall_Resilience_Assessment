package model

import (
	"icw_common/enum"
)

// sqlScanner 抽象 QueryRow 和 Rows 的 Scan 能力
type sqlScanner interface {
	Scan(dest ...interface{}) error
}

// ScanProjectDetectionTask 扫描项目图像检测主任务记录
func ScanProjectDetectionTask(scanner sqlScanner) (*ProjectDetectionTaskRecord, error) {
	task := &ProjectDetectionTaskRecord{}
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

// ScanProjectDetectionSubTask 扫描项目图像检测子任务记录
func ScanProjectDetectionSubTask(scanner sqlScanner) (*ProjectDetectionSubTaskRecord, error) {
	record := &ProjectDetectionSubTaskRecord{}
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
