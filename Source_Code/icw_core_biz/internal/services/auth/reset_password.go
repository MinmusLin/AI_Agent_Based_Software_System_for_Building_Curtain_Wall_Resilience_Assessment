package auth

import (
	"context"
	"errors"

	"golang.org/x/crypto/bcrypt"

	"icw_common/enum"
	"icw_common/gen/core/biz"
	"icw_common/gen/core/common"
	"icw_common/rpc/error"

	"icw_core_biz/internal/services/auth/utils"
	"icw_core_biz/internal/services/common"
)

// ResetPassword 重置密码
func (s *Service) ResetPassword(ctx context.Context, req *bizpb.ResetPasswordRequest) (*bizpb.ResetPasswordResponse, error) {
	resp := &bizpb.ResetPasswordResponse{}
	err := s.CallRPC(req, func() error {
		return s.resetPassword(req, resp)
	})
	return resp, err
}

func (s *Service) resetPassword(req *bizpb.ResetPasswordRequest, _ *bizpb.ResetPasswordResponse) error {
	// 标准化邮箱地址
	email, err := utils.NormalizeEmailAddress(req.Email)
	if err != nil {
		return rpc_error.BadRequest(rpc_error.DetailInvalidEmailAddress, err.Error())
	}
	emailHash := utils.HashEmailAddress(email)

	// 校验密码
	password, err := utils.ValidatePassword(req.NewPassword)
	if err != nil {
		if errors.Is(err, utils.ErrPasswordTooShortOrTooLong) {
			return rpc_error.BadRequest(rpc_error.DetailPasswordTooShortOrTooLong, err.Error())
		} else if errors.Is(err, utils.ErrPasswordTooWeak) {
			return rpc_error.BadRequest(rpc_error.DetailPasswordTooWeak, err.Error())
		}
		return err
	}

	// 校验邮箱验证码，验证成功后即消费，防止同一个验证码被重复使用
	if err := utils.VerifyEmailCode(s.Ctx(), s.Redis(), s.Config().EmailCodeSecret, enum.EmailCodeSceneString(commonpb.EmailCodeScene_Reset), emailHash, req.EmailCode); err != nil {
		if !utils.IsEmailCodeBusinessError(err) {
			return err
		}
		return rpc_error.BadRequest(rpc_error.DetailIncorrectEmailCode, err.Error())
	}

	// 生成密码哈希
	passwordHash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	// 按邮箱更新用户密码
	if err := s.MySQL().UpdatePasswordByEmail(s.Ctx(), email, string(passwordHash)); err != nil {
		return err
	}

	// 重置密码后吊销所有 Refresh Token
	if err := s.MySQL().RevokeRefreshTokensByEmail(s.Ctx(), email); err != nil {
		common.RpcWarn("Revoke refresh tokens by email failed, email: %s, err: %v", email, err)
	}

	return nil
}
