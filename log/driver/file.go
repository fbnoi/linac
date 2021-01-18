package driver

import (
	"os"
)

const (
	_debug = iota
	_info
	_warning
	_error
	_fatal
	_all
)

var (
	fm = map[int]string{
		_debug:   "debug.log",
		_info:    "info.log",
		_warning: "warning.log",
		_error:   "error.log",
		_fatal:   "fatal.log",
	}
)

var filePattern = "2006-01-02"

type fileDriver struct {
	path string
	fps  [_all]*os.File
}

func (driver *fileDriver) close() {
	for _, fp := range driver.fps {
		fp.Close()
	}
}
