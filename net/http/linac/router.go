package linac

import (
	"context"
	"net/http"
	"time"
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

// ServeHTTP 响应http请求
func (router *Router) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	router.handleContext(&Context{
		Writer:  w,
		Request: r,
		index:   -1,
		abort:   false,
	})
}

// handleContext 处理context, 添加超时
func (router *Router) handleContext(ctx *Context) {
	var (
		cancel func()
		tm     time.Duration
	)
	conf, ok := router.engine.GetConfig()
	if ok {
		tm = time.Duration(conf.Timeout)
	}
	if tm > 0 {
		ctx.Context, cancel = context.WithTimeout(context.Background(), tm)
	} else {
		ctx.Context, cancel = context.WithCancel(context.Background())
	}

	defer cancel()
	if route, ok := router.metchRoute(ctx); ok {
		route.handle(ctx)
	} else {
		router.getNotFoundHandler()(ctx)
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
