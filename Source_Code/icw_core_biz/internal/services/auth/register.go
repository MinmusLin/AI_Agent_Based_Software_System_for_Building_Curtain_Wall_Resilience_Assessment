package auth

import (
	"context"
	"errors"

	"golang.org/x/crypto/bcrypt"

	"icw_common/enum"
	"icw_common/gen/core/biz"
	"icw_common/rpc/error"

	"icw_core_biz/internal/services/auth/utils"
	mysqlUtils "icw_core_biz/repositories/mysql/utils"
)

// Register 注册
func (s *Service) Register(ctx context.Context, req *bizpb.RegisterRequest) (*bizpb.RegisterResponse, error) {
	resp := &bizpb.RegisterResponse{}
	err := s.CallRPC(req, func() error {
		return s.register(req, resp)
	})
	return resp, err
}

func (s *Service) register(req *bizpb.RegisterRequest, _ *bizpb.RegisterResponse) error {
	// 标准化邮箱地址
	email, err := utils.NormalizeEmailAddress(req.Email)
	if err != nil {
		return rpc_error.BadRequest(rpc_error.DetailInvalidEmailAddress, err.Error())
	}
	emailHash := utils.HashEmailAddress(email)

	// 校验用户名称
	name, err := utils.ValidateName(req.Name)
	if err != nil {
		if errors.Is(err, utils.ErrNameIsEmpty) {
			return rpc_error.BadRequest(rpc_error.DetailNameRequired, err.Error())
		} else if errors.Is(err, utils.ErrNameTooLong) {
			return rpc_error.BadRequest(rpc_error.DetailNameTooLong, err.Error())
		}
		return err
	}

	// 校验密码
	password, err := utils.ValidatePassword(req.Password)
	if err != nil {
		if errors.Is(err, utils.ErrPasswordTooShortOrTooLong) {
			return rpc_error.BadRequest(rpc_error.DetailPasswordTooShortOrTooLong, err.Error())
		} else if errors.Is(err, utils.ErrPasswordTooWeak) {
			return rpc_error.BadRequest(rpc_error.DetailPasswordTooWeak, err.Error())
		}
		return err
	}

	// 校验邮箱验证码，验证成功后即消费，防止同一个验证码被重复使用
	if err := utils.VerifyEmailCode(s.Ctx(), s.Redis(), s.Config().EmailCodeSecret, enum.EmailCodeSceneString(bizpb.EmailCodeScene_Register), emailHash, req.EmailCode); err != nil {
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

	// 创建用户
	if err := s.MySQL().CreateUser(s.Ctx(), email, string(passwordHash), name); err != nil {
		if mysqlUtils.IsDuplicateEntryError(err) {
			return rpc_error.BadRequest(rpc_error.DetailEmailAlreadyRegistered, "email already registered")
		}
		return err
	}

	return nil
}
