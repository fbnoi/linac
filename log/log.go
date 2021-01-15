package log

import (
	"fmt"
	"linac"
)

// 日志等级
const (
	LevelAll = iota
	LevelDebug
	LevelInfo
	LevelWarning
	LevelError
	LevelFatal
	LevelOff
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
	_FullSourse = "S"
	_finSourse  = "s"
)

var (
	_defaultFormat = "[%D %T][%i.%e][%S][%L]%M"
	_mapLevel      = map[int]string{
		LevelDebug:   "DEBUG",
		LevelInfo:    "INFO",
		LevelWarning: "WARNING",
		LevelError:   "ERROR",
		LevelFatal:   "FATAL",
	}
)

// Driver 日志驱动接口
type Driver interface {
	Write([]byte) error
}

// context 日志上下文
type context struct {
	// Region 地区
	Region string
	// Zone 可用域
	Zone string
	// Hostname 主机名
	Hostname string
	// DeployEnv 部署环境
	DeployEnv string
	// IP 服务IP
	IP string
	// AppID 服务ID
	AppID string
	// AppID 服务名
	AppName string
	// 自定义
	Attech map[string]string
}

type logger struct {
	driver  Driver
	level   int
	render  *render
	context *context
}

// 写入日志
func (l *logger) log(level int, str string) {
	if level < l.level {
		return
	}
	strLevel, ok := _mapLevel[level]
	if !ok {
		fmt.Printf("unsport log level %s", strLevel)
		return
	}
}

func (l *logger) Print(sfmt string, value ...interface{}) {
	str := fmt.Sprintf(sfmt, value...)
	l.driver.Write(linac.StringToBytes(str))
}

// SetFormat
// %T time format at "15:04:05.999" on stdout handler, "15:04:05 MST" on file handler
// %t time format at "15:04:05" on stdout handler, "15:04" on file on file handler
// %D data format at "2006/01/02"
// %d data format at "01/02"
// %L log level e.g. INFO WARN ERROR
// %M log message and additional fields: key=value this is log message
// %f function name and line number e.g. model.Get:121
// %i appid
// %e deploy env e.g. dev prod
// %z zone
// %S full file name and line number: /a/b/c/d.go:23
// %s final file name element and line number: d.go:23
func (l *logger) SetFormat(format string) {
	l.render.parse(format)
}
