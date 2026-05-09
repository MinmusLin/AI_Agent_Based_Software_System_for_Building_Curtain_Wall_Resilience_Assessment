package project_detection

import (
	"icw_common/enum"
	"icw_common/gen/core/common"

	"icw_core_biz/repositories/mysql/model"
)

// projectDetectionSubTaskSchema 项目图像检测子任务表结构映射
type projectDetectionSubTaskSchema struct {
	shouldExecuteColumn string
	table               string
	taskIdColumn        string
}

// projectDetectionSubTaskSchemaByCode 按子任务代码获取子任务表结构映射
func projectDetectionSubTaskSchemaByCode(taskCode string) (*projectDetectionSubTaskSchema, error) {
	switch enum.ParseDetectionTaskCode(taskCode) {
	case commonpb.DetectionTaskCode_Corrosion:
		return &projectDetectionSubTaskSchema{
			shouldExecuteColumn: "corrosion_should_execute",
			table:               "project_detection_corrosion_tasks",
			taskIdColumn:        "corrosion_task_id",
		}, nil
	case commonpb.DetectionTaskCode_Crack:
		return &projectDetectionSubTaskSchema{
			shouldExecuteColumn: "crack_should_execute",
			table:               "project_detection_crack_tasks",
			taskIdColumn:        "crack_task_id",
		}, nil
	case commonpb.DetectionTaskCode_Stain:
		return &projectDetectionSubTaskSchema{
			shouldExecuteColumn: "stain_should_execute",
			table:               "project_detection_stain_tasks",
			taskIdColumn:        "stain_task_id",
		}, nil
	case commonpb.DetectionTaskCode_Flatness:
		return &projectDetectionSubTaskSchema{
			shouldExecuteColumn: "flatness_should_execute",
			table:               "project_detection_flatness_tasks",
			taskIdColumn:        "flatness_task_id",
		}, nil
	case commonpb.DetectionTaskCode_Spalling:
		return &projectDetectionSubTaskSchema{
			shouldExecuteColumn: "spalling_should_execute",
			table:               "project_detection_spalling_tasks",
			taskIdColumn:        "spalling_task_id",
		}, nil
	default:
		return nil, model.ErrUnsupportedDetectionTaskCode
	}
}
