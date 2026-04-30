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

func NewLoginResponse(resp *bizDto.LoginResponse) *LoginResponse {
	if resp == nil {
		return &LoginResponse{}
	}
	return &LoginResponse{
		AccessToken: resp.AccessToken,
		ExpiresIn:   resp.AccessTokenExpiresIn,
		User:        NewUser(resp.User),
	}
}

type LogoutResponse struct{}

func NewLogoutResponse(_ *bizDto.LogoutResponse) *LogoutResponse {
	return &LogoutResponse{}
}

type MeResponse struct {
	User *User `json:"user"`
}

func NewMeResponse(resp *bizDto.MeResponse) *MeResponse {
	if resp == nil {
		return &MeResponse{}
	}
	return &MeResponse{
		User: NewUser(resp.User),
	}
}

type RefreshResponse struct {
	AccessToken string `json:"access_token"`
	ExpiresIn   int    `json:"expires_in"`
	User        *User  `json:"user"`
}

func NewRefreshResponse(resp *bizDto.RefreshResponse) *RefreshResponse {
	if resp == nil {
		return &RefreshResponse{}
	}
	return &RefreshResponse{
		AccessToken: resp.AccessToken,
		ExpiresIn:   resp.AccessTokenExpiresIn,
		User:        NewUser(resp.User),
	}
}

type RegisterRequest struct {
	Email     string `json:"email"`
	EmailCode string `json:"email_code"`
	Password  string `json:"password"`
	Name      string `json:"name"`
}

type RegisterResponse struct{}

func NewRegisterResponse(_ *bizDto.RegisterResponse) *RegisterResponse {
	return &RegisterResponse{}
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

func NewSendEmailCodeResponse(resp *bizDto.SendEmailCodeResponse) *SendEmailCodeResponse {
	if resp == nil {
		return &SendEmailCodeResponse{}
	}
	return &SendEmailCodeResponse{
		ExpiresIn: resp.ExpiresInSeconds,
	}
}
