package linac

import (
	"fmt"
	"net/http"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

var once sync.Once
var SockAddr = "localhost:8089"

func TestTimeOut(t *testing.T) {
	once.Do(serv)
	t.Run("Should timeout by default", func(t *testing.T) {
		c := &http.Client{}
		req, err := http.NewRequest("GET", uri(SockAddr, "/timeout"), nil)
		assert.Nil(t, err)
		start := time.Now()
		resp, err := c.Do(req)
		assert.Nil(t, err)
		assert.Equal(t, 200, resp.StatusCode)
		assert.InDelta(t, float64(time.Second*2), float64(time.Since(start)), float64(time.Second))
	})

	t.Run("Should timeout by methodConfig", func(t *testing.T) {
		c := &http.Client{}
		req, err := http.NewRequest("GET", uri(SockAddr, "/timeout-method-config"), nil)
		assert.Nil(t, err)
		start := time.Now()
		resp, err := c.Do(req)
		assert.Nil(t, err)
		assert.Equal(t, 200, resp.StatusCode)
		assert.InDelta(t, float64(time.Second*3), float64(time.Since(start)), float64(time.Second))
	})
}

func serv() {
	engine := NewEngine()
	regist(engine)
	go engine.Run(":8089")
}

func regist(engine *Engine) {
	engine.GET("/timeout", "testtimeoutdefault", func(ctx *Context) {
		start := time.Now()
		<-ctx.Done()
		ctx.String(200, "Timeout within %s", time.Since(start))
	})
	engine.GET("/timeout-method-config", "testtimeoutconfig", func(ctx *Context) {
		start := time.Now()
		<-ctx.Done()
		ctx.String(200, "Timeout within %s", time.Since(start))
	}).SetConfig(&RouteConfig{Timeout: time.Second * 3})
}

func uri(base, path string) string {
	return fmt.Sprintf("%s://%s%s", "http", base, path)
}
