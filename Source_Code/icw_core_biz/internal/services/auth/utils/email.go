package utils

import (
	"context"
	"errors"
	"net/mail"
	"strings"

	"icw_core_biz/repositories/redis"
)

var (
	// ErrInvalidEmailCode 邮箱验证码格式错误
	ErrInvalidEmailCode = errors.New("invalid email code")
	// ErrEmailCodeNotFound 邮箱验证码不存在或已过期
	ErrEmailCodeNotFound = errors.New("email code not found or expired")
	// ErrIncorrectEmailCode 邮箱验证码错误
	ErrIncorrectEmailCode = errors.New("incorrect email code")
)

// NormalizeEmailAddress 标准化邮箱地址
func NormalizeEmailAddress(value string) (string, error) {
	input := strings.TrimSpace(value)
	if input == "" {
		return "", errors.New("email address is empty")
	}
	address, err := mail.ParseAddress(input)
	if err != nil || address == nil || address.Address == "" {
		return "", errors.New("invalid email address")
	}
	return strings.ToLower(address.Address), nil
}

// VerifyEmailCode 校验邮箱验证码
func VerifyEmailCode(ctx context.Context, redis *redis.Repository, secret, scene, email, code string) error {
	code = strings.TrimSpace(code)
	if len(code) != 6 {
		return ErrInvalidEmailCode
	}
	codeHash, err := redis.GetEmailCode(ctx, scene, email)
	if err != nil {
		return err
	}
	if codeHash == "" {
		return ErrEmailCodeNotFound
	}
	if codeHash != HashEmailCode(code, secret) {
		return ErrIncorrectEmailCode
	}
	if err := redis.ClearEmailCode(ctx, scene, email); err != nil {
		return err
	}
	return nil
}

// IsEmailCodeBusinessError 判断是否为用户输入导致的验证码错误
func IsEmailCodeBusinessError(err error) bool {
	return errors.Is(err, ErrInvalidEmailCode) || errors.Is(err, ErrEmailCodeNotFound) || errors.Is(err, ErrIncorrectEmailCode)
}
