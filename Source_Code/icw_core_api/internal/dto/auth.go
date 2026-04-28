package dto

import (
	bizDto "icw_core_biz/pkg/dto"
)

type LoginRequest struct {
	Email string `json:"email"`
	Scene string `json:"scene"`
	Code  string `json:"code"`
}

type LoginResponse struct {
	AccessToken string `json:"access_token"`
	ExpiresIn   int    `json:"expires_in"`
	User        *User  `json:"user"`
}

func NewLoginResponse(res *bizDto.LoginResponse) *LoginResponse {
	return &LoginResponse{
		AccessToken: res.AccessToken,
		ExpiresIn:   res.AccessTokenExpiresIn,
		User:        NewUser(res.User),
	}
}

type LogoutRequest struct {
	AccessToken  string `json:"-"`
	RefreshToken string `json:"-"`
}

type LogoutResponse struct{}

func NewLogoutResponse(_ *bizDto.LogoutResponse) *LogoutResponse {
	return &LogoutResponse{}
}

type MeRequest struct {
	AccessToken string `json:"-"`
}

type MeResponse struct {
	User *User `json:"user"`
}

func NewMeResponse(res *bizDto.MeResponse) *MeResponse {
	return &MeResponse{
		User: NewUser(res.User),
	}
}

type RefreshRequest struct {
	RefreshToken string `json:"-"`
}

type RefreshResponse struct {
	AccessToken string `json:"access_token"`
	ExpiresIn   int    `json:"expires_in"`
	User        *User  `json:"user"`
}

func NewRefreshResponse(res *bizDto.RefreshResponse) *RefreshResponse {
	return &RefreshResponse{
		AccessToken: res.AccessToken,
		ExpiresIn:   res.AccessTokenExpiresIn,
		User:        NewUser(res.User),
	}
}

type RegisterRequest struct {
	Email     string `json:"email"`
	EmailCode string `json:"email_code"`
	Password  string `json:"password"`
	Name      string `json:"name"`
}

type RegisterResponse struct {
	User *User `json:"user"`
}

func NewRegisterResponse(res *bizDto.RegisterResponse) *RegisterResponse {
	return &RegisterResponse{
		User: NewUser(res.User),
	}
}

type ResetPasswordRequest struct {
	Email       string `json:"email"`
	EmailCode   string `json:"email_code"`
	NewPassword string `json:"new_password"`
}

type ResetPasswordResponse struct{}

func NewResetPasswordResponse(_ *bizDto.ResetPasswordResponse) *ResetPasswordResponse {
	return &ResetPasswordResponse{}
}

type SendEmailCodeRequest struct {
	Email string `json:"email"`
	Scene string `json:"scene"`
}

type SendEmailCodeResponse struct {
	ExpiresIn int `json:"expires_in"`
}

func NewSendEmailCodeResponse(res *bizDto.SendEmailCodeResponse) *SendEmailCodeResponse {
	return &SendEmailCodeResponse{
		ExpiresIn: res.ExpiresInSeconds,
	}
}
