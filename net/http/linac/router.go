package linac

import (
	"context"
	"fmt"
	"net/http"
	"strings"
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
	if route, ok := router.metchRoute(ctx); ok {
		var (
			cancel         func()
			tm             time.Duration
			maxRequestBody int64
		)
		conf := router.engine.GetConfig()
		tm = conf.Timeout
		maxRequestBody = conf.MaxRequestBody
		if conf, ok := route.GetConfig(); ok {
			tm = conf.Timeout
			maxRequestBody = conf.MaxRequestBody
		}

		req := ctx.Request
		ctype := req.Header.Get("Content-Type")
		switch {
		case strings.Contains(ctype, "multipart/form-data"):
			req.ParseMultipartForm(maxRequestBody)
		default:
			req.ParseForm()
		}

		c := context.Background()
		if tm > 0 {
			ctx.Context, cancel = context.WithTimeout(c, tm)
		} else {
			ctx.Context, cancel = context.WithCancel(c)
		}
		defer cancel()
		route.handle(ctx)
	} else {
		router.getNotFoundHandler()(ctx)
	}
}

// metchRoute 匹配context路由并返回
func (router *Router) metchRoute(ctx *Context) (route *Route, ok bool) {
	url := ctx.Request.URL
	for _, cond := range router.RouteGroup.routes {
		if !cond.regex.MatchString(url.Path) {
			continue
		}
		matches := cond.regex.FindStringSubmatch(url.Path)
		//再次检测是否匹配
		if len(matches[0]) != len(url.Path) {
			continue
		}
		// 如果路由模式匹配，检测 http 请求方法
		if ctx.Request.Method != cond.method {
			continue
		}
		route, ok = cond, true
	}
	return
}

//默认 not found handler，返回404状态码
func defaultNotFoundHandler(context *Context) {
	context.String(http.StatusNotFound, fmt.Sprintf("no route found for %s:%s", context.Request.Method, context.Request.URL))
}
