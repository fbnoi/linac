package linac

import (
	"linac/net/http/linac/render"
	"net/http"
)

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

// JSON  返回json response
func (ctx *Context) JSON(data interface{}, code int) {
	render.Write(&render.JSON{
		Code: code,
		Data: data,
		Err:  "",
	}, ctx.Writer)
}

// Handler http 请求处理
type Handler func(*Context)
