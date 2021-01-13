package linac

import (
	"context"
	"fmt"
	xerror "linac/error"
	"linac/net/http/linac/render"
	"net/http"
)

// Context http 请求上下文
type Context struct {
	context.Context

	Writer   http.ResponseWriter
	Request  *http.Request
	Params   map[string]interface{}
	Handlers []Handler
	Error    error

	abort bool
	index int
}

// Next 继续执行下一个handler
// Note: 此方法应该只在中间件中调用
func (ctx *Context) Next() {
	ctx.index++
	for ; ctx.index <= len(ctx.Handlers)-1; ctx.index++ {
		if ctx.IsAbort() {
			return
		}
		ctx.Handlers[ctx.index](ctx)
	}
}

// Abort 停止继续使用handlers处理ctx，但不会停止当前的handler
func (ctx *Context) Abort(code int) {
	ctx.Writer.WriteHeader(code)
	ctx.abort = true
}

// IsAbort 返回context是否终止响应
func (ctx *Context) IsAbort() bool {
	return ctx.abort
}

func (ctx *Context) render(r render.IRender, code int) {
	ctx.Writer.WriteHeader(code)
	ctx.writeContentType(r.ContentType())
	if err := render.Write(r, ctx.Writer); err != nil {
		ctx.Error = err
	}
}

// JSON  将数据 json 编码到response中
// 设置 content type 为 application/json; charset=utf-8
func (ctx *Context) JSON(data interface{}, err error) {
	ctx.Error = err
	bErr := xerror.Cause(err)
	ctx.render(render.JSON{
		Code: bErr.Code(),
		Data: data,
		Err:  bErr.Message(),
	}, http.StatusOK)
}

// JSONMap  将数据 json 编码到response中
// 设置 content type 为 application/json; charset=utf-8
func (ctx *Context) JSONMap(data map[string]interface{}, err error) {
	ctx.Error = err
	bErr := xerror.Cause(err)
	data["message"] = bErr.Message()
	data["code"] = bErr.Code()
	ctx.render(render.JSONMap(data), http.StatusOK)
}

// String 将字符串写入response body中
// 设置 content type 为 text/plain; charset=utf-8
func (ctx *Context) String(code int, sfmt string, value ...interface{}) {
	ctx.render(&render.String{
		Content: fmt.Sprintf(sfmt, value...),
	}, code)
}

func (ctx *Context) writeContentType(ctype string) {
	header := ctx.Writer.Header()
	header.Set("Content-Type", ctype)
}

// Handler http 请求处理
type Handler func(*Context)
