package log

import (
	"fmt"
	"linac"
	"linac/log/driver"
	"runtime"
	"strings"
)

func init() {

}

// 日志等级
const (
	LevelDebug = iota
	LevelInfo
	LevelWarning
	LevelError
	LevelFatal
	LevelOff
)

var (
	_defaultDriver = driver.NewStdOut()
	_mapLevel      = map[int]string{
		LevelDebug:   "DEBUG",
		LevelInfo:    "INFO",
		LevelWarning: "WARNING",
		LevelError:   "ERROR",
		LevelFatal:   "FATAL",
	}
	log *logger
)

// Driver 日志驱动接口
type Driver interface {
	Write([]byte, int) (int, error)
}

// Config 日志配置
type Config struct {
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
	// module set
	module map[string]int
}

type logger struct {
	driver  Driver
	level   int
	render  *render
	context *Config

	attach map[string]interface{}
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

// SetDriver 设置日志驱动
func (l *logger) SetDriver(driver Driver) {
	l.driver = driver
}

// Attach 添加自定义日志选项
// 添加的键值对在每次写入日志时，都会携带写入日志
func (l *logger) Attach(key string, value interface{}) {
	if l.attach == nil {
		l.attach = make(map[string]interface{})
	}
	l.attach[key] = value
}

// SetLevel 设置日志等级
func (l *logger) SetLevel(level int) error {
	if level > LevelOff || level < LevelDebug {
		return fmt.Errorf("SetLevel(%v) error, unknown log level: %v", level, level)
	}
	l.level = level
	return nil
}

func (l *logger) Print(sfmt string, value ...interface{}) {
	str := fmt.Sprintf(sfmt, value...)
	str = l.wrapper(str)
	l.driver.Write(linac.StringToBytes(str), l.level)
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

func (l *logger) wrapper(message interface{}) string {
	m := l.attach
	m["msg"] = message
	d := map[string]interface{}{
		_longTime:   "",
		_shortTime:  "",
		_longDate:   "",
		_shortDate:  "",
		_level:      l.level,
		_appid:      l.context.AppID,
		_env:        l.context.DeployEnv,
		_zone:       l.context.Zone,
		_fullSourse: "",
		_finSourse:  "",
		_function:   "",
		_message:    m,
	}
	return l.render.foramt(d)
}

func (l *logger) sourceFile() (full, fin string, line int) {
	full, line, ok := fileTrace(3)
	if ok {
		arrFile := strings.Split(full, "/")
		fmt.Println(arrFile)
		fin = arrFile[len(arrFile)-1]
	}
	return
}
func fileTrace(dep int) (file string, line int, ok bool) {
	_, file, line, ok = runtime.Caller(dep)
	return
}
