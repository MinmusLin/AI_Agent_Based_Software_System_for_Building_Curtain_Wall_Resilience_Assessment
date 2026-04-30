package auth

import (
	"context"
	"errors"

	"golang.org/x/crypto/bcrypt"

	"icw_core_biz/internal/auth/consts"
	"icw_core_biz/internal/auth/utils"
	"icw_core_biz/internal/rpc_log"
	"icw_core_biz/pkg/dto"
	"icw_core_biz/pkg/rpc_err"
	"icw_core_biz/repositories/mysql"
)

// Register 注册
func (s *Service) Register(req *dto.RegisterRequest, resp *dto.RegisterResponse) (err error) {
	start := rpc_log.Start("AuthService.Register", req)
	defer func() {
		rpc_log.Finish("AuthService.Register", req, resp, start, err)
	}()

	if req == nil {
		return rpc_err.BadRequestDefault("request is nil")
	}
	ctx := context.Background()

	// 标准化邮箱地址
	email, err := utils.NormalizeEmailAddress(req.Email)
	if err != nil {
		return rpc_err.BadRequest(rpc_err.DetailInvalidEmailAddress, err.Error())
	}

	// 校验用户名称
	name, err := utils.ValidateName(req.Name)
	if err != nil {
		if errors.Is(err, utils.ErrNameIsEmpty) {
			return rpc_err.BadRequest(rpc_err.DetailNameRequired, err.Error())
		} else if errors.Is(err, utils.ErrNameTooLong) {
			return rpc_err.BadRequest(rpc_err.DetailNameTooLong, err.Error())
		}
		return err
	}

	// 校验密码
	password, err := utils.ValidatePassword(req.Password)
	if err != nil {
		if errors.Is(err, utils.ErrPasswordTooShortOrTooLong) {
			return rpc_err.BadRequest(rpc_err.DetailPasswordTooShortOrTooLong, err.Error())
		} else if errors.Is(err, utils.ErrPasswordTooWeak) {
			return rpc_err.BadRequest(rpc_err.DetailPasswordTooWeak, err.Error())
		}
		return err
	}

	// 校验邮箱验证码，验证成功后即消费，防止同一个验证码被重复使用
	if err := utils.VerifyEmailCode(ctx, s.redis, s.cfg.EmailCodeSecret, consts.SceneRegister.String(), email, req.EmailCode); err != nil {
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

	// 创建用户
	if err := s.mysql.CreateUser(ctx, email, string(passwordHash), name); err != nil {
		if mysql.IsDuplicateEntryError(err) {
			return rpc_err.BadRequest(rpc_err.DetailEmailAlreadyRegistered, "email already registered")
		}
		return err
	}

	return nil
}
