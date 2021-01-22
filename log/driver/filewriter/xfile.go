package filewriter

import (
	"os"
)

func newFile(path string) (*xfile, error) {
	fp, err := os.OpenFile(path, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return nil, err
	}
	fi, err := fp.Stat()
	if err != nil {
		return nil, err
	}
	return &xfile{fp: fp, name: path, size: fi.Size()}, nil
}

type xfile struct {
	fp   *os.File
	name string
	size int64
}

func (x *xfile) write(bs []byte) (int, error) {
	n, err := x.fp.Write(bs)
	if err != nil {
		return 0, err
	}
	x.size += int64(n)
	return n, nil
}

func (x *xfile) close() error {
	return x.fp.Close()
}
