package linac

import (
	"net/http"
)

// NewRouter 返回一个router
func NewRouter() *Router {
	return &Router{
		RouteGroup: &RouteGroup{
			path: "/",
		},
		notFoundHandler: defaultNotFoundHandler,
	}
}

// Router model
type Router struct {
	*RouteGroup
	engine *Engine

	notFoundHandler Handler
}

// SetNotFoundHandler 设置默认 404 handler
func (router *Router) SetNotFoundHandler(handler Handler) *Router {
	router.notFoundHandler = handler
	return router
}

func (router *Router) getNotFoundHandler() Handler {
	return router.notFoundHandler
}

// ServeHTTP 响应http请求 此处进行context内容的生成
func (router *Router) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	context := &Context{
		Writer:  w,
		Request: r,
	}
	if route, ok := router.metchRoute(context); ok {
		route.handle(context)
	} else {
		router.getNotFoundHandler()(context)
	}
}

// metchRoute 匹配context路由并返回
func (router *Router) metchRoute(ctx *Context) (route *Route, ok bool) {
	r := ctx.Request
	for _, cond := range router.RouteGroup.routes {
		if !cond.regex.MatchString(r.RequestURI) {
			continue
		}
		matches := cond.regex.FindStringSubmatch(r.RequestURI)
		//再次检测是否匹配
		if len(matches[0]) != len(r.RequestURI) {
			continue
		}
		route, ok = cond, true
		// 如果路由模式匹配，并且http方法相同，即刻返回
		if r.Method == route.method {
			return
		}
	}
	return
}

//默认 not found handler，返回404状态码
func defaultNotFoundHandler(context *Context) {
	context.Abort(http.StatusNotFound)
}
