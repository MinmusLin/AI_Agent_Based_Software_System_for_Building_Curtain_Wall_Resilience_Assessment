package dto

import (
	"icw_core_biz/repositories"
)

type User struct {
	Id    uint64
	Email string
	Name  string
}

// UserRecordToDTO 将 MySQL 数据模型转换为 RPC 数据模型
func UserRecordToDTO(user *repositories.UserRecord) *User {
	if user == nil {
		return nil
	}
	return &User{
		Id:    user.Id,
		Email: user.Email,
		Name:  user.Name,
	}
}
