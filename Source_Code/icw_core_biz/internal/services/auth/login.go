package auth

import (
	"context"
	"fmt"
	"time"

	"golang.org/x/crypto/bcrypt"

	"icw_common/enum"
	"icw_common/gen/core/biz"
	"icw_common/rpc/error"

	"icw_core_biz/internal/services/auth/consts"
	"icw_core_biz/internal/services/auth/utils"
)

// Login 登录
func (s *Service) Login(ctx context.Context, req *bizpb.LoginRequest) (*bizpb.LoginResponse, error) {
	resp := &bizpb.LoginResponse{}
	err := s.CallRPC(ctx, req, func() error {
		return s.login(req, resp)
	})
	return resp, err
}

func (s *Service) login(req *bizpb.LoginRequest, resp *bizpb.LoginResponse) error {
	// 标准化邮箱地址
	email, err := utils.NormalizeEmailAddress(req.Email)
	if err != nil {
		return rpc_error.BadRequest(rpc_error.DetailInvalidEmailAddress, err.Error())
	}
	emailHash := utils.HashEmailAddress(email)

	// 获取登录方式枚举
	scene := enum.ParseLoginScene(req.Scene)
	if scene == bizpb.LoginScene_Unknown {
		return rpc_error.BadRequestDefault("invalid login scene")
	}

	// 账号锁定（登录失败次数达上限）时不进行登录操作
	locked, ttl, err := s.Redis().IsLoginLocked(s.Ctx(), enum.LoginSceneString(scene), emailHash, consts.LoginFailureLimit)
	if err != nil {
		return err
	}
	if locked {
		return rpc_error.AccountLockedDefault(fmt.Sprintf("login retry after %s", ttl.Round(time.Second)))
	}

	// 按邮箱查询用户
	user, err := s.MySQL().FindUserByEmail(s.Ctx(), email)
	if err != nil {
		return err
	}
	if user == nil || user.Id == 0 {
		// 用户不存在视作登录失败，避免泄露邮箱是否存在
		if err := s.Redis().RecordLoginFailure(s.Ctx(), enum.LoginSceneString(scene), emailHash, s.Config().LoginFailTTL); err != nil {
			return err
		}
		return rpc_error.BadRequest(rpc_error.DetailInvalidCredentials, "user not found")
	}

	switch scene {
	case bizpb.LoginScene_Password:
		// 密码登录
		if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(req.Code)); err != nil {
			if err := s.Redis().RecordLoginFailure(s.Ctx(), enum.LoginSceneString(scene), emailHash, s.Config().LoginFailTTL); err != nil {
				return err
			}
			return rpc_error.BadRequest(rpc_error.DetailInvalidCredentials, err.Error())
		}
	case bizpb.LoginScene_Email:
		// 邮箱验证码登录
		if err := utils.VerifyEmailCode(s.Ctx(), s.Redis(), s.Config().EmailCodeSecret, enum.EmailCodeSceneString(bizpb.EmailCodeScene_Login), emailHash, req.Code); err != nil {
			if !utils.IsEmailCodeBusinessError(err) {
				return err
			}
			if err := s.Redis().RecordLoginFailure(s.Ctx(), enum.LoginSceneString(scene), emailHash, s.Config().LoginFailTTL); err != nil {
				return err
			}
			return rpc_error.BadRequest(rpc_error.DetailIncorrectEmailCode, err.Error())
		}
	default:
		return rpc_error.BadRequestDefault("invalid login scene")
	}

	// 登录成功后清除登录失败计数
	if err := s.Redis().ClearLoginFailure(s.Ctx(), enum.LoginSceneString(bizpb.LoginScene_Password), emailHash); err != nil {
		return err
	}
	if err := s.Redis().ClearLoginFailure(s.Ctx(), enum.LoginSceneString(bizpb.LoginScene_Email), emailHash); err != nil {
		return err
	}

	// 签发 Access Token 和 Refresh Token，并更新用户最近登录时间
	if err := utils.IssueTokens(s.Ctx(), s.Config(), s.MySQL(), s.tokens, user, resp); err != nil {
		return err
	}

	return nil
}
