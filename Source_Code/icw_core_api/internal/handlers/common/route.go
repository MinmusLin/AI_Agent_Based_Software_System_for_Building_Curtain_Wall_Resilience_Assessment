package common

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"

	"icw_common/utils"
)

// RouteGroup 封装 Gin 路由注册器
type RouteGroup struct {
	router       *gin.RouterGroup
	descriptions map[string]string
}

// NewRouteGroup 创建路由注册器
func NewRouteGroup(router *gin.RouterGroup, descriptions map[string]string) RouteGroup {
	return RouteGroup{
		router:       router,
		descriptions: descriptions,
	}
}

// handle 注册 Gin 路由
func (r RouteGroup) handle(method, path, description string, handlers ...gin.HandlerFunc) {
	if r.router == nil {
		return
	}
	r.router.Handle(method, path, handlers...)
	if r.descriptions != nil {
		r.descriptions[routeKey(method, joinPath(r.router.BasePath(), path))] = description
	}
}

// GET 注册 GET 路由
func (r RouteGroup) GET(path, description string, handlers ...gin.HandlerFunc) {
	r.handle(http.MethodGet, path, description, handlers...)
}

// POST 注册 POST 路由
func (r RouteGroup) POST(path, description string, handlers ...gin.HandlerFunc) {
	r.handle(http.MethodPost, path, description, handlers...)
}

// DELETE 注册 DELETE 路由
func (r RouteGroup) DELETE(path, description string, handlers ...gin.HandlerFunc) {
	r.handle(http.MethodDelete, path, description, handlers...)
}

// FormatRoutesTable 将 Gin 路由表格式化为表格
func FormatRoutesTable(routes gin.RoutesInfo, descriptions map[string]string) string {
	methodValues := make([]string, 0, len(routes))
	pathValues := make([]string, 0, len(routes))
	descriptionValues := make([]string, 0, len(routes))
	handlerValues := make([]string, 0, len(routes))
	for _, route := range routes {
		methodValues = append(methodValues, route.Method)
		pathValues = append(pathValues, route.Path)
		descriptionValues = append(descriptionValues, descriptions[routeKey(route.Method, route.Path)])
		handlerValues = append(handlerValues, strings.TrimSuffix(strings.TrimPrefix(route.Handler, "icw_core_api/internal/handlers/"), "-fm"))
	}
	return utils.FormatTable([]*utils.TableColumn{
		{
			Header: "method",
			Values: methodValues,
		},
		{
			Header: "path",
			Values: pathValues,
		},
		{
			Header: "description",
			Values: descriptionValues,
		},
		{
			Header: "handler",
			Values: handlerValues,
		},
	})
}

// routeKey 生成路由描述 Key
func routeKey(method, path string) string {
	return fmt.Sprintf("%s:%s", method, path)
}

// joinPath 拼接路由路径
func joinPath(basePath, path string) string {
	if basePath == "/" {
		basePath = ""
	}
	fullPath := strings.TrimRight(basePath, "/") + "/" + strings.TrimLeft(path, "/")
	if fullPath == "" {
		return "/"
	}
	return fullPath
}
