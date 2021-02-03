package driver

import (
	"linac/log/driver/filewriter"
	"path/filepath"
)

const (
	_debugIdx = iota
	_infoIdx
	_warnIdx
	_errorIdx
	_fatalIdx
	_totalIdx
)

// NewStdOut new a stdout driver
func NewStdOut() *StdOut {
	io := &StdOut{}
	for idx, o := range _stdOut {
		io.ios[idx] = o
	}
	return io
}

// NewFile new a file driver
func NewFile(dir string, bufferSize, rotateSize int64, maxLogFile int) *File {
	io := &File{}
	var err error
	for idx, f := range _fileNames {
		path := filepath.Join(dir, f)
		io.fws[idx], err = filewriter.New(path, filewriter.MaxFileSize(rotateSize), filewriter.MaxFileList(maxLogFile))
		if err != nil {
			panic(err)
		}
	}
	return io
}
