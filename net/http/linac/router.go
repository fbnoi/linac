package linac

import (
	"net/http"
)

// Router model
type Router struct {
	routes []*Route
	engine *Engine
}

// AddRoute 向路由器中添加路由
func (router *Router) addRoute(path, method string, handler Handler) *Router {
	if path[0] != '/' {
		panic("pattern must start with '/'")
	}
	router.routes = append(router.routes, newRoute(path, method, handler))
	return router
}

// GET 为一个路由注册一个GET方法
func (router *Router) GET(path string, handler Handler) *Router {
	return router.addRoute(path, "GET", handler)
}

// POST 为一个路由注册一个POST方法
func (router *Router) POST(path string, handler Handler) *Router {
	return router.addRoute(path, "POST", handler)
}

// PUT 为一个路由注册一个PUT方法
func (router *Router) PUT(path string, handler Handler) *Router {
	return router.addRoute(path, "PUT", handler)
}

// DELETE 为一个路由注册一个DELETE方法
func (router *Router) DELETE(path string, handler Handler) *Router {
	return router.addRoute(path, "DELETE", handler)
}

// HEAD 为一个路由注册一个HEAD方法
func (router *Router) HEAD(path string, handler Handler) *Router {
	return router.addRoute(path, "HEAD", handler)
}

// ServeHTTP 响应http请求 此处进行context内容的生成
func (router *Router) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	context := &Context{
		Writer:  w,
		Request: r,
	}
	route := router.metchRoute(context)
	route.handle(context)
}

// metchRoute 匹配context路由并返回
func (router *Router) metchRoute(ctx *Context) *Route {
	r := ctx.Request
	var route *Route
	for _, route = range router.routes {
		if !route.Regex.MatchString(r.RequestURI) {
			continue
		}
		matches := route.Regex.FindStringSubmatch(r.RequestURI)
		//double check that the Route matches the URL pattern.
		if len(matches[0]) != len(r.RequestURI) {
			continue
		}
		return route
	}
	return route
}
