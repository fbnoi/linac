package render

import (
	xjson "encoding/json"
)

const (
	// ContextJSON content-type
	ContextJSON = "application/json; charset=utf-8"
)

// JSON 返回json渲染，特定型
type JSON struct {
	Msg  string      `json:"msg"`
	Err  string      `json:"err,omitempty"`
	Data interface{} `json:"data"`
}

func (j *JSON) render() (content []byte, err error) {
	content, err = xjson.Marshal(j)
	return
}

// JSONMap 通用型
type JSONMap map[string]interface{}

func (j *JSONMap) render() (content []byte, err error) {
	content, err = xjson.Marshal(j)
	return
}
