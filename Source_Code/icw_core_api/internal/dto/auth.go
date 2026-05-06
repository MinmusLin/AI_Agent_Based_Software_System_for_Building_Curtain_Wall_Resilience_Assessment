package dto

import (
	"icw_common/gen/core/api"
	"icw_common/gen/core/biz"
)

func NewLoginResponse(resp *bizpb.LoginResponse) *apipb.LoginResponse {
	if resp == nil {
		return nil
	}
	return &apipb.LoginResponse{
		AccessToken: resp.AccessToken,
		ExpiresIn:   resp.AccessTokenExpiresIn,
		User:        resp.User,
	}
}

func NewLogoutResponse(_ *bizpb.LogoutResponse) *apipb.LogoutResponse {
	return &apipb.LogoutResponse{}
}

func NewMeResponse(resp *bizpb.MeResponse) *apipb.MeResponse {
	if resp == nil {
		return nil
	}
	return &apipb.MeResponse{
		User: resp.User,
	}
}

func NewRefreshResponse(resp *bizpb.RefreshResponse) *apipb.RefreshResponse {
	if resp == nil {
		return nil
	}
	return &apipb.RefreshResponse{
		AccessToken: resp.AccessToken,
		ExpiresIn:   resp.AccessTokenExpiresIn,
		User:        resp.User,
	}
}

func NewRegisterResponse(_ *bizpb.RegisterResponse) *apipb.RegisterResponse {
	return &apipb.RegisterResponse{}
}

func NewResetPasswordResponse(_ *bizpb.ResetPasswordResponse) *apipb.ResetPasswordResponse {
	return &apipb.ResetPasswordResponse{}
}

func NewSendEmailCodeResponse(resp *bizpb.SendEmailCodeResponse) *apipb.SendEmailCodeResponse {
	if resp == nil {
		return nil
	}
	return &apipb.SendEmailCodeResponse{
		ExpiresIn: resp.ExpiresInSeconds,
	}
}
