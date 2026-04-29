package rpc_err

import (
	"strings"
)

// Parse 解析 RPC 标准错误
func Parse(err error) (Code, DetailCode, string) {
	if err == nil {
		return CodeInternalError, DetailInternalError, "<nil>"
	}
	return ParseString(err.Error())
}

// ParseString 从 "<CODE>|<DETAIL_CODE>: <error_message>" 中解析 RPC 标准错误
func ParseString(errStr string) (Code, DetailCode, string) {
	header, message, ok := strings.Cut(errStr, ":")
	if !ok || header == "" || message == "" {
		return CodeInternalError, DetailInternalError, strings.TrimSpace(errStr)
	}

	codeStr, detailCodeStr, ok := strings.Cut(strings.TrimSpace(header), "|")
	if !ok || codeStr == "" || detailCodeStr == "" {
		return CodeInternalError, DetailInternalError, strings.TrimSpace(errStr)
	}

	code := Code(strings.TrimSpace(codeStr))
	if !IsCode(code) {
		code = DefaultCode()
	}

	detailCode := DetailCode(strings.TrimSpace(detailCodeStr))
	if !IsDetailCode(detailCode) {
		detailCode = DefaultDetailCode(code)
	}

	message = strings.TrimSpace(message)
	if message == "" {
		message = "<nil>"
	}

	return code, detailCode, message
}
