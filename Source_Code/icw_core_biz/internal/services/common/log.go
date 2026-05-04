package common

import (
	"fmt"
	"log"
	"reflect"
	"strings"
	"time"

	"icw_core_biz/utils"
)

const (
	// LogColorReset ANSI 终端颜色重置码
	LogColorReset = "\033[0m"
	// LogColorBoldGreen ANSI 终端颜色码：绿色
	LogColorBoldGreen = "\033[1;32m"
	// LogColorBoldRed ANSI 终端颜色码：红色
	LogColorBoldRed = "\033[1;31m"
	// LogColorBoldPurple ANSI 终端颜色码：紫色
	LogColorBoldPurple = "\033[1;35m"
	// LogColorBoldYellow ANSI 终端颜色码：黄色
	LogColorBoldYellow = "\033[1;33m"
)

// rpcLog 记录 RPC 请求日志
func rpcLog(method string, req, resp interface{}, start time.Time, err error) {
	requestId := getRequestId(req)
	if requestId == "" {
		requestId = "-"
	}

	if isEmptyError(err) {
		log.Printf("%s [%s] [%s] cost=%s req=%s resp=%s", RpcInfoAndErrorPrefix(err), method, requestId, time.Since(start), utils.JSONF(req), utils.JSONF(resp))
		return
	}
	log.Printf("%s [%s] [%s] cost=%s req=%s resp=%s err=%s", RpcInfoAndErrorPrefix(err), method, requestId, time.Since(start), utils.JSONF(req), utils.JSONF(resp), formatError(err))
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

// RpcInfoAndErrorPrefix RPC 正常和错误日志前缀
func RpcInfoAndErrorPrefix(err interface{}) string {
	if isEmptyError(err) {
		return LogColorBoldGreen + "[RPC INFO]" + LogColorReset
	}
	return LogColorBoldRed + "[RPC ERROR]" + LogColorReset
}

// RpcWarnPrefix RPC 警告日志前缀
func RpcWarnPrefix() string {
	return LogColorBoldYellow + "[WARN]" + LogColorReset
}

// formatError 格式化错误日志
func formatError(err interface{}) string {
	msg := strings.TrimSpace(fmt.Sprint(err))
	if isEmptyError(err) {
		return msg
	}
	return LogColorBoldPurple + msg + LogColorReset
}

// isEmptyError 判断错误是否为空错误
func isEmptyError(err interface{}) bool {
	msg := strings.TrimSpace(fmt.Sprint(err))
	return msg == "" || msg == "nil" || msg == "<nil>"
}
