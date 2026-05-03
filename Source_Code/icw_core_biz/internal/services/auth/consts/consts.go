package consts

import (
	"strings"
)

const (
	// LoginFailureLimit 登录失败最大重试次数
	LoginFailureLimit int = 5
)

const (
	// MinPasswordLength 密码最小字符数
	MinPasswordLength int = 8
	// MaxPasswordLength 密码最大字符数
	MaxPasswordLength int = 24
	// MaxNameLength 用户名称最大字符数
	MaxNameLength int = 8
)

// EmailSendStatus 邮件发送状态枚举
type EmailSendStatus string

const (
	// EmailSendStatusSuccess 邮件发送成功
	EmailSendStatusSuccess EmailSendStatus = "success"
	// EmailSendStatusFailed 邮件发送失败
	EmailSendStatusFailed EmailSendStatus = "failed"
)

// String 将邮件发送状态枚举转换为字符串
func (s EmailSendStatus) String() string {
	return string(s)
}

// ParseEmailSendStatus 将外部输入转换为邮件发送状态枚举
func ParseEmailSendStatus(value string) EmailSendStatus {
	switch status := EmailSendStatus(strings.TrimSpace(value)); status {
	case EmailSendStatusSuccess, EmailSendStatusFailed:
		return status
	default:
		return ""
	}
}

// EmailCodeScene 邮箱验证码业务场景枚举
type EmailCodeScene string

const (
	// SceneRegister 账号注册
	SceneRegister EmailCodeScene = "register"
	// SceneLogin 账号登录
	SceneLogin EmailCodeScene = "login"
	// SceneReset 重置密码
	SceneReset EmailCodeScene = "reset"
)

// String 将邮箱验证码业务场景枚举转换为字符串
func (s EmailCodeScene) String() string {
	return string(s)
}

// ParseEmailCodeScene 将外部输入转换为邮箱验证码业务场景枚举
func ParseEmailCodeScene(value string) EmailCodeScene {
	switch scene := EmailCodeScene(strings.TrimSpace(value)); scene {
	case SceneRegister, SceneLogin, SceneReset:
		return scene
	default:
		return ""
	}
}

// LoginScene 登录方式枚举
type LoginScene string

const (
	// LoginPassword 密码登录方式
	LoginPassword LoginScene = "password"
	// LoginEmail 邮箱验证码登录方式
	LoginEmail LoginScene = "email"
)

// String 将登录方式枚举转换为字符串
func (s LoginScene) String() string {
	return string(s)
}

// ParseLoginScene 将外部输入转换为登录方式枚举
func ParseLoginScene(value string) LoginScene {
	switch scene := LoginScene(strings.TrimSpace(value)); scene {
	case LoginPassword, LoginEmail:
		return scene
	default:
		return ""
	}
}
