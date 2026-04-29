package auth

import (
	"context"
	"fmt"
	"time"

	"icw_core_biz/internal/auth/consts"
	"icw_core_biz/internal/auth/utils"
	"icw_core_biz/pkg/dto"
	"icw_core_biz/pkg/rpc_err"
)

// SendEmailCode 发送邮箱验证码
func (s *Service) SendEmailCode(req *dto.SendEmailCodeRequest, resp *dto.SendEmailCodeResponse) error {
	if req == nil {
		return rpc_err.BadRequestDefault("request is nil")
	}
	ctx := context.Background()

	// 标准化邮箱地址
	email, err := utils.NormalizeEmailAddress(req.Email)
	if err != nil {
		return rpc_err.BadRequest(rpc_err.DetailInvalidEmailAddress, err.Error())
	}

	// 获取邮箱验证码业务场景枚举
	scene := consts.ParseEmailCodeScene(req.Scene)
	if scene == "" {
		return rpc_err.BadRequestDefault("invalid email code scene")
	}
	sceneValue := scene.String()

	// 注册场景：邮箱必须尚未注册
	if scene == consts.SceneRegister {
		exists, err := s.mysql.UserExists(ctx, email)
		if err != nil {
			return err
		}
		if exists {
			return rpc_err.BadRequest(rpc_err.DetailEmailAlreadyRegistered, "email already registered")
		}
	}

	// 登录和重置密码场景：邮箱必须已被注册
	if scene == consts.SceneLogin || scene == consts.SceneReset {
		exists, err := s.mysql.UserExists(ctx, email)
		if err != nil {
			return err
		}
		if !exists {
			return rpc_err.BadRequest(rpc_err.DetailEmailNotRegistered, "email not registered")
		}
	}

	// 账号锁定（登录失败次数达上限）时不发送登录验证码
	if scene == consts.SceneLogin {
		if locked, ttl := s.redis.IsLoginLocked(ctx, consts.LoginEmail.String(), email, consts.LoginFailureLimit); locked {
			return rpc_err.AccountLocked(rpc_err.DetailAccountLocked, fmt.Sprintf("login retry after %s", ttl.Round(time.Second)))
		}
	}

	// 验证码未过期时不发送登录验证码
	exists, err := s.redis.EmailCodeExists(ctx, sceneValue, email)
	if err != nil {
		return err
	}
	if exists {
		return rpc_err.BadRequest(rpc_err.DetailEmailCodeSentTooFrequently, "email code sent too frequently")
	}

	// 生成 6 位数字邮箱验证码
	code, err := utils.NewEmailCode()
	if err != nil {
		return err
	}
	if err := s.redis.SaveEmailCode(ctx, sceneValue, email, utils.HashEmailCode(code, s.cfg.EmailCodeSecret), s.cfg.EmailCodeTTL); err != nil {
		return err
	}

	// 发送邮箱验证码
	if err := s.smtp.SendEmailCode(email, sceneValue, code); err != nil {
		_ = s.redis.ClearEmailCode(ctx, sceneValue, email)
		return rpc_err.InternalError(rpc_err.DetailSendEmailCodeFailed, err.Error())
	}

	// 邮箱验证码发送成功，返回验证码有效期
	resp.ExpiresInSeconds = int(s.cfg.EmailCodeTTL.Seconds())
	return nil
}
