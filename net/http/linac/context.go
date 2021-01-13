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

// Get 获取GET请求参数
// 返回路由参数以及query值，重名是优先路由参数
func (ctx *Context) Get(name string) (value interface{}) {
	value, ok := ctx.Params[name]
	if ok {
		return
	} else if value = ctx.Request.URL.Query().Get(name); value != nil {
		return
	}
	return
}

// Post 获取POST请求参数
func (ctx *Context) Post(name string) (value interface{}) {
	post := ctx.Request.PostForm
	if values, ok := post[name]; ok {
		if len(values) == 1 {
			value = values[0]
		} else {
			value = values
		}
	}
	return
}

// Input 获取请求参数，不管是post还是get
func (ctx *Context) Input(name string) (value interface{}) {
	value = ctx.Get(name)
	if value == nil {
		value = ctx.Post(name)
	}
	return
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
