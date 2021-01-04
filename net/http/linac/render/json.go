package render

import (
	xjson "encoding/json"
)

const (
	// ContentJSON content-type
	ContentJSON = "application/json; charset=utf-8"
)

// JSON 返回json渲染，特定型
type JSON struct {
	Code int         `json:"code"`
	Err  string      `json:"err,omitempty"`
	Data interface{} `json:"data"`
}

// Render Render
func (j JSON) Render() (content []byte, err error) {
	content, err = xjson.Marshal(j)
	return
}

// JSONMap 通用型
type JSONMap map[string]interface{}

// Render Render
func (j JSONMap) Render() (content []byte, err error) {
	content, err = xjson.Marshal(j)
	return
}
