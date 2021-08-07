// +build 386 arm

// Package base14 base16384 的 go 接口
package base14

// #cgo CFLAGS: -I/usr/local/include
// #cgo LDFLAGS: -L/usr/local/lib -lbase14
// #include <stdlib.h>
// #include <base14.h>
import "C"
import (
	"unsafe"
)

func Encode(b []byte) []byte {
	ld := C.encode((*C.uchar)(unsafe.Pointer(&b[0])), (uint)(len(b)))
	ldlen := uintptr(ld.len)
	// []byte 的数据结构
	e := [3]uintptr{uintptr(unsafe.Pointer(ld.data)), ldlen, ldlen}
	return *(*[]byte)(unsafe.Pointer(&e))
}

func Decode(b []byte) []byte {
	ld := C.decode((*C.uchar)(unsafe.Pointer(&b[0])), (uint)(len(b)))
	ldlen := uintptr(ld.len)
	// []byte 的数据结构
	e := [3]uintptr{uintptr(unsafe.Pointer(ld.data)), ldlen, ldlen}
	return *(*[]byte)(unsafe.Pointer(&e))
}
