package driver

import (
	"fmt"
	"log"
	"os"
)

var _stdOut = map[int]*log.Logger{
	_debugIdx: log.New(os.Stdout, "[DEBUG]:", log.LstdFlags),
	_infoIdx:  log.New(os.Stdout, "[INFO]:", log.LstdFlags),
	_warnIdx:  log.New(os.Stdout, "[WARN]:", log.LstdFlags),
	_errorIdx: log.New(os.Stderr, "[ERROR]:", log.LstdFlags),
	_totalIdx: log.New(os.Stderr, "[FATAL]:", log.LstdFlags),
}

// StdOut log stdout driver
type StdOut struct {
	ios [_totalIdx]*log.Logger
}

// Write write to stdio
func (so *StdOut) Write(bs []byte, level int) (int, error) {
	if level >= _totalIdx {
		return 0, fmt.Errorf("unsport log level %d", level)
	}
	so.ios[level].Print(bs)
	return len(bs), nil
}
