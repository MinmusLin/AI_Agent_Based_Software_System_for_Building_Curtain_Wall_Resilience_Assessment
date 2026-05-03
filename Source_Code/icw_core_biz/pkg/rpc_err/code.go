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
	DetailBadRequest                              DetailCode = "BAD_REQUEST"
	DetailUnauthorized                            DetailCode = "UNAUTHORIZED"
	DetailAccountLocked                           DetailCode = "ACCOUNT_LOCKED"
	DetailInternalError                           DetailCode = "INTERNAL_ERROR"
	DetailInvalidEmailAddress                     DetailCode = "INVALID_EMAIL_ADDRESS"
	DetailEmailAlreadyRegistered                  DetailCode = "EMAIL_ALREADY_REGISTERED"
	DetailEmailNotRegistered                      DetailCode = "EMAIL_NOT_REGISTERED"
	DetailEmailCodeSentTooFrequently              DetailCode = "EMAIL_CODE_SENT_TOO_FREQUENTLY"
	DetailSendEmailCodeFailed                     DetailCode = "SEND_EMAIL_CODE_FAILED"
	DetailPasswordTooShortOrTooLong               DetailCode = "PASSWORD_TOO_SHORT_OR_TOO_LONG"
	DetailPasswordTooWeak                         DetailCode = "PASSWORD_TOO_WEAK"
	DetailNameRequired                            DetailCode = "NAME_REQUIRED"
	DetailNameTooLong                             DetailCode = "NAME_TOO_LONG"
	DetailIncorrectEmailCode                      DetailCode = "INCORRECT_EMAIL_CODE"
	DetailInvalidCredentials                      DetailCode = "INVALID_CREDENTIALS"
	DetailInvalidImageContentType                 DetailCode = "INVALID_IMAGE_CONTENT_TYPE"
	DetailProjectNotAccessible                    DetailCode = "PROJECT_NOT_ACCESSIBLE"
	DetailProjectNameTooLong                      DetailCode = "PROJECT_NAME_TOO_LONG"
	DetailProjectBuildingNameTooLong              DetailCode = "PROJECT_BUILDING_NAME_TOO_LONG"
	DetailProjectBuildingLocationTooLong          DetailCode = "PROJECT_BUILDING_LOCATION_TOO_LONG"
	DetailProjectBuildingDescriptionTooLong       DetailCode = "PROJECT_BUILDING_DESCRIPTION_TOO_LONG"
	DetailProjectKnownIssuesTooLong               DetailCode = "PROJECT_KNOWN_ISSUES_TOO_LONG"
	DetailProjectAssessmentGoalTooLong            DetailCode = "PROJECT_ASSESSMENT_GOAL_TOO_LONG"
	DetailProjectAtLeastOneGroupRequired          DetailCode = "PROJECT_AT_LEAST_ONE_GROUP_REQUIRED"
	DetailProjectGroupNameRequired                DetailCode = "PROJECT_GROUP_NAME_REQUIRED"
	DetailProjectGroupNameTooLong                 DetailCode = "PROJECT_GROUP_NAME_TOO_LONG"
	DetailProjectGroupNameDuplicated              DetailCode = "PROJECT_GROUP_NAME_DUPLICATED"
	DetailProjectImageFileNameTooLong             DetailCode = "PROJECT_IMAGE_FILE_NAME_TOO_LONG"
	DetailProjectImageFormatInvalid               DetailCode = "PROJECT_IMAGE_FORMAT_INVALID"
	DetailProjectUploadedImageCountRequired       DetailCode = "PROJECT_UPLOADED_IMAGE_COUNT_REQUIRED"
	DetailProjectEmptyGroupCountInvalid           DetailCode = "PROJECT_EMPTY_GROUP_COUNT_INVALID"
	DetailProjectPendingOrFailedImageCountInvalid DetailCode = "PROJECT_PENDING_OR_FAILED_IMAGE_COUNT_INVALID"
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
		DetailPasswordTooShortOrTooLong,
		DetailPasswordTooWeak,
		DetailNameRequired,
		DetailNameTooLong,
		DetailIncorrectEmailCode,
		DetailInvalidCredentials,
		DetailInvalidImageContentType,
		DetailProjectNotAccessible,
		DetailProjectNameTooLong,
		DetailProjectBuildingNameTooLong,
		DetailProjectBuildingLocationTooLong,
		DetailProjectBuildingDescriptionTooLong,
		DetailProjectKnownIssuesTooLong,
		DetailProjectAssessmentGoalTooLong,
		DetailProjectAtLeastOneGroupRequired,
		DetailProjectGroupNameRequired,
		DetailProjectGroupNameTooLong,
		DetailProjectGroupNameDuplicated,
		DetailProjectImageFileNameTooLong,
		DetailProjectImageFormatInvalid,
		DetailProjectUploadedImageCountRequired,
		DetailProjectEmptyGroupCountInvalid,
		DetailProjectPendingOrFailedImageCountInvalid:
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
