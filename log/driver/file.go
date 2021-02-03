package driver

import (
	"fmt"
	"linac/log/driver/filewriter"
)

var _fileNames = map[int]string{
	_debugIdx: "debug.log",
	_infoIdx:  "info.log",
	_warnIdx:  "warning.log",
	_errorIdx: "error.log",
	_fatalIdx: "fatal.log",
}

// File log file driver
type File struct {
	fws [_totalIdx]*filewriter.Filewriter
}

func (f *File) Write(bs []byte, level int) (int, error) {
	if level >= _totalIdx {
		return 0, fmt.Errorf("unsport log level %d", level)
	}
	n, err := f.fws[level].Write(bs)
	if err != nil {
		return 0, err
	}
	return n, nil
}
