package linac

import (
	"regexp"
	"strings"
)

// newRoute 添加路由处理方法
func newRoute(pattern, method string, handler Handler) *Route {
	regex, params := parseURI(pattern)
	return &Route{
		Regex:   regex,
		Params:  params,
		Method:  method,
		Handler: handler,
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

// Route model
type Route struct {
	Regex   *regexp.Regexp
	Method  string
	Params  map[int]string
	Handler Handler
}
