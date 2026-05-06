package auth

import (
	"context"
	"fmt"
	"time"

	"icw_common/enum"
	"icw_common/gen/core/biz"
	"icw_common/rpc_err"
	"icw_core_biz/internal/services/auth/consts"
	"icw_core_biz/internal/services/auth/utils"
	"icw_core_biz/internal/services/common"
)

// SendEmailCode 发送邮箱验证码
func (s *Service) SendEmailCode(ctx context.Context, req *bizpb.SendEmailCodeRequest) (*bizpb.SendEmailCodeResponse, error) {
	resp := &bizpb.SendEmailCodeResponse{}
	err := s.CallRPC(ctx, req, resp, func() error {
		return s.sendEmailCode(req, resp)
	})
	return resp, err
}

func (s *Service) sendEmailCode(req *bizpb.SendEmailCodeRequest, resp *bizpb.SendEmailCodeResponse) error {
	// 标准化邮箱地址
	email, err := utils.NormalizeEmailAddress(req.Email)
	if err != nil {
		return rpc_err.BadRequest(rpc_err.DetailInvalidEmailAddress, err.Error())
	}
	emailHash := utils.HashEmailAddress(email)

	// 获取邮箱验证码业务场景枚举
	scene := enum.ParseEmailCodeScene(req.Scene)
	if scene == bizpb.EmailCodeScene_EMAIL_CODE_SCENE_UNKNOWN {
		return rpc_err.BadRequestDefault("invalid email code scene")
	}
	sceneValue := enum.EmailCodeSceneString(scene)

	// 注册场景：邮箱必须尚未注册
	if scene == bizpb.EmailCodeScene_EMAIL_CODE_SCENE_REGISTER {
		user, err := s.MySQL().FindUserByEmail(s.Ctx(), email)
		if err != nil {
			return err
		}
		if user != nil && user.Id != 0 {
			return rpc_err.BadRequest(rpc_err.DetailEmailAlreadyRegistered, "email already registered")
		}
	}

	// 登录和重置密码场景：邮箱必须已被注册
	if scene == bizpb.EmailCodeScene_EMAIL_CODE_SCENE_LOGIN || scene == bizpb.EmailCodeScene_EMAIL_CODE_SCENE_RESET {
		user, err := s.MySQL().FindUserByEmail(s.Ctx(), email)
		if err != nil {
			return err
		}
		if user == nil || user.Id == 0 {
			return rpc_err.BadRequest(rpc_err.DetailEmailNotRegistered, "email not registered")
		}
	}

	// 账号锁定（登录失败次数达上限）时不发送登录验证码
	if scene == bizpb.EmailCodeScene_EMAIL_CODE_SCENE_LOGIN {
		locked, ttl, err := s.Redis().IsLoginLocked(s.Ctx(), enum.LoginSceneString(bizpb.LoginScene_LOGIN_SCENE_EMAIL), emailHash, consts.LoginFailureLimit)
		if err != nil {
			return err
		}
		if locked {
			return rpc_err.AccountLockedDefault(fmt.Sprintf("login retry after %s", ttl.Round(time.Second)))
		}
	}

	// 验证码未过期时不发送登录验证码
	exists, err := s.Redis().EmailCodeExists(s.Ctx(), sceneValue, emailHash)
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
	if err := s.Redis().SaveEmailCode(s.Ctx(), sceneValue, emailHash, utils.HashEmailCode(code, s.Config().EmailCodeSecret), s.Config().EmailCodeTTL); err != nil {
		return err
	}

	// 发送邮箱验证码
	if err := s.SMTP().SendEmailCode(email, sceneValue, code); err != nil {
		s.recordEmailSendLog(s.Ctx(), email, sceneValue, code, bizpb.EmailSendStatus_EMAIL_SEND_STATUS_FAILED, err.Error())
		if err := s.Redis().ClearEmailCode(s.Ctx(), sceneValue, emailHash); err != nil {
			common.RpcWarn("Clear email code failed, scene: %s, email: %s, err: %v", sceneValue, email, err)
		}
		return rpc_err.InternalError(rpc_err.DetailSendEmailCodeFailed, err.Error())
	}
	s.recordEmailSendLog(s.Ctx(), email, sceneValue, code, bizpb.EmailSendStatus_EMAIL_SEND_STATUS_SUCCESS, "")

	// 邮箱验证码发送成功，返回验证码有效期
	resp.ExpiresInSeconds = int32(s.Config().EmailCodeTTL.Seconds())

	return nil
}

// recordEmailSendLog 记录邮件发送日志
func (s *Service) recordEmailSendLog(ctx context.Context, receiverEmail, scene, emailCode string, status bizpb.EmailSendStatus, errorMessage string) {
	if err := s.MySQL().CreateEmailSendLog(ctx, receiverEmail, s.Config().SMTPFromEmail, scene, emailCode, status, errorMessage); err != nil {
		common.RpcWarn("Record email send log failed, receiver_email: %s, sender_email: %s, scene: %s, email_code: %s, status: %s, error_message: %s", receiverEmail, s.Config().SMTPFromEmail, scene, emailCode, enum.EmailSendStatusString(status), err.Error())
	}
}
