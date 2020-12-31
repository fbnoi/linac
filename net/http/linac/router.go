package linac

import (
	"net/http"
)

// Router model
type Router struct {
	routes []*Route
	engine *Engine

	notFoundHandler Handler
}

// AddRoute 向路由器中添加路由
func (router *Router) addRoute(path, method string, handler ...Handler) *Router {
	if path[0] != '/' {
		panic("pattern must start with '/'")
	}
	router.routes = append(router.routes, newRoute(path, method, handler...))
	return router
}

// GET 为一个路由注册一个GET方法
func (router *Router) GET(path string, handler ...Handler) *Router {
	return router.addRoute(path, "GET", handler...)
}

// POST 为一个路由注册一个POST方法
func (router *Router) POST(path string, handler ...Handler) *Router {
	return router.addRoute(path, "POST", handler...)
}

// PUT 为一个路由注册一个PUT方法
func (router *Router) PUT(path string, handler ...Handler) *Router {
	return router.addRoute(path, "PUT", handler...)
}

// DELETE 为一个路由注册一个DELETE方法
func (router *Router) DELETE(path string, handler ...Handler) *Router {
	return router.addRoute(path, "DELETE", handler...)
}

// HEAD 为一个路由注册一个HEAD方法
func (router *Router) HEAD(path string, handler ...Handler) *Router {
	return router.addRoute(path, "HEAD", handler...)
}

// SetNotFoundHandler 设置默认 404 handler
func (router *Router) SetNotFoundHandler(handler Handler) *Router {
	router.notFoundHandler = handler
	return router
}

func (router *Router) getNotFoundHandler() Handler {
	if router.notFoundHandler == nil {
		router.notFoundHandler = defaultNotFoundHandler
	}
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
	for _, cond := range router.routes {
		if !cond.Regex.MatchString(r.RequestURI) {
			continue
		}
		matches := cond.Regex.FindStringSubmatch(r.RequestURI)
		//double check that the Route matches the URL pattern.
		if len(matches[0]) != len(r.RequestURI) {
			continue
		}
		route, ok = cond, true
		if r.Method == route.Method {
			return
		}
	}
	return
}

func defaultNotFoundHandler(context *Context) {

}
