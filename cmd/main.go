package main

import (
	"linac/net/http/linac"
	"log"
	"time"
)

func main() {
	engine := linac.NewEngine()
	engine.GET("/timeout", "testtimeout", func(ctx *linac.Context) {
		start := time.Now()
		ctx.String(200, "Timeout within %s", time.Since(start))
	})
	engine.GET("/test1/:param1", "testGetParam1", func(ctx *linac.Context) {
		for name, value := range ctx.Params {
			log.Printf("%s=%s", name, value)
		}
		ctx.String(200, "Get Param %s", ctx.Get("param1"))
	})
	engine.GET("/test2/:param2", "testGetParam2", func(ctx *linac.Context) {
		url := ctx.Request.URL
		log.Printf("Scheme=%s", url.Scheme)
		log.Printf("Opaque=%s", url.Opaque)
		log.Printf("Host=%s", url.Host)
		log.Printf("Path=%s", url.Path)
		log.Printf("RawPath=%s", url.RawPath)
		log.Printf("ForceQuery=%v", url.ForceQuery)
		log.Printf("RawQuery=%s", url.RawQuery)
		log.Printf("Fragment=%s", url.Fragment)
		log.Printf("RawFragment=%s", url.RawFragment)
		ctx.String(200, "Get Param %s", ctx.Get("param2"))
	})
	engine.GET("/test3/:param3", "testGetParam3", func(ctx *linac.Context) {
		log.Printf("param3=%s", ctx.Get("param3"))
		log.Printf("param2=%s", ctx.Get("param2"))
		ctx.String(200, "Get Param %s", ctx.Post("param2"))
	})
	engine.Run(":8089")
}
