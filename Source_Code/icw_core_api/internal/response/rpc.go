package response

import (
	"net/http"

	"icw_common/rpc/error"
)

// rpcErrorMessages 业务错误枚举提示
var rpcErrorMessages = map[rpc_error.DetailCode]string{
	rpc_error.DetailBadRequest:                              "请求参数错误",
	rpc_error.DetailUnauthorized:                            "登录状态已失效，请重新登录",
	rpc_error.DetailAccountLocked:                           "登录失败次数过多，请稍后重试",
	rpc_error.DetailInternalError:                           "服务暂时不可用，请稍后重试",
	rpc_error.DetailInvalidEmailAddress:                     "邮箱格式错误",
	rpc_error.DetailEmailAlreadyRegistered:                  "邮箱已被注册，请登录账号",
	rpc_error.DetailEmailNotRegistered:                      "邮箱尚未注册，请注册账号",
	rpc_error.DetailEmailCodeSentTooFrequently:              "验证码发送过于频繁，请稍后重试",
	rpc_error.DetailSendEmailCodeFailed:                     "验证码发送失败，请稍后重试",
	rpc_error.DetailPasswordTooShortOrTooLong:               "密码长度必须不小于 8 位，不多于 24 位",
	rpc_error.DetailPasswordTooWeak:                         "密码必须同时包含大小写英文字母、数字和符号",
	rpc_error.DetailNameRequired:                            "请输入用户名称",
	rpc_error.DetailNameTooLong:                             "用户名称不能超过 8 个字符",
	rpc_error.DetailIncorrectEmailCode:                      "验证码错误",
	rpc_error.DetailInvalidCredentials:                      "邮箱或密码错误",
	rpc_error.DetailInvalidImageContentType:                 "请上传正确的图像格式",
	rpc_error.DetailProjectNotAccessible:                    "无项目访问权限",
	rpc_error.DetailProjectNameTooLong:                      "项目名称不能超过 32 个字符",
	rpc_error.DetailProjectBuildingNameTooLong:              "建筑名称不能超过 32 个字符",
	rpc_error.DetailProjectBuildingLocationTooLong:          "建筑地址不能超过 128 个字符",
	rpc_error.DetailProjectBuildingDescriptionTooLong:       "建筑描述不能超过 5000 个字符",
	rpc_error.DetailProjectKnownIssuesTooLong:               "已知问题或人工先验描述不能超过 5000 个字符",
	rpc_error.DetailProjectAssessmentGoalTooLong:            "评估目标或重点关注方向不能超过 5000 个字符",
	rpc_error.DetailProjectAtLeastOneGroupRequired:          "项目应至少存在一个图像组",
	rpc_error.DetailProjectGroupNameRequired:                "图像组名称不能为空",
	rpc_error.DetailProjectGroupNameTooLong:                 "图像组名称不能超过 32 个字符",
	rpc_error.DetailProjectGroupNameDuplicated:              "图像组名称已存在",
	rpc_error.DetailProjectImageFileNameTooLong:             "图像文件名不能超过 255 个字符",
	rpc_error.DetailProjectImageFormatInvalid:               "图像格式不合法",
	rpc_error.DetailProjectImageExpired:                     "图片已失效，请重新上传",
	rpc_error.DetailProjectUploadedImageCountRequired:       "项目应至少存在一张已上传图像",
	rpc_error.DetailProjectEmptyGroupCountInvalid:           "项目中不能存在空图像组",
	rpc_error.DetailProjectPendingOrFailedImageCountInvalid: "项目中不能存在上传中或上传失败的图像",
	rpc_error.DetailProjectDetectionTasksNotSucceeded:       "项目图像检测任务尚未全部完成",
}

// errorStatus 获取错误状态码
// 400 - BAD_REQUEST - 无效请求
// 401 - UNAUTHORIZED - 身份验证未通过
// 423 - ACCOUNT_LOCKED - 账号锁定（登录失败次数达上限）
// 500 - INTERNAL_ERROR - 服务器内部错误
func errorStatus(code rpc_error.Code) int {
	switch code {
	case rpc_error.CodeBadRequest:
		return http.StatusBadRequest
	case rpc_error.CodeUnauthorized:
		return http.StatusUnauthorized
	case rpc_error.CodeAccountLocked:
		return http.StatusLocked
	default:
		return http.StatusInternalServerError
	}
}

// errorCode 获取错误业务代码
func errorCode(code rpc_error.Code, detailCode rpc_error.DetailCode) string {
	if !rpc_error.IsCode(code) {
		code = rpc_error.DefaultCode()
	}
	if !rpc_error.IsDetailCode(detailCode) {
		detailCode = rpc_error.DefaultDetailCode(code)
	}
	return string(detailCode)
}

// errorMessage 获取错误业务信息
func errorMessage(detailCode rpc_error.DetailCode) string {
	if message, ok := rpcErrorMessages[detailCode]; ok && message != "" {
		return message
	}
	return rpcErrorMessages[rpc_error.DetailInternalError]
}
