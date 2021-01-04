package linac

import (
	"reflect"
	"unsafe"
)

// StringToBytes string 转化为 byte
func StringToBytes(str string) (bs []byte) {
	sp := *(*reflect.StringHeader)(unsafe.Pointer(&str))
	bp := (*reflect.SliceHeader)(unsafe.Pointer(&bs))
	bp.Data, bp.Len, bp.Cap = sp.Data, sp.Len, sp.Len
	return
}

// BytesToString byte 转化为 string
func BytesToString(bs []byte) string {
	return *(*string)(unsafe.Pointer(&bs))
}
