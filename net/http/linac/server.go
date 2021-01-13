package linac

import (
	"log"
	"net/http"
	"sync/atomic"
	"time"
)

var (
	_defaultConfig = &ServerConfig{
		Address:      ":8089",
		Timeout:      time.Second * time.Duration(2),
		ReadTimeout:  time.Second * time.Duration(10),
		WriteTimeout: time.Second * time.Duration(30),
	}
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
		config: &atomic.Value{},
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
func (engine *Engine) GetConfig() (conf *ServerConfig) {
	if conf, ok := engine.config.Load().(*ServerConfig); ok {
		return conf
	}
	return _defaultConfig
}

// Run 运行 http server engine
func (engine *Engine) Run(address string) {
	conf := engine.GetConfig()
	serve := &http.Server{
		Addr:         address,
		Handler:      engine.Router,
		ReadTimeout:  conf.ReadTimeout,
		WriteTimeout: conf.WriteTimeout,
	}
	engine.server.Store(serve)
	log.Print("http server run at:" + address + "...")
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
