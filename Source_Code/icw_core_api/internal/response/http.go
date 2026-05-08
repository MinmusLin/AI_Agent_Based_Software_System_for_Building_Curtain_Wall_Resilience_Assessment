package response

import (
	"encoding/json"
	"net/http"

	"github.com/gin-gonic/gin"
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/proto"

	"icw_common/rpc/error"
	"icw_common/utils"
)

// okEnvelope HTTP 标准成功响应
type okEnvelope[T any] struct {
	Code    string `json:"code"`
	Message string `json:"message"`
	Data    T      `json:"data"`
}

// errorEnvelope HTTP 标准失败响应
type errorEnvelope struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

// OK 写入 HTTP 标准成功响应
func OK[T any](c *gin.Context, data T) {
	if protoData, ok := marshalProtoData(data); ok {
		c.JSON(http.StatusOK, okEnvelope[any]{
			Code:    "OK",
			Message: "success",
			Data:    protoData,
		})
		return
	}
	c.JSON(http.StatusOK, okEnvelope[T]{
		Code:    "OK",
		Message: "success",
		Data:    data,
	})
}

// Error 写入 HTTP 标准失败响应
func Error(c *gin.Context, err error) {
	code, detailCode, _ := rpc_error.Parse(err)
	c.JSON(errorStatus(code), errorEnvelope{
		Code:    errorCode(code, detailCode),
		Message: errorMessage(detailCode),
	})
}

// BindJSON 绑定 JSON 请求体
func BindJSON(c *gin.Context, out interface{}) bool {
	if err := c.ShouldBindJSON(out); err != nil {
		c.JSON(http.StatusBadRequest, errorEnvelope{
			Code:    string(rpc_error.DetailBadRequest),
			Message: errorMessage(rpc_error.DetailBadRequest),
		})
		return false
	}
	return true
}

// marshalProtoData 序列化协议数据
func marshalProtoData[T any](data T) (any, bool) {
	message, ok := any(data).(proto.Message)
	if !ok {
		return nil, false
	}
	if utils.IsNil(message) {
		return nil, true
	}
	bytes, err := protojson.MarshalOptions{
		EmitUnpopulated: true,
		UseEnumNumbers:  true,
		UseProtoNames:   true,
	}.Marshal(message)
	if err != nil {
		return nil, false
	}
	var payload any
	if err := json.Unmarshal(bytes, &payload); err != nil {
		return nil, false
	}
	return payload, true
}
