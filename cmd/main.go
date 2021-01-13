package main

import (
	"linac/net/http/linac"
	"log"
)

func main() {
	engine := linac.NewEngine()

	engine.Group("/home", "home", func(router *linac.RouteGroup) *linac.RouteGroup {
		router.GET("/test", "test", func(ctx *linac.Context) {
			ctx.String(200, "hello world")
		})
		router.GET("/:id", "test1", func(ctx *linac.Context) {
			ctx.JSONMap(ctx.Params, nil)
		})
		return router
	}, func(ctx *linac.Context) {
		log.Print("in middleware")
	})
	engine.Group("/index", "index", func(router *linac.RouteGroup) *linac.RouteGroup {
		router.GET("/test", "test1", func(ctx *linac.Context) {
			log.Print("123456")
		})
		router.GET("/:id(\\d+)", "test2", func(ctx *linac.Context) {
			log.Print("123456")
		})
		return router
	})
	engine.Group("/panic", "panic", func(router *linac.RouteGroup) *linac.RouteGroup {
		router.GET("/fin", "test1", func(ctx *linac.Context) {
			panic("Oops! in final")
		})
		router.GET("/mid", "test2", func(ctx *linac.Context) {
			panic("Oops! in middle")
		}, func(ctx *linac.Context) {
			ctx.String(200, "safe!")
		})
		return router
	})
	engine.Run(":8089")
}
