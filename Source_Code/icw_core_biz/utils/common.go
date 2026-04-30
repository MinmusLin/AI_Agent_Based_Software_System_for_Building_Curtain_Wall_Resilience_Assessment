package utils

import (
	"icw_core_biz/pkg/dto"
	"icw_core_biz/repositories"
)

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
