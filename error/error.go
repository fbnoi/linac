package error

import (
	"fmt"
	"strconv"
	"sync/atomic"

	"github.com/pkg/errors"
)

var (
	_messages = atomic.Value{}     // NOTE: struct map[int]string
	_codes    = map[int]struct{}{} // register codes.
)

// Register 注册错误码信息.
func Register(cm map[int]string) {
	_messages.Store(cm)
}

// New 注册一个新的错误
// 新的错误码必须是全局独一无二的，否则 panic
func New(e int) Error {
	if e <= 0 {
		panic("error code must greater than zero")
	}
	return add(e)
}

func add(e int) Error {
	if _, ok := _codes[e]; ok {
		panic(fmt.Sprintf("error: %d already exist", e))
	}
	_codes[e] = struct{}{}
	return Int(e)
}

//IError error 接口
type IError interface {
	Error() string
	Code() int
	Message() string
}

// Error error code
type Error int

// Error Error
func (code Error) Error() string {
	return strconv.FormatInt(int64(code), 10)
}

// Code error code
func (code Error) Code() int {
	return int(code)
}

// Message Message
func (code Error) Message() string {
	if cm, ok := _messages.Load().(map[int]string); ok {
		if msg, ok := cm[code.Code()]; ok {
			return msg
		}
	}
	return code.Error()
}

// Cause 将错误转化为错误码
func Cause(e error) IError {
	if e == nil {
		return OK
	}
	ec, ok := errors.Cause(e).(Error)
	if ok {
		return ec
	}
	return String(e.Error())
}

// EqualError 错误是否一致
func EqualError(code IError, err error) bool {
	return Cause(err).Code() == code.Code()
}

// String 将字符串转化为error
// error 的 error code 不总是为int，这里进行转化
func String(e string) IError {
	if e == "" {
		return OK
	}
	// try error string
	i, err := strconv.Atoi(e)
	if err != nil {
		return ServerErr
	}
	return Error(i)
}

// Int 将数字转化为Error
func Int(i int) Error { return Error(i) }
