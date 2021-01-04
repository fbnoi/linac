package render

import (
	"fmt"
	"io"
)

// IRender render
type IRender interface {
	Render() ([]byte, error)
}

// Write 将 render 渲染到 io 中
func Write(render IRender, w io.Writer) {
	bs, err := render.Render()
	if err != nil {
		panic(fmt.Sprintf("render error: %s", err.Error()))
	}
	_, err = w.Write(bs)
	if err != nil {
		panic(fmt.Sprintf("write to io error: %s", err.Error()))
	}
}
