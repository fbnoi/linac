package filewriter

import (
	"time"
)

type opt struct {
	RotateFormat   string
	RotateInterval time.Duration
	WriteTimeout   time.Duration
	WriteInterval  time.Duration
	MaxFileSize    int64
	MaxFileList    int
}

//Option Option
type Option func(*opt)

// RotateFormat RotateFormat
func RotateFormat(format string) Option {
	return func(o *opt) {
		o.RotateFormat = format
	}
}

// RotateInterval RotateInterval
func RotateInterval(i time.Duration) Option {
	return func(o *opt) {
		o.RotateInterval = i
	}
}

// WriteTimeout WriteTimeout
func WriteTimeout(i time.Duration) Option {
	return func(o *opt) {
		o.WriteTimeout = i
	}
}

// WriteInterval WriteInterval
func WriteInterval(i time.Duration) Option {
	return func(o *opt) {
		o.WriteInterval = i
	}
}

// MaxFileSize MaxFileSize
func MaxFileSize(i int64) Option {
	return func(o *opt) {
		o.MaxFileSize = i
	}
}

// MaxFileList MaxFileList
func MaxFileList(i int) Option {
	return func(o *opt) {
		o.MaxFileList = i
	}
}
