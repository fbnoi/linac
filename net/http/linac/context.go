package linac

import "net/http"

// Context http 请求上下文
type Context struct {
	writer  http.ResponseWriter
	request *http.Request
}

// Handler http 请求处理
type Handler func(*Context)
