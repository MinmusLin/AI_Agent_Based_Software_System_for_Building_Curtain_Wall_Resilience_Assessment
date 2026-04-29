package rpc_err

import (
	"fmt"
	"strings"
)

// Error 标准错误结构体
type Error struct {
	Code    Code
	Detail  DetailCode
	Message string
}

// Error 将标准错误结构体转换为输出格式
func (e *Error) Error() string {
	return fmt.Sprintf("%s|%s: %s", e.Code, e.Detail, e.Message)
}

// newError 创建标准错误
func newError(code Code, detail DetailCode, message string) error {
	if !IsCode(code) {
		code = DefaultCode()
	}
	if !IsDetailCode(detail) {
		detail = DefaultDetailCode(code)
	}

	message = strings.TrimSpace(message)
	if message == "" {
		message = "<nil>"
	}

	return &Error{
		Code:    code,
		Detail:  detail,
		Message: message,
	}
}

// BadRequest 创建无效请求标准错误
func BadRequest(detail DetailCode, message string) error {
	return newError(CodeBadRequest, detail, message)
}

// BadRequestDefault 创建无效请求标准错误（默认错误业务代码）
func BadRequestDefault(message string) error {
	return newError(CodeBadRequest, DefaultDetailCode(CodeBadRequest), message)
}

// Unauthorized 创建身份验证未通过标准错误
func Unauthorized(detail DetailCode, message string) error {
	return newError(CodeUnauthorized, detail, message)
}

// UnauthorizedDefault 创建身份验证未通过标准错误（默认错误业务代码）
func UnauthorizedDefault(message string) error {
	return newError(CodeUnauthorized, DefaultDetailCode(CodeUnauthorized), message)
}

// AccountLocked 创建账号锁定标准错误
func AccountLocked(detail DetailCode, message string) error {
	return newError(CodeAccountLocked, detail, message)
}

// AccountLockedDefault 创建账号锁定标准错误（默认错误业务代码）
func AccountLockedDefault(message string) error {
	return newError(CodeAccountLocked, DefaultDetailCode(CodeAccountLocked), message)
}

// InternalError 创建服务器内部错误标准错误
func InternalError(detail DetailCode, message string) error {
	return newError(CodeInternalError, detail, message)
}

// InternalErrorDefault 创建服务器内部错误标准错误（默认错误业务代码）
func InternalErrorDefault(message string) error {
	return newError(CodeInternalError, DefaultDetailCode(CodeInternalError), message)
}
