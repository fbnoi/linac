package filewriter

import (
	"bytes"
	"container/list"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
	"sync"
	"time"
)

var _defaultOpts = opt{
	RotateFormat:   "2006-01-02",
	RotateInterval: 30 * time.Second,
	WriteTimeout:   10 * time.Second,
	WriteInterval:  10 * time.Microsecond,
	MaxFileSize:    1 << 18,
	MaxFileList:    999,
}

// New New
func New(path string, opts ...Option) (*Filewriter, error) {
	name := filepath.Base(path)
	if name == "" {
		return nil, fmt.Errorf("file name connot be nil")
	}
	dir := filepath.Dir(path)
	di, err := os.Stat(dir)
	if err == nil && !di.IsDir() {
		return nil, fmt.Errorf("path %s is already exist and is not a dir", dir)
	}
	if os.IsNotExist(err) {
		if err := os.MkdirAll(dir, 755); err != nil {
			return nil, err
		}
	}

	options := &_defaultOpts
	for _, opt := range opts {
		opt(options)
	}

	stdout := log.New(os.Stderr, "file log: ", log.LstdFlags)
	xfile, err := newFile(path)
	if err != nil {
		return nil, err
	}
	fList, err := scanRotateFiles(dir, name, options.RotateFormat)
	lastRotateName := time.Now().Format(options.RotateFormat)
	lastRotateIdx := 1
	if err != nil {
		fList = list.New().Init()
	} else {
		rfi := fList.Front()
		if rfi != nil {
			rf := rfi.Value.(*rotateFile)
			lastRotateName = rf.rotateTime.Format(options.RotateFormat)
			lastRotateIdx = rf.rotateIdx
		}
	}
	wwg := &sync.WaitGroup{}
	cwg := &sync.WaitGroup{}
	w := &Filewriter{
		opt:     options,
		dir:     dir,
		fname:   name,
		bufChan: make(chan *bytes.Buffer),
		pool: &sync.Pool{New: func() interface{} {
			return &bytes.Buffer{}
		}},
		stdOut:         stdout,
		out:            xfile,
		closed:         false,
		lastRotateName: lastRotateName,
		lastRotateIdx:  lastRotateIdx,
		rotateList:     fList,
		wwg:            wwg,
		cwg:            cwg,
	}
	go w.deamon()
	wwg.Add(1)
	return w, nil
}

// Filewriter Filewriter
type Filewriter struct {
	opt   *opt
	dir   string
	fname string

	bufChan chan *bytes.Buffer
	pool    *sync.Pool
	stdOut  *log.Logger
	out     *xfile

	lastRotateName string
	lastRotateIdx  int
	closed         bool
	rotateList     *list.List

	wwg *sync.WaitGroup
	cwg *sync.WaitGroup
}

// deamon
func (f *Filewriter) deamon() {
	eggBuf := &bytes.Buffer{}
	rotateTicker := time.NewTicker(f.opt.RotateInterval)
	writeTicker := time.NewTicker(f.opt.WriteInterval)
	for {
		select {
		case t := <-rotateTicker.C:
			f.checkAndRotate(t)
		case <-writeTicker.C:
			if eggBuf.Len() > 0 {
				if _, err := f.writeToFile(eggBuf.Bytes()); err != nil {
					f.stdOut.Printf("write file error: %s\n", err)
					f.stdOut.Printf("file log: %s \n", eggBuf.Bytes())
				}
				eggBuf.Reset()
			}
		case buf, ok := <-f.bufChan:
			if ok {
				if _, err := eggBuf.Write(buf.Bytes()); err != nil {
					f.stdOut.Printf("write to pool error: %s\n", err)
					f.stdOut.Printf("file log: %s \n", buf.Bytes())
				}
				f.putBuf(buf)
			}
		}
		if f.closed {
			if _, err := f.writeToFile(eggBuf.Bytes()); err != nil {
				f.stdOut.Printf("log file write error: %s", err)
			}
			for buf := range f.bufChan {
				if _, err := f.writeToFile(buf.Bytes()); err != nil {
					f.stdOut.Printf("log file write error: %s", err)
				}
			}
			break
		}
	}
	f.wwg.Done()
}

// checkAndRotate 此方法在后台运行
func (f *Filewriter) checkAndRotate(t time.Time) {
	var err error
	if f.opt.MaxFileList > 0 {
		for f.opt.MaxFileList <= f.rotateList.Len() {
			ofile := f.rotateList.Remove(f.rotateList.Back()).(*rotateFile)
			if err = os.Remove(ofile.path); err != nil {
				f.stdOut.Printf("remove file %s error: %s", ofile.path, err)
			}
		}
	}
	// f.out 在 checkAndRotate 之后可能为空
	if f.out == nil {
		f.stdOut.Printf("check file size error, file nil")
		return
	}
	rotate := t.Format(f.opt.RotateFormat)
	size := f.out.size
	if rotate != f.lastRotateName || (size > f.opt.MaxFileSize && f.opt.MaxFileSize > 0) {
		f.lastRotateIdx++
		if f.lastRotateIdx > f.opt.MaxFileList {
			f.lastRotateIdx = 1
		}
		if err = f.out.close(); err != nil {
			f.stdOut.Printf("close file %s error: %s", f.out.name, err)
		}
		newName := fmt.Sprintf("%s.%s.%03d", f.out.name, rotate, f.lastRotateIdx)
		if err = os.Rename(f.out.name, newName); err != nil {
			f.stdOut.Printf("rename file %s error: %s", f.out.name, err)
			return
		}
		f.rotateList.PushFront(&rotateFile{path: newName, rotateTime: t, rotateIdx: f.lastRotateIdx})

		f.lastRotateName = rotate
		f.out, err = newFile(filepath.Join(f.dir, f.fname))
		if err != nil {
			f.stdOut.Printf("open file %s error: %s", filepath.Join(f.dir, f.fname), err)
		}
	}
}

// Write 写入文件
// 线程安全，将 bs 写入缓存并尝试放入通道，返回写入成功的字节数以及错误
// 当Filewriter已关闭，或者超时将返回错误信息
func (f *Filewriter) Write(bs []byte) (int, error) {
	if f.closed {
		return 0, fmt.Errorf("Filewriter already closed")
	}
	// 在这段到将buf发送到bufChan之间，可能bufChan已经关闭导致致命错误，添加 cwg 等待
	f.cwg.Add(1)
	buf := f.getBuf()
	n, err := buf.Write(bs)
	if err != nil {
		return 0, err
	}
	if f.opt.WriteTimeout > 0 {
		tm := time.NewTimer(f.opt.WriteTimeout)
		select {
		case f.bufChan <- buf:
			f.cwg.Done()
			return n, nil
		case <-tm.C:
			return 0, fmt.Errorf("write file time out")
		}
	}
	select {
	case f.bufChan <- buf:
		f.cwg.Done()
		return n, nil
	default:
		return 0, fmt.Errorf("write file time out")
	}
}

// Close Close
func (f *Filewriter) Close() error {
	f.closed = true
	// cwg 等待关闭 bufChan 后，已经从池中获取的buf写入 bufChan，写完之后关闭
	f.cwg.Wait()
	close(f.bufChan)
	f.wwg.Wait()
	return nil
}

func (f *Filewriter) writeToFile(bs []byte) (int, error) {
	// f.out 在 checkAndRotate 之后可能为空
	if f.out == nil {
		return 0, fmt.Errorf("write file error, file nil")
	}
	n, err := f.out.write(bs)
	if err != nil {
		return 0, err
	}
	return n, nil
}

func (f *Filewriter) getBuf() *bytes.Buffer {
	return f.pool.Get().(*bytes.Buffer)
}

func (f *Filewriter) putBuf(buf *bytes.Buffer) {
	buf.Reset()
	f.pool.Put(buf)
}

func scanRotateFiles(dir, filename, rotateFormat string) (*list.List, error) {
	fis, err := ioutil.ReadDir(dir)
	if err != nil {
		return nil, err
	}
	var ls []*rotateFile
	for _, fi := range fis {
		if fi.IsDir() {
			continue
		}
		name := fi.Name()
		if strings.HasPrefix(name, filename) && name != filename {
			arr := strings.Split(strings.TrimLeft(name[len(filename):], "."), ".")
			if len(arr) != 2 {
				continue
			}
			rtime, err := time.Parse(rotateFormat, arr[0])
			if err != nil {
				continue
			}
			ridx, err := strconv.Atoi(arr[1])
			if err != nil {
				continue
			}
			rfile := &rotateFile{
				path:       filepath.Join(dir, name),
				rotateTime: rtime,
				rotateIdx:  ridx,
			}
			ls = append(ls, rfile)
		}
	}
	sort.Slice(ls, func(i, j int) bool {
		return ls[i].rotateTime.Before(ls[j].rotateTime) || ls[i].rotateIdx < ls[j].rotateIdx
	})
	list := list.New().Init()
	for _, r := range ls {
		list.PushFront(r)
	}
	return list, nil
}

type rotateFile struct {
	path       string
	rotateTime time.Time
	rotateIdx  int
}
