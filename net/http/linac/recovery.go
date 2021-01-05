package linac

import (
	"fmt"
	"net/http/httputil"
	"os"
	"runtime"

	"github.com/prometheus/common/log"
)

// Recovery 从 panic 中恢复 server，并记录下 ctx
func Recovery() Handler {
	return func(c *Context) {
		defer func() {
			var rawReq []byte
			if err := recover(); err != nil {
				const size = 64 << 10
				buf := make([]byte, size)
				buf = buf[:runtime.Stack(buf, false)]
				if c.Request != nil {
					rawReq, _ = httputil.DumpRequest(c.Request, false)
				}
				pl := fmt.Sprintf("http call panic: %s\n%v\n%s\n", string(rawReq), err, buf)
				fmt.Fprintf(os.Stderr, pl)
				log.Error(pl)
				c.Abort(500)
			}
		}()
		c.Next()
	}
}
