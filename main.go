package main

import (
	"linac/net/http/linac"
	"log"
)

func main() {
	engine := linac.NewEngine()
	engine.Group("/home", func(router *linac.RouteGroup) *linac.RouteGroup {
		router.GET("/test", func(ctx *linac.Context) {
			log.Print("test")
		})
		router.GET("/:id", func(ctx *linac.Context) {
			log.Print(ctx.Params)
		})
		return router
	}, func(ctx *linac.Context) {
		log.Print("in middleware")
	})
	engine.Group("/index", func(router *linac.RouteGroup) *linac.RouteGroup {
		router.GET("/test", func(ctx *linac.Context) {
			log.Print("123456")
		})
		router.GET("/:id(\\d+)", func(ctx *linac.Context) {
			log.Print("123456")
		})
		return router
	})
	engine.Run(":8089")
}
