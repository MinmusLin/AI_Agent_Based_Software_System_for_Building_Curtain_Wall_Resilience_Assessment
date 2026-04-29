package rpc_err

// Code 错误状态代码
type Code string

const (
	CodeBadRequest    Code = "BAD_REQUEST"
	CodeUnauthorized  Code = "UNAUTHORIZED"
	CodeAccountLocked Code = "ACCOUNT_LOCKED"
	CodeInternalError Code = "INTERNAL_ERROR"
)

// IsCode 判断是否为已定义错误状态代码
func IsCode(code Code) bool {
	switch code {
	case CodeBadRequest, CodeUnauthorized, CodeAccountLocked, CodeInternalError:
		return true
	default:
		return false
	}
}

// DefaultCode 返回默认错误状态代码
func DefaultCode() Code {
	return CodeInternalError
}

// DetailCode 错误业务代码
type DetailCode string

const (
	DetailBadRequest                 DetailCode = "BAD_REQUEST"
	DetailUnauthorized               DetailCode = "UNAUTHORIZED"
	DetailAccountLocked              DetailCode = "ACCOUNT_LOCKED"
	DetailInternalError              DetailCode = "INTERNAL_ERROR"
	DetailInvalidEmailAddress        DetailCode = "INVALID_EMAIL_ADDRESS"
	DetailEmailAlreadyRegistered     DetailCode = "EMAIL_ALREADY_REGISTERED"
	DetailEmailNotRegistered         DetailCode = "EMAIL_NOT_REGISTERED"
	DetailEmailCodeSentTooFrequently DetailCode = "EMAIL_CODE_SENT_TOO_FREQUENTLY"
	DetailSendEmailCodeFailed        DetailCode = "SEND_EMAIL_CODE_FAILED"
	DetailPasswordTooShort           DetailCode = "PASSWORD_TOO_SHORT"
	DetailPasswordTooWeak            DetailCode = "PASSWORD_TOO_WEAK"
	DetailNameRequired               DetailCode = "NAME_REQUIRED"
	DetailNameTooLong                DetailCode = "NAME_TOO_LONG"
	DetailIncorrectEmailCode         DetailCode = "INCORRECT_EMAIL_CODE"
	DetailInvalidCredentials         DetailCode = "INVALID_CREDENTIALS"
)

// IsDetailCode 判断是否为已定义错误业务代码
func IsDetailCode(detailCode DetailCode) bool {
	switch detailCode {
	case DetailBadRequest,
		DetailUnauthorized,
		DetailAccountLocked,
		DetailInternalError,
		DetailInvalidEmailAddress,
		DetailEmailAlreadyRegistered,
		DetailEmailNotRegistered,
		DetailEmailCodeSentTooFrequently,
		DetailSendEmailCodeFailed,
		DetailPasswordTooShort,
		DetailPasswordTooWeak,
		DetailNameRequired,
		DetailNameTooLong,
		DetailIncorrectEmailCode,
		DetailInvalidCredentials:
		return true
	default:
		return false
	}
}

// DefaultDetailCode 返回错误状态代码对应的默认错误业务代码
func DefaultDetailCode(code Code) DetailCode {
	switch code {
	case CodeBadRequest:
		return DetailBadRequest
	case CodeUnauthorized:
		return DetailUnauthorized
	case CodeAccountLocked:
		return DetailAccountLocked
	default:
		return DetailInternalError
	}
}
