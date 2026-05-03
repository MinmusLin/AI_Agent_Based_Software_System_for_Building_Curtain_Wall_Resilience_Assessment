package auth

import (
	"errors"
	"log"

	"golang.org/x/crypto/bcrypt"

	"icw_core_biz/internal/services/auth/consts"
	"icw_core_biz/internal/services/auth/utils"
	"icw_core_biz/internal/services/common"
	"icw_core_biz/pkg/dto"
	"icw_core_biz/pkg/rpc_err"
)

// ResetPassword 重置密码
func (s *Service) ResetPassword(req *dto.ResetPasswordRequest, resp *dto.ResetPasswordResponse) error {
	return s.CallRPC("AuthService.ResetPassword", req, resp, func() error {
		return s.resetPassword(req, resp)
	})
}

func (s *Service) resetPassword(req *dto.ResetPasswordRequest, _ *dto.ResetPasswordResponse) error {
	// 标准化邮箱地址
	email, err := utils.NormalizeEmailAddress(req.Email)
	if err != nil {
		return rpc_err.BadRequest(rpc_err.DetailInvalidEmailAddress, err.Error())
	}

	// 校验密码
	password, err := utils.ValidatePassword(req.NewPassword)
	if err != nil {
		if errors.Is(err, utils.ErrPasswordTooShortOrTooLong) {
			return rpc_err.BadRequest(rpc_err.DetailPasswordTooShortOrTooLong, err.Error())
		} else if errors.Is(err, utils.ErrPasswordTooWeak) {
			return rpc_err.BadRequest(rpc_err.DetailPasswordTooWeak, err.Error())
		}
		return err
	}

	// 校验邮箱验证码，验证成功后即消费，防止同一个验证码被重复使用
	if err := utils.VerifyEmailCode(s.Ctx(), s.Redis(), s.Config().EmailCodeSecret, consts.SceneReset.String(), email, req.EmailCode); err != nil {
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
	if err := s.MySQL().UpdatePasswordByEmail(s.Ctx(), email, string(passwordHash)); err != nil {
		return err
	}

	// 重置密码后吊销所有 Refresh Token
	if err := s.MySQL().RevokeRefreshTokensByEmail(s.Ctx(), email); err != nil {
		log.Printf("%s Revoke refresh tokens by email failed, email: %s, err: %v", common.WarnPrefix(), email, err)
	}

	return nil
}
