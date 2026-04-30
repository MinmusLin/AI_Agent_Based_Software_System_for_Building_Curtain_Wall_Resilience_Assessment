package utils

import (
	"context"
	"errors"
	"net/mail"
	"strings"

	"icw_core_biz/repositories"
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
func VerifyEmailCode(ctx context.Context, redis *repositories.RedisRepository, secret, scene, email, code string) error {
	code = strings.TrimSpace(code)
	if len(code) != 6 {
		return errors.New("invalid email code")
	}
	codeHash, err := redis.GetEmailCode(ctx, scene, email)
	if err != nil {
		return err
	}
	if codeHash == "" {
		return errors.New("email code not found or expired")
	}
	if codeHash != HashEmailCode(code, secret) {
		return errors.New("incorrect email code")
	}
	if err := redis.ClearEmailCode(ctx, scene, email); err != nil {
		return err
	}
	return nil
}
