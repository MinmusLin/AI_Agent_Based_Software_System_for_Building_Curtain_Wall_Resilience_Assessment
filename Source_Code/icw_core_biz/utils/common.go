package utils

import (
	"encoding/json"

	"icw_core_biz/pkg/dto"
	"icw_core_biz/pkg/dto/project"
	"icw_core_biz/repositories/mysql"
)

// JSONF 将任意结构格式化为 JSON 字符串
func JSONF(v interface{}) string {
	bytes, err := json.Marshal(v)
	if err != nil {
		return ""
	}
	return string(bytes)
}

// UserRecordToDTO 将 MySQL 数据模型转换为 RPC 数据模型
func UserRecordToDTO(user *mysql.UserRecord) *dto.User {
	if user == nil {
		return nil
	}
	return &dto.User{
		Id:    user.Id,
		Email: user.Email,
		Name:  user.Name,
	}
}

// ProjectRecordToDTO 将 MySQL 数据模型转换为 RPC 数据模型
func ProjectRecordToDTO(record *mysql.ProjectRecord) *project.Project {
	if record == nil {
		return nil
	}
	return &project.Project{
		Id:                  record.Id,
		Name:                record.Name,
		BuildingName:        record.BuildingName,
		BuildingLocation:    record.BuildingLocation,
		BuiltYear:           mysql.NullInt64ToUint16(record.BuiltYear),
		BuildingDescription: record.BuildingDescription,
		KnownIssues:         record.KnownIssues,
		AssessmentGoal:      record.AssessmentGoal,
		Progress:            record.Progress,
		CreatedAt:           mysql.TimeToString(record.CreatedAt),
		UpdatedAt:           mysql.TimeToString(record.UpdatedAt),
	}
}

// ProjectRecordsToListItemsDTO 将 MySQL 项目模型批量转换为 RPC 项目列表模型
func ProjectRecordsToListItemsDTO(records []*mysql.ProjectRecord) []*project.ProjectListItem {
	if records == nil {
		return make([]*project.ProjectListItem, 0)
	}

	items := make([]*project.ProjectListItem, 0, len(records))
	for _, record := range records {
		if record == nil {
			continue
		}
		items = append(items, &project.ProjectListItem{
			Id:               record.Id,
			Name:             record.Name,
			BuildingName:     record.BuildingName,
			BuildingLocation: record.BuildingLocation,
			Progress:         record.Progress,
			CreatedAt:        mysql.TimeToString(record.CreatedAt),
		})
	}
	return items
}
