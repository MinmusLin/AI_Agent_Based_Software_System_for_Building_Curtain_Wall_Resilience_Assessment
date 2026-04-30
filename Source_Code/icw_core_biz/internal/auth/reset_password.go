package auth

import (
	"context"
	"errors"

	"golang.org/x/crypto/bcrypt"

	"icw_core_biz/internal/auth/consts"
	"icw_core_biz/internal/auth/utils"
	"icw_core_biz/pkg/dto"
	"icw_core_biz/pkg/rpc_err"
)

// ResetPassword 重置密码
func (s *Service) ResetPassword(req *dto.ResetPasswordRequest, _ *dto.ResetPasswordResponse) error {
	if req == nil {
		return rpc_err.BadRequestDefault("request is nil")
	}
	ctx := context.Background()

	// 标准化邮箱地址
	email, err := utils.NormalizeEmailAddress(req.Email)
	if err != nil {
		return rpc_err.BadRequest(rpc_err.DetailInvalidEmailAddress, err.Error())
	}

	// 校验密码
	password, err := utils.ValidatePassword(req.NewPassword)
	if err != nil {
		if errors.Is(err, utils.ErrPasswordTooShort) {
			return rpc_err.BadRequest(rpc_err.DetailPasswordTooShort, err.Error())
		} else if errors.Is(err, utils.ErrPasswordTooWeak) {
			return rpc_err.BadRequest(rpc_err.DetailPasswordTooWeak, err.Error())
		}
		return err
	}

	// 校验邮箱验证码，验证成功后即消费，防止同一个验证码被重复使用
	if err := utils.VerifyEmailCode(ctx, s.redis, s.cfg.EmailCodeSecret, consts.SceneReset.String(), email, req.EmailCode); err != nil {
		if !utils.IsEmailCodeBusinessError(err) {
			return err
		}
		return rpc_err.BadRequest(rpc_err.DetailIncorrectEmailCode, err.Error())
	}

	// 生成密码哈希
	passwordHash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	// 按邮箱更新用户密码
	if err := s.mysql.UpdatePasswordByEmail(ctx, email, string(passwordHash)); err != nil {
		return err
	}

	// 重置密码后吊销所有 Refresh Token
	_ = s.mysql.RevokeUserRefreshTokensByEmail(ctx, email)

	return nil
}
