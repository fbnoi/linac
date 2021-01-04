package render

import "linac"

const (
	// ContentString content-type
	ContentString = "text/plain; charset=utf-8"
)

// String string response
type String struct {
	Content string
}

// Render Render
func (str *String) Render() (res []byte, err error) {
	return linac.StringToBytes(str.Content), nil
}
