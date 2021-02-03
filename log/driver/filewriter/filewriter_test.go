package filewriter

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"os"
	"strconv"
	"testing"
	"time"
)

func TestWriteFile(t *testing.T) {
	os.Remove("./test_data/test1.log")
	w, err := New("./test_data/test1.log")
	if err != nil {
		t.Error(err)
	}
	w.Write([]byte("hello world\n"))
	w.Close()
	ex, err := ioutil.ReadFile("./test_data/test1.log")
	if err != nil {
		t.Error(err)
	}
	if string(ex) != "hello world" {
		t.Errorf(fmt.Sprintf("expcted hello world, get %s", ex))
	}
}

func TestEruptWrite(t *testing.T) {
	os.Remove("./test_data/test2.log")
	w, err := New("./test_data/test2.log", WriteTimeout(5*time.Second))
	if err != nil {
		t.Error(err)
	}
	for i := 0; i < 100000; i++ {
		go func() {
			if _, err := w.Write([]byte("hello world\n")); err != nil {
				t.Error(err)
			}
		}()
	}
	w.Close()
	file, _ := os.Open("./test_data/test2.log")
	fd := bufio.NewReader(file)
	count := 0
	for {
		_, err := fd.ReadString('\n')
		if err != nil {
			break
		}
		count++
	}
	if 100000 != count {
		t.Errorf("expected %v, get %v", 100000, count)
	}
}

func TestScanFile(t *testing.T) {
	os.Create("./test_data/test.log3.2021-01-19.001")
	os.Create("./test_data/test.log3.2021-01-19.002")
	os.Create("./test_data/test.log3.2021-01-19.003")
	fl, err := scanRotateFiles("./test_data/", "test.log3", "2006-01-02")
	if err != nil {
		t.Error(err)
	}
	if fl.Len() != 3 {
		t.Error(fmt.Errorf("scan error, expected 3 files, get %v", fl.Len()))
	}
}

func TestRotateFile(t *testing.T) {
	os.Remove("./test_data/test4.log")
	w, err := New("./test_data/test4.log", WriteTimeout(5*time.Second), MaxFileSize(1<<20), RotateInterval(500*time.Millisecond), MaxFileList(5))
	if err != nil {
		t.Error(err)
	}
	for i := 0; i < 2000000; i++ {
		if _, err := w.Write([]byte("hello world" + strconv.Itoa(i) + "\n")); err != nil {
			t.Error(err)
		}
	}
	w.Close()
	if w.rotateList.Len() != 5 {
		t.Error(fmt.Errorf("rotate file error, expected 5 files, get %v", w.rotateList.Len()))
	}
}
