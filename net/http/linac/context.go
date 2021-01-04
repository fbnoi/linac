package linac

import (
	xerror "linac/error"
	"linac/net/http/linac/render"
	"net/http"
)

// Context http 请求上下文
type Context struct {
	Writer  http.ResponseWriter
	Request *http.Request
	Params  map[string]interface{}

	abort bool

	Err error
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

// JSON  将数据 json 编码到response中
// 设置 content type 为 application/json; charset=utf-8
func (ctx *Context) JSON(data interface{}, err error) {
	ctx.Err = err
	bErr := xerror.Cause(err)
	ctx.writeContentType(render.ContentJSON)
	render.Write(render.JSON{
		Code: bErr.Code(),
		Data: data,
		Err:  "",
	}, ctx.Writer)
}

// JSONMap  将数据 json 编码到response中
// 设置 content type 为 application/json; charset=utf-8
func (ctx *Context) JSONMap(data map[string]interface{}, err error) {
	ctx.Err = err
	bErr := xerror.Cause(err)
	data["message"] = bErr.Message()
	data["code"] = bErr.Code()
	ctx.writeContentType(render.ContentJSON)
	render.Write(render.JSONMap(data), ctx.Writer)
}

// String 将字符串写入response body中
// 设置 content type 为 text/plain; charset=utf-8
func (ctx *Context) String(str string, err error) {
	ctx.Err = err
	ctx.writeContentType(render.ContentString)
	render.Write(&render.String{
		Content: str,
	}, ctx.Writer)
}

func (ctx *Context) writeContentType(ctype string) {
	header := ctx.Writer.Header()
	header.Set("Content-Type", ctype)
}

// Handler http 请求处理
type Handler func(*Context)
