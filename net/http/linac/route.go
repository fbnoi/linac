package linac

import (
	"fmt"
	xpath "path"
	"regexp"
	"strings"
	"sync/atomic"
	"time"
)

// newRoute 添加路由处理方法
func newRoute(pattern, method, name string, handler ...Handler) *Route {
	regex, params := parseURI(pattern)
	return &Route{
		name:     name,
		regex:    regex,
		params:   params,
		method:   method,
		handlers: handler,
	}
}

// RouteGroup Route集合
type RouteGroup struct {
	name     string
	path     string
	handlers []Handler
	routes   map[string]*Route
}

// AddRoute 向路由器中添加路由
func (group *RouteGroup) addRoute(path, method, name string, handler ...Handler) (route *Route) {
	if path[0] != '/' {
		panic("pattern must start with '/'")
	}
	name = group.fullName(name)
	path = group.absPath(path)
	handler = group.mergeHandlers(handler...)
	if _, ok := group.GetRoute(name); ok {
		panic(fmt.Errorf("add route error, name '%s' already exist", name))
	}
	route = newRoute(path, method, name, handler...)
	group.routes[name] = route
	return
}

// Use 向 Group 中添加全局的handler
// 该方法必须只用于添加中间件
// NOTE: 该方法添加的中间件，只对调用了 Use 方法之后注册的路由有效，而对之前注册的路由无效
func (group *RouteGroup) Use(handlers ...Handler) {
	group.handlers = group.mergeHandlers(handlers...)
}

// Group 新建分组
func (group *RouteGroup) Group(path, name string, register func(*RouteGroup) *RouteGroup, handlers ...Handler) {
	name = group.fullName(name)
	path = group.absPath(path)
	handlers = group.mergeHandlers(handlers...)
	newGroup := &RouteGroup{
		name:     name,
		path:     path,
		handlers: handlers,
		routes:   make(map[string]*Route),
	}
	newGroup = register(newGroup)
	for name, route := range newGroup.routes {
		if _, ok := group.GetRoute(name); ok {
			panic(fmt.Errorf("add route error, name '%s' already exist", name))
		}
		group.routes[name] = route
	}
}

// GET 为一个路由注册一个GET方法
func (group *RouteGroup) GET(path, name string, handler ...Handler) *Route {
	return group.addRoute(path, "GET", name, handler...)
}

// POST 为一个路由注册一个POST方法
func (group *RouteGroup) POST(path, name string, handler ...Handler) *Route {
	return group.addRoute(path, "POST", name, handler...)
}

// PUT 为一个路由注册一个PUT方法
func (group *RouteGroup) PUT(path, name string, handler ...Handler) *Route {
	return group.addRoute(path, "PUT", name, handler...)
}

// DELETE 为一个路由注册一个DELETE方法
func (group *RouteGroup) DELETE(path, name string, handler ...Handler) *Route {
	return group.addRoute(path, "DELETE", name, handler...)
}

// HEAD 为一个路由注册一个HEAD方法
func (group *RouteGroup) HEAD(path, name string, handler ...Handler) *Route {
	return group.addRoute(path, "HEAD", name, handler...)
}

//GetRoute 获取route
func (group *RouteGroup) GetRoute(name string) (route *Route, ok bool) {
	route, ok = group.routes[name]
	return
}

func (group *RouteGroup) absPath(path string) string {
	if path == "" {
		return group.path
	}
	finalPath := xpath.Join(group.path, path)
	appendSlash := path[len(path)-1] == '/' && finalPath[len(finalPath)-1] != '/'
	if appendSlash {
		return finalPath + "/"
	}
	return finalPath
}

func (group *RouteGroup) fullName(name string) string {
	if group.name == "" {
		return name
	}
	return group.name + "." + name
}

func (group *RouteGroup) mergeHandlers(handlers ...Handler) []Handler {
	return append(group.handlers, handlers...)
}

// RouteConfig 路由配置
// 为路由定制配置选项
type RouteConfig struct {
	Timeout   time.Duration
	MaxMemery int
}

// Route model
type Route struct {
	name     string
	regex    *regexp.Regexp
	method   string
	params   map[int]string
	handlers []Handler

	config *atomic.Value
}

// SetConfig 为路由添加特定的配置
func (route *Route) SetConfig(config *RouteConfig) {
	route.config.Store(config)
}

// GetConfig 获取路由配置
func (route *Route) GetConfig() (config *RouteConfig, ok bool) {
	config, ok = route.config.Load().(*RouteConfig)
	return
}

// handle 处理http请求
// 1.解析路由参数
// 2.调用 Handler 处理 context
func (route *Route) handle(ctx *Context) {
	r := ctx.Request
	matches := route.regex.FindStringSubmatch(r.RequestURI)
	params := make(map[string]interface{})
	if len(route.params) > 0 {
		for i, match := range matches[1:] {
			params[route.params[i]] = match
		}
	}
	ctx.Params = params
	ctx.Handlers = route.handlers
	ctx.Next()
}

// pattern 路由模式
// 如：'/users' 或者 '/users/:id'
// 其中 :id 将被解析为路由参数。也可以为参数添加正则验证，
// 如：'/user/:id([0-9]+)'
func parseURI(pattern string) (*regexp.Regexp, map[int]string) {
	parts := strings.Split(pattern, "/")
	params := make(map[int]string)
	j := 0
	for index, part := range parts {
		if strings.HasPrefix(part, ":") {
			expr := "([^/]+)"
			if i := strings.Index(part, "("); i != -1 {
				expr = part[i:]
				part = part[:i]
			}
			parts[index] = expr
			params[j] = part
			j++
		}
	}

	pattern = strings.Join(parts, "/")
	regex, regexErr := regexp.Compile(pattern)
	if regexErr != nil {
		panic(regexErr)
	}
	return regex, params
}
