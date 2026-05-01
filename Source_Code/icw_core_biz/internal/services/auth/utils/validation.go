package utils

import (
	"errors"
	"strings"
	"unicode"
	"unicode/utf8"

	"icw_core_biz/internal/services/auth/consts"
)

var (
	// ErrPasswordTooShortOrTooLong 密码长度过短或过长
	ErrPasswordTooShortOrTooLong = errors.New("password too short or too long")
	// ErrPasswordTooWeak 密码强度过弱
	ErrPasswordTooWeak = errors.New("password too weak")
	// ErrNameIsEmpty 用户名称为空
	ErrNameIsEmpty = errors.New("name is empty")
	// ErrNameTooLong 用户名称长度过长
	ErrNameTooLong = errors.New("name too long")
)

// ValidatePassword 校验密码
// 密码长度必须不小于 8 位，不多于 24 位，且必须同时包含大小写英文字母、数字和符号
func ValidatePassword(password string) (string, error) {
	password = strings.TrimSpace(password)
	if utf8.RuneCountInString(password) < consts.MinPasswordLength || utf8.RuneCountInString(password) > consts.MaxPasswordLength {
		return "", ErrPasswordTooShortOrTooLong
	}
	var hasUpper, hasLower, hasDigit, hasSymbol bool
	for _, item := range password {
		switch {
		case unicode.IsUpper(item):
			hasUpper = true
		case unicode.IsLower(item):
			hasLower = true
		case unicode.IsDigit(item):
			hasDigit = true
		case unicode.IsPunct(item) || unicode.IsSymbol(item):
			hasSymbol = true
		}
	}
	if !hasUpper || !hasLower || !hasDigit || !hasSymbol {
		return "", ErrPasswordTooWeak
	}
	return password, nil
}

// ValidateName 校验用户名称
// 用户名称不能超过 8 个字符，且不能为空
func ValidateName(name string) (string, error) {
	name = strings.TrimSpace(name)
	if name == "" {
		return "", ErrNameIsEmpty
	}
	if utf8.RuneCountInString(name) > consts.MaxNameLength {
		return "", ErrNameTooLong
	}
	return name, nil
}
