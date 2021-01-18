package driver

import (
	"bytes"
	"errors"
	"fmt"
	"log"
	"os"
	"sync"
	"time"
)

var (
	_defaultConfig = &config{
		RotateFormat: "2016-01-02",
		MaxFile:      999,
		MaxSize:      1 << 30,
		ChanSize:     1 << 16,
		RotateTick:   time.Millisecond * 100,
		WriteTimeout: time.Second * 2,
	}
)

// New 返回一个新的
// func New(path string, cnfs ...Config) (*FileWriter, error) {
// 	conf := _defaultConfig
// 	for _, cnf := range cnfs {
// 		cnf(conf)
// 	}

// 	fname := filepath.Base(path)
// 	if fname == "" {
// 		return nil, errors.New("file name connot be nil")
// 	}

// 	dir := filepath.Dir(path)
// 	fi, err := os.Stat(dir)
// 	if err == nil && !fi.IsDir() {
// 		return nil, fmt.Errorf("path %s already exists and not a directory", dir)
// 	}

// 	if os.IsNotExist(err) {
// 		if err = os.MkdirAll(dir, 0755); err != nil {
// 			return nil, fmt.Errorf("create dir %s error: %s", dir, err.Error())
// 		}
// 	}
// }

// FileWriter 文件io
// 采用缓冲池的方式写入文件
// 现将写入的文件放入pool中的buf上，然后传入 chan，FileWriter后台的线程获取buf内容，并将buf放回pool中
type FileWriter struct {
	dir    string
	fname  string
	wc     chan *bytes.Buffer
	pool   *sync.Pool
	closed bool

	writeTimeout time.Duration

	fp *xfile

	// 等待所有的buf被写入文件，平滑关闭文件日志
	wg sync.WaitGroup
}

type xfile struct {
	filezise int64
	fp       *os.File
}

func (x *xfile) write(p []byte) (n int, err error) {
	n, err = x.fp.Write(p)
	x.filezise += int64(n)
	return
}

func (x *xfile) size() (n int, err error) {
	return x.size()
}

func (f *FileWriter) write(bs []byte) (int, error) {
	if f.closed {
		return 0, errors.New("file witer has been closed")
	}
	buf := f.getBuf()
	buf.Write(bs)
	if f.writeTimeout > 0 {
		timeout := time.NewTimer(f.writeTimeout)
		select {
		case f.wc <- buf:
			return len(bs), nil
		case <-timeout.C:
			return 0, errors.New("file witer error, log channel is full")
		}
	}
	select {
	case f.wc <- buf:
		return len(bs), nil
	default:
		return 0, fmt.Errorf("file witer error, log channel is full")
	}
}

func (f *FileWriter) writeProcess() {
	connBuf := &bytes.Buffer{}
	writeTick := time.NewTimer(10 * time.Millisecond)
	var err error
	for {
		select {
		case buf, ok := <-f.wc:
			if ok {
				connBuf.Write(buf.Bytes())
				f.releaseBuf(buf)
			}
		case <-writeTick.C:
			if connBuf.Len() > 0 {
				if err = f.writeToFile(connBuf.Bytes()); err != nil {
					log.Printf("file write error: %s", err)
				}
				connBuf.Reset()
			}
		}

		if false == f.closed {
			continue
		}

		if err := f.writeToFile(connBuf.Bytes()); err != nil {
			log.Printf("file write error: %s", err)
		}
		for buf := range f.wc {
			if err = f.writeToFile(buf.Bytes()); err != nil {
				log.Printf("file write error: %s", err)
			}
			f.releaseBuf(buf)
		}
		break
	}
}

func (f *FileWriter) releaseBuf(buf *bytes.Buffer) {
	buf.Reset()
	f.pool.Put(buf)
}

func (f *FileWriter) getBuf() *bytes.Buffer {
	return f.pool.Get().(*bytes.Buffer)
}

func (f *FileWriter) close() error {
	f.closed = true
	close(f.wc)
	f.wg.Wait()
	return nil
}

func (f *FileWriter) writeToFile(bs []byte) error {
	if f.fp == nil {
		log.Printf("file write error, file handler nil")
		log.Printf("%s", bs)
	}
	_, err := f.fp.write(bs)
	return err
}

// config FileWriter 的配置选项
type config struct {
	RotateFormat string
	MaxFile      int
	MaxSize      int64
	ChanSize     int
	RotateTick   time.Duration
	WriteTimeout time.Duration
}

// Config config
type Config func(conf *config)

// RotateFormat 设置 rotate format
func RotateFormat(format string) Config {
	return func(conf *config) {
		conf.RotateFormat = format
	}
}

// MaxFile 设置 MaxFile
func MaxFile(n int) Config {
	return func(conf *config) {
		conf.MaxFile = n
	}
}

// MaxSize 设置 MaxSize
func MaxSize(size int64) Config {
	return func(conf *config) {
		conf.MaxSize = size
	}
}

// ChanSize 设置 ChanSize
// chan size 不应设置过高以避免程序占用太多的内存
// chan size 默认 1 << 16
func ChanSize(size int) Config {
	return func(conf *config) {
		conf.ChanSize = size
	}
}
