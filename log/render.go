package log

import (
	"bytes"
	"fmt"
	"path"
	"runtime"
	"strconv"
	"strings"
	"sync"
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
	_fullSourse = "S"
	_finSourse  = "s"
)

var (
	funcMap        sync.Map
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
		_fullSourse: fullSource,
		_finSourse:  finSource,
		_function:   funcName,
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

func fullSource(map[string]interface{}) string {
	if _, file, line, ok := runtime.Caller(3); ok {
		return fmt.Sprintf("%s:%d", file, line)
	}
	return "unknown:0"
}

func finSource(map[string]interface{}) string {
	if _, file, line, ok := runtime.Caller(3); ok {
		return fmt.Sprintf("%s:%d", path.Base(file), line)
	}
	return "unknown:0"
}

func funcName(map[string]interface{}) (name string) {
	if pc, _, line, ok := runtime.Caller(3); ok {
		if v, ok := funcMap.Load(pc); ok {
			name = v.(string)
		} else {
			name = runtime.FuncForPC(pc).Name() + ":" + strconv.FormatInt(int64(line), 10)
			funcMap.Store(pc, name)
		}
	}
	return
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
