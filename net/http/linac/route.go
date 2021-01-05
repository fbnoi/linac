package linac

import (
	"net/http"
	xpath "path"
	"regexp"
	"strings"
)

// newRoute 添加路由处理方法
func newRoute(pattern, method string, handler ...Handler) *Route {
	regex, params := parseURI(pattern)
	return &Route{
		regex:    regex,
		params:   params,
		method:   method,
		handlers: handler,
	}
}

// RouteGroup Route集合
type RouteGroup struct {
	path     string
	handlers []Handler
	routes   []*Route
}

// AddRoute 向路由器中添加路由
func (group *RouteGroup) addRoute(path, method string, handler ...Handler) *RouteGroup {
	if path[0] != '/' {
		panic("pattern must start with '/'")
	}
	path = group.absPath(path)
	handler = group.mergeHandlers(handler...)
	group.routes = append(group.routes, newRoute(path, method, handler...))

	return group
}

// Use 向 Group 中添加全局的handler
// 该方法必须只用于添加中间件
// NOTE: 该方法添加的中间件，只对调用了 Use 方法之后注册的路由有效，而对之前注册的路由无效
func (group *RouteGroup) Use(handlers ...Handler) {
	group.handlers = group.mergeHandlers(handlers...)
}

// Group 新建分组
func (group *RouteGroup) Group(path string, register func(*RouteGroup) *RouteGroup, handlers ...Handler) {
	path = group.absPath(path)
	handlers = group.mergeHandlers(handlers...)
	newGroup := &RouteGroup{
		path:     path,
		handlers: handlers,
		routes:   []*Route{},
	}
	newGroup = register(newGroup)
	group.routes = append(group.routes, newGroup.routes...)
}

// GET 为一个路由注册一个GET方法
func (group *RouteGroup) GET(path string, handler ...Handler) *RouteGroup {
	return group.addRoute(path, "GET", handler...)
}

// POST 为一个路由注册一个POST方法
func (group *RouteGroup) POST(path string, handler ...Handler) *RouteGroup {
	return group.addRoute(path, "POST", handler...)
}

// PUT 为一个路由注册一个PUT方法
func (group *RouteGroup) PUT(path string, handler ...Handler) *RouteGroup {
	return group.addRoute(path, "PUT", handler...)
}

// DELETE 为一个路由注册一个DELETE方法
func (group *RouteGroup) DELETE(path string, handler ...Handler) *RouteGroup {
	return group.addRoute(path, "DELETE", handler...)
}

// HEAD 为一个路由注册一个HEAD方法
func (group *RouteGroup) HEAD(path string, handler ...Handler) *RouteGroup {
	return group.addRoute(path, "HEAD", handler...)
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

func (group *RouteGroup) mergeHandlers(handlers ...Handler) []Handler {
	return append(group.handlers, handlers...)
}

// Route model
type Route struct {
	regex    *regexp.Regexp
	method   string
	params   map[int]string
	handlers []Handler
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
	for i, handler := range route.handlers {
		if ctx.IsAbort() {
			return
		}
		if i == len(route.handlers) && ctx.Request.Method != route.method {
			ctx.Abort(http.StatusMethodNotAllowed)
		}
		handler(ctx)
	}
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
