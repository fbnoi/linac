package render

import (
	"io"
)

// IRender render
type IRender interface {
	Render() ([]byte, error)
}

// Write 将 render 渲染到 io 中
func Write(render IRender, w io.Writer) error {
	bs, err := render.Render()
	if err != nil {
		return err
	}
	_, err = w.Write(bs)
	if err != nil {
		return err
	}
	return nil
}
