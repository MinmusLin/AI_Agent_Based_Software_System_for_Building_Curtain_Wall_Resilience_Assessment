package dto

import (
	"icw_core_biz/pkg/dto"
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

func NewLoginResponse(resp *dto.LoginResponse) *LoginResponse {
	if resp == nil {
		return nil
	}
	return &LoginResponse{
		AccessToken: resp.AccessToken,
		ExpiresIn:   resp.AccessTokenExpiresIn,
		User:        NewUser(resp.User),
	}
}

type LogoutResponse struct{}

func NewLogoutResponse(_ *dto.LogoutResponse) *LogoutResponse {
	return nil
}

type MeResponse struct {
	User *User `json:"user"`
}

func NewMeResponse(resp *dto.MeResponse) *MeResponse {
	if resp == nil {
		return nil
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

func NewRefreshResponse(resp *dto.RefreshResponse) *RefreshResponse {
	if resp == nil {
		return nil
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

func NewRegisterResponse(_ *dto.RegisterResponse) *RegisterResponse {
	return nil
}

type ResetPasswordRequest struct {
	Email       string `json:"email"`
	EmailCode   string `json:"email_code"`
	NewPassword string `json:"new_password"`
}

type ResetPasswordResponse struct{}

func NewResetPasswordResponse(_ *dto.ResetPasswordResponse) *ResetPasswordResponse {
	return nil
}

type SendEmailCodeRequest struct {
	Email string `json:"email"`
	Scene string `json:"scene"`
}

type SendEmailCodeResponse struct {
	ExpiresIn int `json:"expires_in"`
}

func NewSendEmailCodeResponse(resp *dto.SendEmailCodeResponse) *SendEmailCodeResponse {
	if resp == nil {
		return nil
	}
	return &SendEmailCodeResponse{
		ExpiresIn: resp.ExpiresInSeconds,
	}
}
