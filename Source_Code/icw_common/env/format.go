package env

import (
	"fmt"
	"reflect"
	"strings"

	"icw_common/utils"
)

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
	valueType := value.Type()
	nameValues := make([]string, 0, value.NumField())
	configValues := make([]string, 0, value.NumField())
	for i := 0; i < value.NumField(); i++ {
		field := valueType.Field(i)
		if field.PkgPath != "" {
			continue
		}
		name := envConfigName(field)
		if name == "" {
			continue
		}
		nameValues = append(nameValues, name)
		configValues = append(configValues, envConfigValue(value.Field(i)))
	}
	return utils.FormatTable([]*utils.TableColumn{
		{
			Header: "env",
			Values: nameValues,
		},
		{
			Header: "value",
			Values: configValues,
		},
	})
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
