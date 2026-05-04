package common

import (
	"reflect"
	"strings"
	"time"

	"icw_core_biz/consts"
	"icw_core_biz/utils"
)

// rpcLog 输出 RPC 请求日志
func rpcLog(method string, req, resp interface{}, start time.Time, err error) {
	requestId := getRequestId(req)
	if requestId == "" {
		requestId = "-"
	}

	if utils.IsEmptyError(err) {
		RpcInfo("[%s] [%s] cost=%s req=%s resp=%s", method, requestId, time.Since(start), utils.JSONF(req), utils.JSONF(resp))
		return
	}
	RpcError("[%s] [%s] cost=%s req=%s resp=%s err=%s", method, requestId, time.Since(start), utils.JSONF(req), utils.JSONF(resp), utils.FormatErrorLog(err))
}

// getRequestId 从 RPC 元数据中获取请求 ID
func getRequestId(req interface{}) string {
	if req == nil {
		return ""
	}

	reqValue := reflect.ValueOf(req)
	if reqValue.Kind() != reflect.Ptr || reqValue.IsNil() {
		return ""
	}

	elem := reqValue.Elem()
	if elem.Kind() != reflect.Struct {
		return ""
	}

	metaField := elem.FieldByName("Meta")
	if !metaField.IsValid() || metaField.Kind() != reflect.Ptr || metaField.IsNil() {
		return ""
	}

	metaValue := metaField.Elem()
	if metaValue.Kind() != reflect.Struct {
		return ""
	}

	requestIdField := metaValue.FieldByName("RequestId")
	if !requestIdField.IsValid() || requestIdField.Kind() != reflect.String {
		return ""
	}

	return strings.TrimSpace(requestIdField.String())
}

// RpcInfo 输出标准 RPC 日志
func RpcInfo(format string, args ...interface{}) {
	utils.LogInfo(consts.LogScopeRPC, consts.LogColorBoldGreen, format, args...)
}

// RpcWarn 输出警告 RPC 日志
func RpcWarn(format string, args ...interface{}) {
	utils.LogWarn(consts.LogScopeRPC, format, args...)
}

// RpcError 输出错误 RPC 日志
func RpcError(format string, args ...interface{}) {
	utils.LogError(consts.LogScopeRPC, format, args...)
}

// RpcFault 输出致命错误 RPC 日志并退出进程
func RpcFault(format string, args ...interface{}) {
	utils.LogFault(consts.LogScopeRPC, format, args...)
}
