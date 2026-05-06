package auth

import (
	"context"

	"icw_common/gen/core/biz"
	"icw_core_api/rpc/common"
	"icw_core_api/rpc/icw_core_biz"
)

// Login 登录
func Login(ctx context.Context, client *icw_core_biz.Client, req *bizpb.LoginRequest, resp *bizpb.LoginResponse) error {
	return common.CallGRPC[bizpb.LoginRequest, bizpb.LoginResponse](ctx, client, req, resp, client.Auth().Login)
}

// Logout 登出
func Logout(ctx context.Context, client *icw_core_biz.Client, req *bizpb.LogoutRequest, resp *bizpb.LogoutResponse) error {
	return common.CallGRPC[bizpb.LogoutRequest, bizpb.LogoutResponse](ctx, client, req, resp, client.Auth().Logout)
}

// Me 获取用户信息
func Me(ctx context.Context, client *icw_core_biz.Client, req *bizpb.MeRequest, resp *bizpb.MeResponse) error {
	return common.CallGRPC[bizpb.MeRequest, bizpb.MeResponse](ctx, client, req, resp, client.Auth().Me)
}

// Refresh 刷新 Token
func Refresh(ctx context.Context, client *icw_core_biz.Client, req *bizpb.RefreshRequest, resp *bizpb.RefreshResponse) error {
	return common.CallGRPC[bizpb.RefreshRequest, bizpb.RefreshResponse](ctx, client, req, resp, client.Auth().Refresh)
}

// Register 注册
func Register(ctx context.Context, client *icw_core_biz.Client, req *bizpb.RegisterRequest, resp *bizpb.RegisterResponse) error {
	return common.CallGRPC[bizpb.RegisterRequest, bizpb.RegisterResponse](ctx, client, req, resp, client.Auth().Register)
}

// ResetPassword 重置密码
func ResetPassword(ctx context.Context, client *icw_core_biz.Client, req *bizpb.ResetPasswordRequest, resp *bizpb.ResetPasswordResponse) error {
	return common.CallGRPC[bizpb.ResetPasswordRequest, bizpb.ResetPasswordResponse](ctx, client, req, resp, client.Auth().ResetPassword)
}

// SendEmailCode 发送邮箱验证码
func SendEmailCode(ctx context.Context, client *icw_core_biz.Client, req *bizpb.SendEmailCodeRequest, resp *bizpb.SendEmailCodeResponse) error {
	return common.CallGRPC[bizpb.SendEmailCodeRequest, bizpb.SendEmailCodeResponse](ctx, client, req, resp, client.Auth().SendEmailCode)
}
