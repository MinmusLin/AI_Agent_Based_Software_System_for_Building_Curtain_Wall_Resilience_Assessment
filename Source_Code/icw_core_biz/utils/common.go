package utils

import (
	"encoding/json"

	"icw_core_biz/pkg/dto"
	"icw_core_biz/repositories"
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
func UserRecordToDTO(user *repositories.UserRecord) *dto.User {
	if user == nil {
		return nil
	}
	return &dto.User{
		Id:    user.Id,
		Email: user.Email,
		Name:  user.Name,
	}
}
