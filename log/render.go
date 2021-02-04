package log

import (
	"bytes"
	"fmt"
	"time"
)

const (
	_longTime   = "T"
	_shortTime  = "t"
	_longDate   = "D"
	_shortDate  = "d"
	_level      = "L"
	_message    = "M"
	_function   = "f"
	_appid      = "i"
	_env        = "e"
	_zone       = "z"
	_fullSource = "S"
	_finSource  = "s"
)

var (
	_defaultFormat = "[%D %T][%i.%e][%S][%L]%M"
	_mapFormetFunc = map[string]func(map[string]interface{}) string{
		_longTime:   longTime,
		_shortTime:  shortTime,
		_longDate:   longDate,
		_shortDate:  shortDate,
		_level:      keyFormatFuncFactory(_level),
		_env:        keyFormatFuncFactory(_env),
		_zone:       keyFormatFuncFactory(_zone),
		_appid:      keyFormatFuncFactory(_appid),
		_fullSource: keyFormatFuncFactory(_fullSource),
		_finSource:  keyFormatFuncFactory(_finSource),
		_function:   keyFormatFuncFactory(_function),
		_message:    keyFormatFuncFactory(_message),
	}
)

type render struct {
	sli []func(map[string]interface{}) string
}

func (render *render) foramt(d map[string]interface{}) string {
	var buf bytes.Buffer
	for _, fun := range render.sli {
		buf.WriteString(fun(d))
	}
	return buf.String()
}

func (render *render) parse(format string) {
	if format == "" {
		return
	}
	var buf bytes.Buffer
	for i := 0; i < len(format); i++ {
		if format[i] != '%' {
			buf.WriteByte(format[i])
			continue
		}
		if i+1 >= len(format) {
			buf.WriteByte(format[i])
			continue
		}
		fun, ok := _mapFormetFunc[string(format[i+1])]
		if !ok {
			buf.WriteByte(format[i])
			continue
		} else {
			if buf.Len() > 0 {
				render.sli = append(render.sli, defaultFormatFuncFactory(buf.String()))
				buf.Reset()
			}
			render.sli = append(render.sli, fun)
			i++
		}
	}
	if buf.Len() > 0 {
		render.sli = append(render.sli, defaultFormatFuncFactory(buf.String()))
		buf.Reset()
	}
}

func longTime(map[string]interface{}) string {
	return time.Now().Format("15:04:05.000")
}

func shortTime(map[string]interface{}) string {
	return time.Now().Format("15:04:05")
}

func longDate(map[string]interface{}) string {
	return time.Now().Format("2006/01/02")
}

func shortDate(map[string]interface{}) string {
	return time.Now().Format("01/02")
}

func defaultFormatFuncFactory(s string) func(map[string]interface{}) string {
	return func(map[string]interface{}) string {
		return s
	}
}

func keyFormatFuncFactory(key string) func(map[string]interface{}) string {
	return func(d map[string]interface{}) string {
		if v, ok := d[key]; ok {
			if s, ok := v.(string); ok {
				return s
			}
			return fmt.Sprint(v)
		}
		return ""
	}
}
