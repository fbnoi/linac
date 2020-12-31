package linac

import "net/http"

// Context http 请求上下文
type Context struct {
	Writer  http.ResponseWriter
	Request *http.Request
	params  map[string]string
}

// Handler http 请求处理
type Handler func(*Context)
