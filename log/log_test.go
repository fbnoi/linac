package log

import (
	"fmt"
	"testing"
)

func TestLogWrapper(t *testing.T) {
	logger := &logger{
		level:   LevelDebug,
		render:  &render{},
		context: &context{},
		attach:  make(map[string]interface{}),
	}
	logger.SetFormat("[%D %t][%f]%M 1213")
	fl := logger.wrapper("test")
	d := make(map[string]interface{})
	m := make(map[string]interface{})
	m["msg"] = "test"
	d["M"] = m
	l := fmt.Sprintf("[%s %s]%s 1213", longDate(d), shortTime(d), message(d))
	if fl != l {
		t.Errorf("log.wrapper error, expected %s, get %s \n", l, fl)
	}
}
