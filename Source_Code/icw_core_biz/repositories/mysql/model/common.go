package model

import (
	"database/sql"
	"time"

	"icw_common/gen/core/biz"
)

// UserRecord 用户记录
type UserRecord struct {
	Id           uint64
	Email        string
	PasswordHash string
	Name         string
	LastLoginAt  sql.NullTime
	CreatedAt    time.Time
	UpdatedAt    time.Time
}

// UserRecordToDTO 将 MySQL 数据模型转换为 RPC 数据模型
func UserRecordToDTO(user *UserRecord) *bizpb.User {
	if user == nil {
		return nil
	}
	return &bizpb.User{
		Id:    user.Id,
		Email: user.Email,
		Name:  user.Name,
	}
}
