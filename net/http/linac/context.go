package linac

import "net/http"

// Context http 请求上下文
type Context struct {
	Writer  http.ResponseWriter
	Request *http.Request
	Params  map[string]string

	abort bool
}

// Abort 终止http响应，并设置http code
func (ctx *Context) Abort(code int) {
	ctx.Writer.WriteHeader(code)
	ctx.abort = true
}

// IsAbort 返回context是否终止响应
func (ctx *Context) IsAbort() bool {
	return ctx.abort
}

// Handler http 请求处理
type Handler func(*Context)
