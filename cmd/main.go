package main

import (
	"linac/net/http/linac"
	"time"
)

func main() {
	engine := linac.NewEngine()
	engine.GET("/timeout", "testtimeout", func(ctx *linac.Context) {
		start := time.Now()
		<-ctx.Done()
		ctx.String(200, "Timeout within %s", time.Since(start))
	})
	engine.Run(":8089")
}
