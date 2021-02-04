package log

import (
	"fmt"
	"linac"
	"linac/log/driver"
	"path"
	"runtime"
	"sync"
)

// Driver 日志驱动接口
type Driver interface {
	Write([]byte, int) (int, error)
}

type kv struct {
	key   string
	value interface{}
}

// 日志等级
const (
	LevelDebug = iota
	LevelInfo
	LevelWarn
	LevelError
	LevelFatal
	LevelOff
)

var (
	funcMap sync.Map
)

var (
	_defaultDriver = driver.NewStdOut()
	_mapLevel      = map[int]string{
		LevelDebug: "DEBUG",
		LevelInfo:  "INFO",
		LevelWarn:  "WARN",
		LevelError: "ERROR",
		LevelFatal: "FATAL",
	}
	_log = "log"
)

var (
	_r = &render{}
	_d Driver
	_v int
	_c = &Config{}
)

// Config Config
type Config struct {
	Driver     string
	MaxLogFile int
	RotateSize int64
	Specs      map[string]int
}

func log(level int, kvs ...*kv) {
	if level >= LevelOff || level < 0 {
		return
	}
	pc, file, line, ok := runtime.Caller(3)
	var fname, funcname string
	if ok {
		fname = path.Base(file)
		if v, ok := funcMap.Load(pc); ok {
			funcname = v.(string)
		} else {
			funcname = runtime.FuncForPC(pc).Name()
			funcMap.Store(pc, funcname)
		}
	} else {
		fname = "unknown"
	}
	fl := flevel(fname)
	if level < fl {
		return
	}
	kvs = append(kvs, kV(_fullSource, file), kV(_finSource, fmt.Sprintf("%s:%d", path.Base(file), line)), kV(_function, funcname))
	m := make(map[string]interface{})
	for _, kv := range kvs {
		m[kv.key] = kv.value
	}
	_d.Write(linac.StringToBytes(_r.foramt(m)), level)
}

// Print Print
func Print(sfmt string, v ...interface{}) {
	log(LevelInfo, kV(_message, fmt.Sprintf(sfmt, v...)))
}

// Info Info
func Info(sfmt string, v ...interface{}) {
	log(LevelInfo, kV(_message, fmt.Sprintf(sfmt, v...)))
}

// Debug Debug
func Debug(sfmt string, v ...interface{}) {
	log(LevelDebug, kV(_message, fmt.Sprintf(sfmt, v...)))
}

// Warn Warn
func Warn(sfmt string, v ...interface{}) {
	log(LevelWarn, kV(_message, fmt.Sprintf(sfmt, v...)))
}

// Error Error
func Error(sfmt string, v ...interface{}) {
	log(LevelError, kV(_message, fmt.Sprintf(sfmt, v...)))
}

// Fatal Fatal
func Fatal(sfmt string, v ...interface{}) {
	log(LevelFatal, kV(_message, fmt.Sprintf(sfmt, v...)))
}

func kV(key string, value interface{}) *kv {
	return &kv{key: key, value: value}
}

func flevel(fname string) int {
	return LevelDebug
}
