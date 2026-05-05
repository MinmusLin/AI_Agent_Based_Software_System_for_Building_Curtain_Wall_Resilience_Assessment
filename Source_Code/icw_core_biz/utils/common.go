package utils

import (
	"encoding/json"
	"fmt"
	"reflect"
	"strings"
	"unicode"

	"icw_core_biz/consts"
)

// JSONF 将任意结构格式化为 JSON 字符串
func JSONF(v interface{}) string {
	bytes, err := json.Marshal(v)
	if err != nil {
		return ""
	}
	return string(bytes)
}

// IsEmptyError 判断错误是否为空错误
func IsEmptyError(err interface{}) bool {
	msg := strings.TrimSpace(fmt.Sprint(err))
	return msg == "" || msg == "nil" || msg == "<nil>"
}

// FormatErrorLog 格式化错误日志内容
func FormatErrorLog(err interface{}) string {
	msg := strings.TrimSpace(fmt.Sprint(err))
	if IsEmptyError(err) {
		return msg
	}
	return consts.LogColorBoldPurple + msg + consts.LogColorReset
}

// FormatEnvConfig 将已加载的环境变量配置格式化为标准输出
func FormatEnvConfig(config interface{}) string {
	value := reflect.ValueOf(config)
	if !value.IsValid() {
		return ""
	}
	if value.Kind() == reflect.Ptr {
		if value.IsNil() {
			return ""
		}
		value = value.Elem()
	}
	if value.Kind() != reflect.Struct {
		return ""
	}
	type envConfigItem struct {
		Name  string
		Value string
	}
	valueType := value.Type()
	items := make([]envConfigItem, 0, value.NumField())
	maxNameWidth := 0
	for i := 0; i < value.NumField(); i++ {
		field := valueType.Field(i)
		if field.PkgPath != "" {
			continue
		}
		name := envConfigName(field)
		if name == "" {
			continue
		}
		items = append(items, envConfigItem{
			Name:  name,
			Value: envConfigValue(value.Field(i)),
		})
		maxNameWidth = max(maxNameWidth, len(name))
	}
	var builder strings.Builder
	for index, item := range items {
		if index > 0 {
			builder.WriteString("\n")
		}
		builder.WriteString(fmt.Sprintf("%-*s = %s", maxNameWidth, item.Name, item.Value))
	}
	return builder.String()
}

// envConfigName 标准化环境变量名
func envConfigName(field reflect.StructField) string {
	tag := strings.TrimSpace(field.Tag.Get("env"))
	if tag == "-" {
		return ""
	}
	name, _, _ := strings.Cut(tag, ",")
	name = strings.TrimSpace(name)
	if name == "" {
		return field.Name
	}
	return name
}

// envConfigValue 标准化环境变量值
func envConfigValue(value reflect.Value) string {
	if !value.IsValid() || !value.CanInterface() {
		return ""
	}
	return fmt.Sprint(value.Interface())
}

// PadRight 按终端显示宽度补齐右侧空格
func PadRight(value string, width int) string {
	padding := width - DisplayWidth(value)
	if padding <= 0 {
		return value
	}
	return value + strings.Repeat(" ", padding)
}

// DisplayWidth 计算终端显示宽度
func DisplayWidth(value string) int {
	width := 0
	for _, r := range value {
		if unicode.Is(unicode.Han, r) || unicode.Is(unicode.Hiragana, r) ||
			unicode.Is(unicode.Katakana, r) || unicode.Is(unicode.Hangul, r) ||
			(r >= 0xFF01 && r <= 0xFF60) || (r >= 0xFFE0 && r <= 0xFFE6) {
			width += 2
			continue
		}
		width++
	}
	return width
}
