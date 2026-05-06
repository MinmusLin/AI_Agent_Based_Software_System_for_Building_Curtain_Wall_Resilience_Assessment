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
