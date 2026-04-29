package consts

import (
	"strings"
)

// LoginFailureLimit 登录失败最大重试次数
const LoginFailureLimit int = 5

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
