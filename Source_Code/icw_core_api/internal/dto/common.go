package dto

import (
	"icw_core_biz/pkg/dto"
)

type User struct {
	Id    uint64 `json:"id"`
	Email string `json:"email"`
	Name  string `json:"name"`
}

func NewUser(user *dto.User) *User {
	if user == nil {
		return &User{}
	}
	return &User{
		Id:    user.Id,
		Email: user.Email,
		Name:  user.Name,
	}
}
