package auth

import (
	"context"
	"fmt"
	"log"
	"time"

	"icw_core_biz/internal/services/auth/consts"
	"icw_core_biz/internal/services/auth/utils"
	"icw_core_biz/pkg/dto"
	"icw_core_biz/pkg/rpc_err"
)

// SendEmailCode 发送邮箱验证码
func (s *Service) SendEmailCode(req *dto.SendEmailCodeRequest, resp *dto.SendEmailCodeResponse) error {
	return s.CallRPC("AuthService.SendEmailCode", req, resp, func() error {
		return s.sendEmailCode(req, resp)
	})
}

func (s *Service) sendEmailCode(req *dto.SendEmailCodeRequest, resp *dto.SendEmailCodeResponse) (err error) {
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
		user, err := s.MySQL().FindUserByEmail(s.Ctx, email)
		if err != nil {
			return err
		}
		if user != nil && user.Id != 0 {
			return rpc_err.BadRequest(rpc_err.DetailEmailAlreadyRegistered, "email already registered")
		}
	}

	// 登录和重置密码场景：邮箱必须已被注册
	if scene == consts.SceneLogin || scene == consts.SceneReset {
		user, err := s.MySQL().FindUserByEmail(s.Ctx, email)
		if err != nil {
			return err
		}
		if user == nil || user.Id == 0 {
			return rpc_err.BadRequest(rpc_err.DetailEmailNotRegistered, "email not registered")
		}
	}

	// 账号锁定（登录失败次数达上限）时不发送登录验证码
	if scene == consts.SceneLogin {
		locked, ttl, err := s.Redis().IsLoginLocked(s.Ctx, consts.LoginEmail.String(), email, consts.LoginFailureLimit)
		if err != nil {
			return err
		}
		if locked {
			return rpc_err.AccountLockedDefault(fmt.Sprintf("login retry after %s", ttl.Round(time.Second)))
		}
	}

	// 验证码未过期时不发送登录验证码
	exists, err := s.Redis().EmailCodeExists(s.Ctx, sceneValue, email)
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
	if err := s.Redis().SaveEmailCode(s.Ctx, sceneValue, email, utils.HashEmailCode(code, s.Config().EmailCodeSecret), s.Config().EmailCodeTTL); err != nil {
		return err
	}

	// 发送邮箱验证码
	if err := s.SMTP().SendEmailCode(email, sceneValue, code); err != nil {
		s.recordEmailSendLog(s.Ctx, email, sceneValue, code, consts.EmailSendStatusFailed, err.Error())
		if err := s.Redis().ClearEmailCode(s.Ctx, sceneValue, email); err != nil {
			log.Printf("[WARN] Clear email code failed, scene: %s, email: %s, err: %v", sceneValue, email, err)
		}
		return rpc_err.InternalError(rpc_err.DetailSendEmailCodeFailed, err.Error())
	}
	s.recordEmailSendLog(s.Ctx, email, sceneValue, code, consts.EmailSendStatusSuccess, "")

	// 邮箱验证码发送成功，返回验证码有效期
	resp.ExpiresInSeconds = int(s.Config().EmailCodeTTL.Seconds())

	return nil
}

// recordEmailSendLog 记录邮件发送日志
func (s *Service) recordEmailSendLog(ctx context.Context, receiverEmail, scene, emailCode string, status consts.EmailSendStatus, errorMessage string) {
	if err := s.MySQL().CreateEmailSendLog(ctx, receiverEmail, s.Config().SMTPFromEmail, scene, emailCode, status, errorMessage); err != nil {
		log.Printf("[WARN] Record email send log failed, receiver_email: %s, sender_email: %s, scene: %s, email_code: %s, status: %s, error_message: %s", receiverEmail, s.Config().SMTPFromEmail, scene, emailCode, status.String(), err.Error())
	}
}
