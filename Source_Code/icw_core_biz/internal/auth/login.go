package auth

import (
	"context"
	"fmt"
	"time"

	"golang.org/x/crypto/bcrypt"

	"icw_core_biz/internal/auth/consts"
	"icw_core_biz/internal/auth/utils"
	"icw_core_biz/pkg/dto"
	"icw_core_biz/pkg/rpc_err"
)

// Login 登录
func (s *Service) Login(req *dto.LoginRequest, resp *dto.LoginResponse) error {
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
	locked, ttl, err := s.redis.IsLoginLocked(ctx, scene.String(), email, consts.LoginFailureLimit)
	if err != nil {
		return err
	}
	if locked {
		return rpc_err.AccountLockedDefault(fmt.Sprintf("login retry after %s", ttl.Round(time.Second)))
	}

	// 按邮箱查询用户
	user, err := s.mysql.FindUserByEmail(ctx, email)
	if err != nil {
		return err
	}
	if user == nil || user.Id == 0 {
		// 用户不存在视作登录失败，避免泄露邮箱是否存在
		if err := s.redis.RecordLoginFailure(ctx, scene.String(), email, s.cfg.LoginFailTTL); err != nil {
			return err
		}
		return rpc_err.Unauthorized(rpc_err.DetailInvalidCredentials, "user not found")
	}

	switch scene {
	case consts.LoginPassword:
		// 密码登录
		if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(req.Code)); err != nil {
			if err := s.redis.RecordLoginFailure(ctx, scene.String(), email, s.cfg.LoginFailTTL); err != nil {
				return err
			}
			return rpc_err.Unauthorized(rpc_err.DetailInvalidCredentials, err.Error())
		}
	case consts.LoginEmail:
		// 邮箱验证码登录
		if err := utils.VerifyEmailCode(ctx, s.redis, s.cfg.EmailCodeSecret, consts.SceneLogin.String(), email, req.Code); err != nil {
			if !utils.IsEmailCodeBusinessError(err) {
				return err
			}
			if err := s.redis.RecordLoginFailure(ctx, scene.String(), email, s.cfg.LoginFailTTL); err != nil {
				return err
			}
			return rpc_err.BadRequest(rpc_err.DetailIncorrectEmailCode, err.Error())
		}
	default:
		return rpc_err.BadRequestDefault("invalid login scene")
	}

	// 登录成功后清除登录失败计数
	if err := s.redis.ClearLoginFailure(ctx, scene.String(), email); err != nil {
		return err
	}

	// 签发新的 Refresh Token 和 Access Token
	if err := utils.IssueTokens(ctx, s.cfg, s.mysql, s.tokens, user, resp); err != nil {
		return err
	}

	return nil
}
