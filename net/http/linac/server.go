package linac

import (
	"errors"
	"net/http"
	"sync/atomic"
	"time"
)

// ServerConfig 服务器配置
type ServerConfig struct {
	Address      string
	Timeout      time.Duration
	ReadTimeout  time.Duration
	WriteTimeout time.Duration
}

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

	config *atomic.Value
}

// SetConfig 设置服务器配置
func (engine *Engine) SetConfig(conf *ServerConfig) {
	engine.config.Store(conf)
}

// GetConfig 获取服务器配置
func (engine *Engine) GetConfig() (conf *ServerConfig, ok bool) {
	conf, ok = engine.config.Load().(*ServerConfig)
	return
}

// Run 运行 http server engine
func (engine *Engine) Run(address string) {
	if conf, ok := engine.GetConfig(); ok {
		serve := &http.Server{
			Addr:         address,
			Handler:      engine.Router,
			ReadTimeout:  conf.ReadTimeout,
			WriteTimeout: conf.WriteTimeout,
		}
		engine.server.Store(serve)
		if err := serve.ListenAndServe(); err != nil {
			panic(err)
		}
	} else {
		panic(errors.New("server engine should have server config, get nil"))
	}
}

// Server 返回 engine 的 http server
func (engine *Engine) Server() *http.Server {
	if server, ok := engine.server.Load().(*http.Server); ok {
		return server
	}
	return nil
}
