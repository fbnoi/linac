package linac

import (
	"net/http"
	"sync/atomic"
)

// NewEngine 返回一个新的 http server engine
func NewEngine() *Engine {
	engine := &Engine{
		Router: NewRouter(),
		server: &atomic.Value{},
	}
	engine.Router.engine = engine
	engine.Use(Recovery())
	return engine
}

// Engine http server engine
type Engine struct {
	*Router
	server *atomic.Value
}

// Run 运行 http server engine
func (engine *Engine) Run(address string) {
	serve := &http.Server{
		Addr:    address,
		Handler: engine.Router,
	}
	engine.server.Store(serve)
	if err := serve.ListenAndServe(); err != nil {
		panic(err)
	}
}

// Server 返回 engine 的 http server
func (engine *Engine) Server() *http.Server {
	if server, ok := engine.server.Load().(*http.Server); ok {
		return server
	}
	return nil
}
