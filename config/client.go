package config

import (
	"sync/atomic"
)

type config struct {
}

// Namespace namespace
type Namespace struct {
	Name string            `json:"name"`
	Data map[string]string `json:"data"`
}

// Client is config client.
type Client struct {
	ver   int64 // NOTE: for config v1
	data  atomic.Value
	event chan string

	useV2     bool
	watchFile map[string]struct{}
	watchAll  bool
}

// Value 返回 config 值
func (c *Client) Value(key string) (val string, ok bool) {
	var (
		m map[string]*Namespace
		n *Namespace
	)
	if m, ok = c.data.Load().(map[string]*Namespace); !ok {
		return
	}
	if n, ok = m[""]; !ok {
		return
	}
	val, ok = n.Data[key]
	return
}
