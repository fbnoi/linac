package linac

import (
	"net/http"
	"regexp"
	"strings"
)

// newRoute 添加路由处理方法
func newRoute(pattern, method string, handler ...Handler) *Route {
	regex, params := parseURI(pattern)
	return &Route{
		Regex:    regex,
		Params:   params,
		Method:   method,
		Handlers: handler,
	}
}

// Route model
type Route struct {
	Regex    *regexp.Regexp
	Method   string
	Params   map[int]string
	Handlers []Handler
}

// handle 处理http请求
// 1.解析路由参数
// 2.调用 Handler 处理 context
func (route *Route) handle(ctx *Context) {
	r := ctx.Request
	matches := route.Regex.FindStringSubmatch(r.RequestURI)
	params := make(map[string]string)
	if len(route.Params) > 0 {
		for i, match := range matches[1:] {
			params[route.Params[i]] = match
		}
	}
	ctx.Params = params
	for i, handler := range route.Handlers {
		if ctx.IsAbort() {
			return
		}
		if i == len(route.Handlers) && ctx.Request.Method != route.Method {
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
