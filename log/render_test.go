package log

import (
	"bytes"
	"testing"
	"time"
)

func TestParseLongTime(t *testing.T) {
	r := &render{}
	r.parse("%T")
	d := map[string]interface{}{}
	l := time.Now().Format("15:04:05.000")
	fl := r.foramt(d)
	if l != fl {
		t.Errorf("long time error, expected %s, get %s \n", l, fl)
	}
}

func TestParseShortTime(t *testing.T) {
	r := &render{}
	r.parse("%t")
	d := map[string]interface{}{}
	l := time.Now().Format("15:04:05")
	fl := r.foramt(d)
	if l != fl {
		t.Errorf("short time error, expected %s, get %s \n", l, fl)
	}
}

func TestParseLongDate(t *testing.T) {
	r := &render{}
	r.parse("%D")
	d := map[string]interface{}{}
	l := time.Now().Format("2006/01/02")
	fl := r.foramt(d)
	if l != fl {
		t.Errorf("long date error, expected %s, get %s \n", l, fl)
	}
}

func TestParseShortDate(t *testing.T) {
	r := &render{}
	r.parse("%d")
	d := map[string]interface{}{}
	l := time.Now().Format("01/02")
	fl := r.foramt(d)
	if l != fl {
		t.Errorf("short date error, expected %s, get %s \n", l, fl)
	}
}

func TestDefaultFormatFuncFactory(t *testing.T) {
	str := "[12345 45567]"
	var sli []func(map[string]interface{}) string
	for i := 0; i < len(str); i++ {
		f := defaultFormatFuncFactory(string(str[i]))
		sli = append(sli, f)
	}
	var buf bytes.Buffer
	for _, fun := range sli {
		buf.WriteString(fun(map[string]interface{}{}))
	}
	if str != buf.String() {
		t.Errorf("default func error, expected %s, get %s \n", str, buf.String())
	}
}

func TestLevel(t *testing.T) {
	r := &render{}
	r.parse("%L")
	d := map[string]interface{}{
		"L": "DEBUG",
	}
	l := "DEBUG"
	fl := r.foramt(d)
	if l != fl {
		t.Errorf("log level error, expected %s, get %s \n", l, fl)
	}
}

func TestMessage(t *testing.T) {
	r := &render{}
	r.parse("%M")
	d := map[string]interface{}{
		"L": "DEBUG",
		"M": "test",
	}
	l := "test"
	fl := r.foramt(d)
	if l != fl {
		t.Errorf("message error, expected %s, get %s \n", l, fl)
	}
}

func TestMessageMap(t *testing.T) {
	r := &render{}
	r.parse("%M")
	d := map[string]interface{}{
		"L": "DEBUG",
		"M": map[string]interface{}{
			"foo": "bar",
		},
	}
	l := "foo=bar"
	fl := r.foramt(d)
	if l != fl {
		t.Errorf("message map error, expected %s, get %s \n", l, fl)
	}
}

func TestFunc(t *testing.T) {
	r := &render{}
	r.parse("%f")
	d := map[string]interface{}{
		"f": "func",
	}
	l := "func"
	fl := r.foramt(d)
	if l != fl {
		t.Errorf("function error, expected %s, get %s \n", l, fl)
	}
}

func TestAppId(t *testing.T) {
	r := &render{}
	r.parse("%i")
	d := map[string]interface{}{
		"i": "appid",
	}
	l := "appid"
	fl := r.foramt(d)
	if l != fl {
		t.Errorf("appid error, expected %s, get %s \n", l, fl)
	}
}

func TestEnv(t *testing.T) {
	r := &render{}
	r.parse("%e")
	d := map[string]interface{}{
		"e": "dev",
	}
	l := "dev"
	fl := r.foramt(d)
	if l != fl {
		t.Errorf("env error, expected %s, get %s \n", l, fl)
	}
}

func TestZone(t *testing.T) {
	r := &render{}
	r.parse("%z")
	d := map[string]interface{}{
		"z": "wh01",
	}
	l := "wh01"
	fl := r.foramt(d)
	if l != fl {
		t.Errorf("zone error, expected %s, get %s \n", l, fl)
	}
}

func TestFullSource(t *testing.T) {
	r := &render{}
	r.parse("%S")
	d := map[string]interface{}{
		"S": "/dev/sss.go:80",
	}
	l := "/dev/sss.go:80"
	fl := r.foramt(d)
	if l != fl {
		t.Errorf("full source error, expected %s, get %s \n", l, fl)
	}
}

func TestFinSource(t *testing.T) {
	r := &render{}
	r.parse("%s")
	d := map[string]interface{}{
		"s": "sss.go:80",
	}
	l := "sss.go:80"
	fl := r.foramt(d)
	if l != fl {
		t.Errorf("fin source error, expected %s, get %s \n", l, fl)
	}
}
