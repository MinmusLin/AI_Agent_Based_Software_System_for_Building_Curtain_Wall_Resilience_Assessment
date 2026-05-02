package response

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"

	"icw_core_biz/pkg/rpc_err"
)

// RpcErrorMessages 业务错误枚举提示
var RpcErrorMessages = map[rpc_err.DetailCode]string{
	rpc_err.DetailBadRequest:                 "请求参数错误",
	rpc_err.DetailUnauthorized:               "登录状态已失效，请重新登录",
	rpc_err.DetailAccountLocked:              "登录失败次数过多，请稍后重试",
	rpc_err.DetailInternalError:              "服务暂时不可用，请稍后重试",
	rpc_err.DetailInvalidEmailAddress:        "邮箱格式错误",
	rpc_err.DetailEmailAlreadyRegistered:     "邮箱已被注册，请登录账号",
	rpc_err.DetailEmailNotRegistered:         "邮箱尚未注册，请注册账号",
	rpc_err.DetailEmailCodeSentTooFrequently: "验证码发送过于频繁，请稍后重试",
	rpc_err.DetailSendEmailCodeFailed:        "验证码发送失败，请稍后重试",
	rpc_err.DetailPasswordTooShortOrTooLong:  "密码长度必须不小于 8 位，不多于 24 位",
	rpc_err.DetailPasswordTooWeak:            "密码必须同时包含大小写英文字母、数字和符号",
	rpc_err.DetailNameRequired:               "请输入用户名称",
	rpc_err.DetailNameTooLong:                "用户名称不能超过 8 个字符",
	rpc_err.DetailIncorrectEmailCode:         "验证码错误",
	rpc_err.DetailInvalidCredentials:         "邮箱或密码错误",
	rpc_err.DetailInvalidAvatarContentType:   "请上传正确的图像格式",
	rpc_err.DetailProjectNotAccessible:       "无项目访问权限",
}

// WriteError 将 RPC 标准错误转换为 API 层的 HTTP 响应
func WriteError(c *gin.Context, err error) {
	code, detailCode, message := rpc_err.Parse(err)
	log.Printf("[ERROR] %s|%s: %s", code, detailCode, message)
	Error(c, errorStatus(code), errorCode(code, detailCode), errorMessage(detailCode))
}

// errorStatus 获取错误状态码
// 400 - BAD_REQUEST - 无效请求
// 401 - UNAUTHORIZED - 身份验证未通过
// 423 - ACCOUNT_LOCKED - 账号锁定（登录失败次数达上限）
// 500 - INTERNAL_ERROR - 服务器内部错误
func errorStatus(code rpc_err.Code) int {
	switch code {
	case rpc_err.CodeBadRequest:
		return http.StatusBadRequest
	case rpc_err.CodeUnauthorized:
		return http.StatusUnauthorized
	case rpc_err.CodeAccountLocked:
		return http.StatusLocked
	default:
		return http.StatusInternalServerError
	}
}

// errorCode 获取错误业务代码
func errorCode(code rpc_err.Code, detailCode rpc_err.DetailCode) string {
	if !rpc_err.IsCode(code) {
		code = rpc_err.DefaultCode()
	}
	if !rpc_err.IsDetailCode(detailCode) {
		detailCode = rpc_err.DefaultDetailCode(code)
	}
	return string(detailCode)
}

// errorMessage 获取错误业务信息
func errorMessage(detailCode rpc_err.DetailCode) string {
	if message, ok := RpcErrorMessages[detailCode]; ok && message != "" {
		return message
	}
	return RpcErrorMessages[rpc_err.DetailInternalError]
}
