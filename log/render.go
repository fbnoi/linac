package log

import (
	"fmt"
	"linac"
	"strings"
	"time"
)

var (
	_suportFormet  = []string{"T", "t", "D", "d", "L", "M", "f", "e", "z", "i", "S", "s"}
	_defaultFormat = "[%D %T][%i.%e][%S][%L]%M"
	_mapFormetFunc = map[string]func(map[string]interface{}) string{
		_longTime:   longTime,
		_shortTime:  shortTime,
		_longDate:   longDate,
		_shortDate:  shortDate,
		_level:      keyFactory("L"),
		_function:   keyFactory("f"),
		_env:        keyFactory("e"),
		_zone:       keyFactory("z"),
		_appid:      keyFactory("i"),
		_FullSourse: keyFactory("S"),
		_finSourse:  keyFactory("s"),
		_message:    message,
	}
)

type render struct {
	sli map[string]interface{}
}

func (render *render) parse(format string) {
	if render.sli == nil {
		render.sli = make(map[string]interface{})
	}
	var bs = make([]byte, 1)
	for i := 0; i < len(format); i++ {
		b := format[i]
		if b != '%' {
			continue
		}
		bs = append(bs, b)
		k := linac.BytesToString(bs)
		if fun, ok := _mapFormetFunc[k]; ok {
			render.sli[k] = fun
		}
		bs = bs[:0]
	}
}

func longTime(map[string]interface{}) string {
	return time.Now().Format("00:00:00.000")
}

func shortTime(map[string]interface{}) string {
	return time.Now().Format("00:00:00")
}

func longDate(map[string]interface{}) string {
	return time.Now().Format("2020/01/01")
}

func shortDate(map[string]interface{}) string {
	return time.Now().Format("01/01")
}

func keyFactory(key string) func(map[string]interface{}) string {
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
	var m string
	var s []string
	for k, v := range d {
		if k == _message {
			m = fmt.Sprint(v)
			continue
		}
		s = append(s, fmt.Sprintf("%s=%v", k, v))
	}
	s = append(s, m)
	return strings.Join(s, " ")
}
