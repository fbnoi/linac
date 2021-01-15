package log

import (
	"bytes"
	"fmt"
	"strings"
	"time"
)

var (
	_mapFormetFunc = map[string]func(map[string]interface{}) string{
		_longTime:   longTime,
		_shortTime:  shortTime,
		_longDate:   longDate,
		_shortDate:  shortDate,
		_level:      keyFormatFuncFactory(_level),
		_function:   keyFormatFuncFactory(_function),
		_env:        keyFormatFuncFactory(_env),
		_zone:       keyFormatFuncFactory(_zone),
		_appid:      keyFormatFuncFactory(_appid),
		_fullSourse: keyFormatFuncFactory(_fullSourse),
		_finSourse:  keyFormatFuncFactory(_finSourse),
		_message:    message,
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

func message(d map[string]interface{}) string {
	var s []string
	if m, ok := d[_message]; ok {
		if mv, ok := m.(map[string]interface{}); ok {
			for k, v := range mv {
				s = append(s, fmt.Sprintf("%s=%v", k, v))
			}
		} else {
			s = append(s, fmt.Sprint(m))
		}
	}
	return strings.Join(s, " ")
}
