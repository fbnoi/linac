package driver

import (
	"bufio"
	"os"
	"sync"
	"testing"
	"time"
)

func TestWriter(t *testing.T) {
	fw, err := New("./test_data/test.log")
	if err != nil {
		t.Errorf("TestWriter faild, %s", err)
	}
	fw.write([]byte("hello world"))
	time.Sleep(2 * time.Second)
}

func TestErupt(t *testing.T) {
	fw, err := New("./test_data/erupt_test.log")
	if err != nil {
		t.Errorf("TestWriter faild, %s", err)
	}
	wg := &sync.WaitGroup{}
	line := 1000000
	for i := 0; i < line; i++ {
		wg.Add(1)
		go func() {
			if _, err := fw.write([]byte("hello world\n")); err != nil {
				t.Errorf("err: %v", err)
			}
			wg.Done()
		}()
	}
	wg.Wait()
	time.Sleep(3 * time.Second)
	file, _ := os.Open("./test_data/erupt_test.log")
	fd := bufio.NewReader(file)
	count := 0
	for {
		_, err := fd.ReadString('\n')
		if err != nil {
			break
		}
		count++
	}
	if count != line {
		t.Errorf("expected %v, get %v", line, count)
	}
}
