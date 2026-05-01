package auth

import (
	"context"
	"fmt"
	"time"

	"golang.org/x/crypto/bcrypt"

	"icw_core_biz/internal/rpc_log"
	"icw_core_biz/internal/services/auth/consts"
	"icw_core_biz/internal/services/auth/utils"
	"icw_core_biz/pkg/dto"
	"icw_core_biz/pkg/rpc_err"
)

// Login 登录
func (s *Service) Login(req *dto.LoginRequest, resp *dto.LoginResponse) (err error) {
	start := rpc_log.Start("AuthService.Login", req)
	defer func() {
		rpc_log.Finish("AuthService.Login", req, resp, start, err)
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

	// 获取登录方式枚举
	scene := consts.ParseLoginScene(req.Scene)
	if scene == "" {
		return rpc_err.BadRequestDefault("invalid login scene")
	}

	// 账号锁定（登录失败次数达上限）时不进行登录操作
	locked, ttl, err := s.Redis().IsLoginLocked(ctx, scene.String(), email, consts.LoginFailureLimit)
	if err != nil {
		return err
	}
	if locked {
		return rpc_err.AccountLockedDefault(fmt.Sprintf("login retry after %s", ttl.Round(time.Second)))
	}

	// 按邮箱查询用户
	user, err := s.MySQL().FindUserByEmail(ctx, email)
	if err != nil {
		return err
	}
	if user == nil || user.Id == 0 {
		// 用户不存在视作登录失败，避免泄露邮箱是否存在
		if err := s.Redis().RecordLoginFailure(ctx, scene.String(), email, s.Config().LoginFailTTL); err != nil {
			return err
		}
		return rpc_err.BadRequest(rpc_err.DetailInvalidCredentials, "user not found")
	}

	switch scene {
	case consts.LoginPassword:
		// 密码登录
		if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(req.Code)); err != nil {
			if err := s.Redis().RecordLoginFailure(ctx, scene.String(), email, s.Config().LoginFailTTL); err != nil {
				return err
			}
			return rpc_err.BadRequest(rpc_err.DetailInvalidCredentials, err.Error())
		}
	case consts.LoginEmail:
		// 邮箱验证码登录
		if err := utils.VerifyEmailCode(ctx, s.Redis(), s.Config().EmailCodeSecret, consts.SceneLogin.String(), email, req.Code); err != nil {
			if !utils.IsEmailCodeBusinessError(err) {
				return err
			}
			if err := s.Redis().RecordLoginFailure(ctx, scene.String(), email, s.Config().LoginFailTTL); err != nil {
				return err
			}
			return rpc_err.BadRequest(rpc_err.DetailIncorrectEmailCode, err.Error())
		}
	default:
		return rpc_err.BadRequestDefault("invalid login scene")
	}

	// 登录成功后清除登录失败计数
	if err := s.Redis().ClearLoginFailure(ctx, consts.LoginPassword.String(), email); err != nil {
		return err
	}
	if err := s.Redis().ClearLoginFailure(ctx, consts.LoginEmail.String(), email); err != nil {
		return err
	}

	// 签发 Access Token 和 Refresh Token，并更新用户最近登录时间
	if err := utils.IssueTokens(ctx, s.Config(), s.MySQL(), s.tokens, user, resp); err != nil {
		return err
	}

	return nil
}
